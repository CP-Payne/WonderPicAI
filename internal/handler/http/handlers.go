package http

import (
	"github.com/CP-Payne/wonderpicai/internal/service"
	"go.uber.org/zap"
)

type ApiHandlers struct {
	AuthHandler    *AuthHandler
	LandingHandler *LandingHandler
}

func NewApiHandlers(authService service.AuthService, logger *zap.Logger) *ApiHandlers {
	return &ApiHandlers{
		AuthHandler:    NewAuthHandler(authService, logger),
		LandingHandler: NewLandingHandler(logger),
	}
}
