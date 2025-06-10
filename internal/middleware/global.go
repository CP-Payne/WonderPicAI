package middleware

import (
	"net/http"

	httpHandler "github.com/CP-Payne/wonderpicai/internal/handler/http"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

func CustomRecoverer(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
					errorID := xid.New().String()
					logger.Error("Panic recovered", zap.Any("panic", rvr), zap.String("errorID", errorID), zap.Stack("stack"))

					message := "An unexpected problem occured. We are looking into it."

					httpHandler.HxRedirectErrorPage(w, r, http.StatusInternalServerError, errorID, message)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
