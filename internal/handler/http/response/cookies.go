package response

import (
	"net/http"
	"time"

	"github.com/CP-Payne/wonderpicai/internal/config"
)

func SetAuthCookie(w http.ResponseWriter, r *http.Request, accessToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		MaxAge:   config.Cfg.JWT.ExpiryMinutes * 60,
		Expires:  time.Now().Add(time.Duration(config.Cfg.JWT.ExpiryMinutes) * time.Minute),
	})
}
