package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"unicode"

	"github.com/CP-Payne/wonderpicai/internal/service"
	authComponents "github.com/CP-Payne/wonderpicai/web/template/components/auth"
	authPages "github.com/CP-Payne/wonderpicai/web/template/pages/auth"
	"github.com/CP-Payne/wonderpicai/web/template/viewmodel"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService service.AuthService
	logger      *zap.Logger
	validate    *validator.Validate
}

func NewAuthHandler(authService service.AuthService, logger *zap.Logger, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{authService: authService, logger: logger, validate: validate}
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

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
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
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				fieldName := fieldError.Field()
				tag := fieldError.Tag()

				var errorMessage string

				switch tag {
				case "required":
					errorMessage = fmt.Sprintf("%s is required.", fieldName)
				case "alphanum":
					errorMessage = fmt.Sprintf("%s must contain only letters and numbers.", fieldName)
				case "email":
					errorMessage = "Please enter a valid email address."
				case "min":
					errorMessage = fmt.Sprintf("%s must be at least %s characters long.", fieldName, fieldError.Param())
				case "max":
					errorMessage = fmt.Sprintf("%s cannot be longer than %s characters.", fieldName, fieldError.Param())
				case "eqfield":
					errorMessage = "Passwords do not match."
				case "passwordcomplexity":
					errorMessage = "Password must be at least 8 characters long and include an uppercase letter, lowercase letter, number, and special character."
				default:
					errorMessage = fmt.Sprintf("%s is invalid.", fieldName)
				}

				if existing, ok := vm.Errors[fieldName]; ok {
					vm.Errors[lowerFirst(fieldName)] = existing + " " + errorMessage
				} else {
					vm.Errors[lowerFirst(fieldName)] = errorMessage
				}
			}
		} else {
			h.logger.Error("Validation error", zap.Error(err))
			vm.Error = "An unexpected error occurred during validation."
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `<div id="form-errors" class="text-error">%s</div>`, vm.Error)
			return
		}

		err = authComponents.SignupForm(vm).Render(r.Context(), w)
		if err != nil {
			h.logger.Error("Failed to render login page", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return

	}

	user, err := h.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		// TODO: implement error handling to differentiate different types of errors
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{
		User: UserResponse{ID: user.ID, Username: user.Username},
	})
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
