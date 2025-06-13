package http

import (
	"errors"
	"net/http"

	"github.com/CP-Payne/wonderpicai/internal/config"
	"github.com/CP-Payne/wonderpicai/internal/domain"
	"github.com/CP-Payne/wonderpicai/internal/handler/http/response"
	"github.com/CP-Payne/wonderpicai/internal/service"
	"github.com/CP-Payne/wonderpicai/internal/validation"
	authComponents "github.com/CP-Payne/wonderpicai/web/template/components/auth"
	authPages "github.com/CP-Payne/wonderpicai/web/template/pages/auth"
	"github.com/CP-Payne/wonderpicai/web/template/viewmodel"
	"github.com/go-playground/validator/v10"
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
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

func (h *AuthHandler) ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	// Viewmodel empty on initial load
	vm := viewmodel.LoginFormComponentData{
		Form: viewmodel.LoginFormData{
			Email: "",
		},
		Errors: make(map[string]string),
		Error:  "",
	}
	err := authPages.AuthPage(authComponents.LoginForm(vm), config.Cfg.GoogleAuth.ClientSecret).Render(r.Context(), w)
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

	err := authPages.AuthPage(authComponents.SignupForm(vm), config.Cfg.GoogleAuth.ClientSecret).Render(r.Context(), w)
	if err != nil {
		h.logger.Error("Failed to render login page", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.logger.Error("Failed to parse form", zap.Error(err))
		response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, xid.New().String(), "")
		return
	}

	req := SignupRequest{
		Username:        r.FormValue("username"),
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}
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
			h.logger.Error("General signup validation error", zap.String("error", vm.Error), zap.Error(err))

			toastID, loadErr := response.LoadErrorToast(w, r, h.logger, vm.Error)
			if loadErr != nil {
				response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, toastID, "")
				return
			}
			// Load form with validation errors (if any) before ending request processing
		}

		h.logger.Warn("Signup validation errors", zap.Any("errors", vm.Errors))

		loadErr := response.LoadSignupForm(w, r, h.logger, vm)
		if loadErr != nil {
			response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
			return
		}

		// Toast and SignupForm written to response and ready to end request processing
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
		loadErr := response.LoadSignupForm(w, r, h.logger, vm)
		if loadErr != nil {
			response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
			return
		}
		return
	}

	h.logger.Info("User registered successfully", zap.String("userID", user.ID.String()), zap.String("email", user.Email))

	response.SetAuthCookie(w, r, token)

	response.HxRedirect(w, r, "/gen")
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.logger.Error("Failed to parse form", zap.Error(err))
		response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, xid.New().String(), "")
		return
	}

	req := LoginRequest{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	vm := viewmodel.LoginFormComponentData{
		Form: viewmodel.LoginFormData{
			Email: req.Email,
		},
		Errors: make(map[string]string),
		Error:  "",
	}

	err := h.validate.Struct(req) // This seems to be throwing a null pointer
	if err != nil {
		fieldErrors, generalValError := validation.TranslateValidationErrors(err)
		vm.Errors = fieldErrors
		vm.Error = generalValError

		if vm.Error != "" {
			h.logger.Error("General login validation error", zap.String("error", vm.Error), zap.Error(err))

			toastID, loadErr := response.LoadErrorToast(w, r, h.logger, vm.Error)
			if loadErr != nil {
				response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, toastID, "")
				return
			}

		}

		h.logger.Warn("Login validation errors", zap.Any("errors", vm.Errors))

		loadErr := response.LoadLoginForm(w, r, h.logger, vm)
		if loadErr != nil {
			response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
			return
		}
		// Toast and Form loaded and ready to be returned
		return
	}

	// -- Validation Passed ---
	//
	// Clear validation errors
	vm.Errors = make(map[string]string)
	vm.Error = ""

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			vm.Error = "Invalid Credentials"

			loadErr := response.LoadLoginForm(w, r, h.logger, vm)
			if loadErr != nil {
				response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
				return
			}

			return
		}

		vm.Error = "Something went wrong. Please try again."

		h.logger.Error("General login validation error", zap.Error(err))

		toastID, loadErr := response.LoadErrorToast(w, r, h.logger, vm.Error)
		if loadErr != nil {
			response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, toastID, "")
			return
		}

		loadErr = response.LoadLoginForm(w, r, h.logger, vm)
		if loadErr != nil {
			response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
			return
		}

		// Toast and LoginForm rendered (loaded) end request processing
		return
	}

	h.logger.Info("User authenticated successfully", zap.String("userID", user.ID.String()))

	response.SetAuthCookie(w, r, token)
	response.HxRedirect(w, r, "/gen")
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {

	response.SetEmptyAuthCookie(w, r)

	response.HxRedirect(w, r, "/auth/login")
}

func (h *AuthHandler) HandleExternalAuth(w http.ResponseWriter, r *http.Request) {

	vm := viewmodel.LoginFormComponentData{
		Form:   viewmodel.LoginFormData{},
		Errors: make(map[string]string),
		Error:  "",
	}

	user, token, err := h.authService.HandleExternalAuthCallback(r)
	if err != nil {

		vm.Error = "Something went wrong. Please try again."

		h.logger.Error("General login validation error", zap.Error(err))

		toastID, loadErr := response.LoadErrorToast(w, r, h.logger, vm.Error)
		if loadErr != nil {
			response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, toastID, "")
			return
		}

		loadErr = response.LoadLoginForm(w, r, h.logger, vm)
		if loadErr != nil {
			response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
			return
		}

		return
	}

	h.logger.Info("User authenticated successfully", zap.String("userID", user.ID.String()))

	response.SetAuthCookie(w, r, token)
	response.HxRedirect(w, r, "/gen")
}
