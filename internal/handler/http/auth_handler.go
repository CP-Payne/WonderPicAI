package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/CP-Payne/wonderpicai/internal/config"
	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/service"
	"github.com/CP-Payne/wonderpicai/internal/validation"
	authComponents "github.com/CP-Payne/wonderpicai/web/template/components/auth"
	"github.com/CP-Payne/wonderpicai/web/template/components/ui"
	authPages "github.com/CP-Payne/wonderpicai/web/template/pages/auth"
	"github.com/CP-Payne/wonderpicai/web/template/viewmodel"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService service.AuthService
	logger      *zap.Logger
	validate    *validator.Validate
}

func NewAuthHandler(authService service.AuthService, logger *zap.Logger, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{authService: authService, logger: logger.With(zap.String("component", "AuthHandler")), validate: validate}
}

type SignupRequest struct {
	Username        string `validate:"required,alphanum,min=3,max=30"`
	Email           string `validate:"required,email"`
	Password        string `validate:"required,passwordcomplexity"`
	ConfirmPassword string `validate:"required,eqfield=Password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token,omitempty"` // omitempty for registration response
}

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func (h *AuthHandler) ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	err := authPages.AuthPage(authComponents.LoginForm()).Render(r.Context(), w)
	if err != nil {
		h.logger.Error("Failed to render login page", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) ShowSignupPage(w http.ResponseWriter, r *http.Request) {
	// Viewmodel empty on initial load
	vm := viewmodel.SignupFormComponentData{
		Form: viewmodel.SignupFormData{
			Username: "",
			Email:    "",
		},
		Errors: make(map[string]string),
		Error:  "",
	}

	err := authPages.AuthPage(authComponents.SignupForm(vm)).Render(r.Context(), w)
	if err != nil {
		h.logger.Error("Failed to render login page", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.logger.Error("Failed to parse form", zap.Error(err))
		http.Error(w, "Bad Request: Could not parse form data", http.StatusBadRequest)
		return
	}

	req := SignupRequest{
		Username:        r.FormValue("username"),
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}
	// TODO: Do validations and Generate View
	vm := viewmodel.SignupFormComponentData{
		Form: viewmodel.SignupFormData{
			Username: req.Username,
			Email:    req.Email,
		},
		Errors: make(map[string]string),
		Error:  "",
	}

	err := h.validate.Struct(req)
	if err != nil {
		fieldErrors, generalValError := validation.TranslateValidationErrors(err)
		vm.Errors = fieldErrors
		vm.Error = generalValError

		if vm.Error != "" {
			h.logger.Error("General validation error", zap.String("error", vm.Error), zap.Error(err))
			toastData := viewmodel.ToastComponentData{
				Message: vm.Error,
				Type:    viewmodel.ToastError,
				ToastID: xid.New().String(),
			}

			err := ui.ToastNotification(toastData).Render(r.Context(), w)
			if err != nil {
				h.logger.Error("Failed to render toast notification", zap.Error(err))
				// No message = default
				HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
				// Don't return yet, first render the form with other potential validation errors
			}
		}

		h.logger.Warn("Signup validation errors", zap.Any("errors", vm.Errors))

		err = authComponents.SignupForm(vm).Render(r.Context(), w)
		if err != nil {
			h.logger.Error("Failed to render login page", zap.Error(err))
			HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
			return
		}
		// Toast and Form now ready to be returned
		return
	}

	// -- Validation Passed ---

	// Clear validation errors
	vm.Errors = make(map[string]string)
	vm.Error = ""

	user, token, err := h.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			vm.Errors["email"] = "This email address is already registered."
		} else {
			vm.Error = "An error occured while creating your account. Please try again"
			h.logger.Error("Unexpected error from AuthService.Register", zap.Error(err))
		}
		// Consider using error page or Toast for the general, unexpected error
		err = authComponents.SignupForm(vm).Render(r.Context(), w)
		if err != nil {
			h.logger.Error("Failed to render login page", zap.Error(err))
			HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
			return
		}
	}

	// TODO: Set token and redirect user to home page
	h.logger.Info("User registered successfully", zap.String("userID", user.ID.String()), zap.String("email", user.Email))

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		MaxAge:   config.Cfg.JWT.ExpiryMinutes * 60,
		Expires:  time.Now().Add(time.Duration(config.Cfg.JWT.ExpiryMinutes) * time.Minute),
	})

	HxRedirect(w, r, "/")
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invlid request body", http.StatusBadRequest)
		return
	}

	user, token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		// TODO: Differentiate errors
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		User:  UserResponse{ID: user.ID, Username: user.Username},
		Token: token,
	})
}
