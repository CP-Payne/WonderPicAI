package response

import (
	"fmt"
	"net/http"

	creditPages "github.com/CP-Payne/wonderpicai/web/template/pages/credits"
	"github.com/CP-Payne/wonderpicai/web/template/viewmodel"
	"go.uber.org/zap"
)

func LoadPurchasePage(w http.ResponseWriter, r *http.Request, logger *zap.Logger, vm viewmodel.PurchaseViewData) (renderErr error) {
	err := creditPages.PurchasePage(vm).Render(r.Context(), w)
	if err != nil {
		logger.Error("Failed to render purchase page", zap.Error(err))
		return fmt.Errorf("failed to render purchase page: %w", err)
	}
	return nil
}
