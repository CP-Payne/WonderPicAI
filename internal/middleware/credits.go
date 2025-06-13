package middleware

import (
	"net/http"

	"github.com/CP-Payne/wonderpicai/internal/context/auth"
	"github.com/CP-Payne/wonderpicai/internal/context/credits"
	"github.com/CP-Payne/wonderpicai/internal/handler/http/response"
	"github.com/CP-Payne/wonderpicai/internal/service"
	"go.uber.org/zap"
)

func WithCredits(logger *zap.Logger, walletService service.WalletService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			userID, err := auth.UserID(r.Context())
			if err != nil {
				response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
				return
			}

			wallet, err := walletService.GetWallet(r.Context(), userID)
			if err != nil {
				logger.Error("Failed to retrieve user Wallet", zap.String("userID", userID.String()), zap.Error(err))
				response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
				return

			}

			ctx := credits.NewContextWithCredits(r.Context(), int(wallet.Credits))

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
