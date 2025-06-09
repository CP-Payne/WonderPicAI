package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	genPages "github.com/CP-Payne/wonderpicai/web/template/pages/gen"
	"github.com/CP-Payne/wonderpicai/web/template/viewmodel"
)

type GenHandler struct {
	logger   *zap.Logger
	validate *validator.Validate
}

func NewGenHandler(logger *zap.Logger, validate *validator.Validate) *GenHandler {
	return &GenHandler{logger: logger.With(zap.String("component", "GenHandler")), validate: validate}
}

func (h *GenHandler) ShowGenPage(w http.ResponseWriter, r *http.Request) {
	// TODO: Get images from database
	genPageData := viewmodel.GenPageData{
		GalleryData: viewmodel.GalleryComponentData{
			Images: []viewmodel.Image{},
		},
		GenFormData: viewmodel.GenFormComponentData{
			Form:   viewmodel.GenFormData{},
			Errors: map[string]string{},
			Error:  "",
		},
	}
	err := genPages.GenPage(genPageData).Render(r.Context(), w)
	if err != nil {
		h.logger.Error("Failed to render login page", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
