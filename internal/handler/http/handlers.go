package http

import (
	"github.com/CP-Payne/wonderpicai/internal/service"
	"github.com/CP-Payne/wonderpicai/internal/validation"
	"go.uber.org/zap"
)

type ApiHandlers struct {
	AuthHandler     *AuthHandler
	LandingHandler  *LandingHandler
	ErrorHandler    *ErrorHandler
	GenHandler      *GenHandler
	PurchaseHandler *PurchaseHandler
}

func NewApiHandlers(authService service.AuthService, genService service.GenService, purchaseService service.PurcaseService, logger *zap.Logger) *ApiHandlers {

	appValidator := validation.New()

	return &ApiHandlers{
		AuthHandler:     NewAuthHandler(authService, logger, appValidator),
		LandingHandler:  NewLandingHandler(logger),
		ErrorHandler:    NewErrorHandler(logger),
		GenHandler:      NewGenHandler(logger, appValidator, genService),
		PurchaseHandler: NewPurchaseHandler(logger, appValidator, purchaseService),
	}
}
