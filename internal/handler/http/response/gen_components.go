package response

import (
	"fmt"
	"net/http"

	genComponents "github.com/CP-Payne/wonderpicai/web/template/components/gen"
	"github.com/CP-Payne/wonderpicai/web/template/viewmodel"
	"go.uber.org/zap"
)

func LoadGenForm(w http.ResponseWriter, r *http.Request, logger *zap.Logger, vm viewmodel.GenFormComponentData) (renderErr error) {
	err := genComponents.GenForm(vm).Render(r.Context(), w)
	if err != nil {
		logger.Error("Failed to render GenForm component", zap.Error(err))
		return fmt.Errorf("failed to render gen form: %w", err)
	}
	return nil
}

func LoadPendingImage(w http.ResponseWriter, r *http.Request, logger *zap.Logger, vm viewmodel.Image) (renderErr error) {
	err := genComponents.PendingImageCard(vm).Render(r.Context(), w)
	if err != nil {
		logger.Error("Failed to render PendingImageCard component", zap.Error(err))
		return fmt.Errorf("failed to render pending image: %w", err)
	}
	return nil
}
