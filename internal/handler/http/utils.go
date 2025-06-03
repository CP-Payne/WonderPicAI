package http

import "net/http"

func HxRedirect(w http.ResponseWriter, r *http.Request, to string) {
	// if len(r.Header.Get("HX-Request")) > 0 {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", to)
		w.WriteHeader(http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, to, http.StatusSeeOther)
}
