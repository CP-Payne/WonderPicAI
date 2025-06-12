package middleware

import (
	"errors"
	"net/http"

	"github.com/CP-Payne/wonderpicai/internal/context/auth"
	"github.com/CP-Payne/wonderpicai/internal/handler/http/response"
	"github.com/CP-Payne/wonderpicai/internal/port"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func WithAuth(logger *zap.Logger, tokenService port.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			cookie, err := r.Cookie("auth_token")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					response.HxRedirect(w, r, "/auth/login")
					return
				} else {
					response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
					return
				}
			}

			jwtCookie := cookie.Value

			token, err := tokenService.ValidateToken(jwtCookie)
			if err != nil {
				logger.Warn("Token failed verification", zap.Error(err))
				response.HxRedirect(w, r, "/auth/login")
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				// This can happen if the token is valid but the claims are not what you expect.
				logger.Error("Invalid token claims type", zap.String("token", jwtCookie))
				response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
				return
			}

			subClaim, ok := claims["sub"].(string)
			if !ok {
				logger.Warn("Token missing or has malformed 'sub' claim", zap.Any("claims", claims))
				response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
				return
			}

			userID, err := uuid.Parse(subClaim)
			if err != nil {
				logger.Warn("Invalid userID UUID in claims", zap.String("sub", subClaim), zap.Error(err))
				response.HxRedirectErrorPage(w, r, http.StatusInternalServerError, "", "")
				return
			}

			ctx := auth.NewContextWithUserID(r.Context(), userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
