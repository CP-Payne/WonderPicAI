package response

import (
	"fmt"
	"net/http"

	"github.com/CP-Payne/wonderpicai/web/template/viewmodel"
	"go.uber.org/zap"

	authComponents "github.com/CP-Payne/wonderpicai/web/template/components/auth"
)

// LoadLoginForm prepares and writes (render) the LoginForm component, typically with validation errors.
// Returns any error encountered during rendering.
// The caller should check the error and typically return if it is not nil.
func LoadLoginForm(w http.ResponseWriter, r *http.Request, logger *zap.Logger, vm viewmodel.LoginFormComponentData) (renderErr error) {
	err := authComponents.LoginForm(vm).Render(r.Context(), w)
	if err != nil {
		logger.Error("Failed to render LoginForm component", zap.Error(err))
		return fmt.Errorf("failed to render login form: %w", err)
	}
	return nil
}

// LoadSignupForm prepares and writes (render) the SignupForm component, typically with validation errors.
// Returns any error encountered during rendering.
// The caller should check the error and typically return if it is not nil.
func LoadSignupForm(w http.ResponseWriter, r *http.Request, logger *zap.Logger, vm viewmodel.SignupFormComponentData) (renderErr error) {
	err := authComponents.SignupForm(vm).Render(r.Context(), w)
	if err != nil {
		logger.Error("Failed to render SignupForm component", zap.Error(err))
		return fmt.Errorf("failed to render signup form: %w", err)
	}
	return nil
}
