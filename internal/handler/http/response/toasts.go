package response

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/CP-Payne/wonderpicai/web/template/components/ui"
	"github.com/CP-Payne/wonderpicai/web/template/viewmodel"
	"github.com/rs/xid"
)

// LoadErrorToast prepares and writes (render) the error toast component to the response.
// Returns the toastID and any error encountered during rendering.
// The caller should check the error and typically return if it is not nil.
func LoadErrorToast(w http.ResponseWriter, r *http.Request, logger *zap.Logger, message string) (toastID string, loadErr error) {
	toastData := viewmodel.ToastComponentData{
		Message: message,
		Type:    viewmodel.ToastError,
		ToastID: xid.New().String(),
	}

	err := ui.ToastNotification(toastData).Render(r.Context(), w)
	if err != nil {
		logger.Error("Failed to render toast component", zap.Error(err), zap.String("toastID", toastData.ToastID))
		return toastData.ToastID, err
	}

	return toastData.ToastID, nil
}
