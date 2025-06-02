package http

import (
	"github.com/CP-Payne/wonderpicai/internal/service"
	"github.com/CP-Payne/wonderpicai/internal/validation"
	"go.uber.org/zap"
)

type ApiHandlers struct {
	AuthHandler    *AuthHandler
	LandingHandler *LandingHandler
}

func NewApiHandlers(authService service.AuthService, logger *zap.Logger) *ApiHandlers {

	appValidator := validation.New()

	return &ApiHandlers{
		AuthHandler:    NewAuthHandler(authService, logger, appValidator),
		LandingHandler: NewLandingHandler(logger),
	}
}
