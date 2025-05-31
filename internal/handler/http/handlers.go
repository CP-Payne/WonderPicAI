package http

import (
	"github.com/CP-Payne/wonderpicai/internal/service"
	"go.uber.org/zap"
)

type ApiHandlers struct {
	AuthHandler *AuthHandler
	PageHandler *PageHandler
}

func NewApiHandlers(authService service.AuthService, logger *zap.Logger) *ApiHandlers {
	return &ApiHandlers{
		AuthHandler: NewAuthHandler(authService, logger),
		PageHandler: NewPageHandler(),
	}
}
