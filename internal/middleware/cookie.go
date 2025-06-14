package middleware

import (
	"net/http"

	"github.com/CP-Payne/wonderpicai/internal/handler/http/response"
)

func RedirectIfAuthCookie(targetPath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			if targetPath == "" {
				next.ServeHTTP(w, r)
				return
			}

			cookie, err := r.Cookie("auth_token")
			if err == nil && cookie != nil && cookie.Value != "" {
				response.HxRedirect(w, r, targetPath)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
