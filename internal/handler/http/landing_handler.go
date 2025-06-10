package http

import (
	"net/http"

	"github.com/CP-Payne/wonderpicai/web/template/pages/landing"
	"go.uber.org/zap"
)

type LandingHandler struct {
	logger *zap.Logger
}

func NewLandingHandler(logger *zap.Logger) *LandingHandler {
	return &LandingHandler{
		logger: logger,
	}
}

func (h *LandingHandler) ShowLandingPage(w http.ResponseWriter, r *http.Request) {
	err := landing.LandingPage().Render(r.Context(), w)
	if err != nil {
		h.logger.Error("Failed to render landing page", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
