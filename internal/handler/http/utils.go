package http

import (
	"fmt"
	"net/http"
)

func HxRedirect(w http.ResponseWriter, r *http.Request, to string) {
	// if len(r.Header.Get("HX-Request")) > 0 {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", to)
		w.WriteHeader(http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, to, http.StatusSeeOther)
}

func HxRedirectErrorPage(w http.ResponseWriter, r *http.Request, statusCode int, errorID, message string) {
	redirectString := fmt.Sprintf("/error?statusCode=%d&errorID=%s&message=%s", http.StatusInternalServerError, errorID, message)

	HxRedirect(w, r, redirectString)
}
