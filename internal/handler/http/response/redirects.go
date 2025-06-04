package response

import (
	"fmt"
	"net/http"
)

// HxRedirect adds the appropriate HTMX redirect headers
// It performs a normal http redirect if the HX-Request header is not true.
func HxRedirect(w http.ResponseWriter, r *http.Request, to string) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", to)
		w.WriteHeader(http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, to, http.StatusSeeOther)
}

// HxRedirectErrorPage prepares the query parameters for the error page.
// It calls HxRedirect with the `to` parameter set to the the error page with the prepared query.
func HxRedirectErrorPage(w http.ResponseWriter, r *http.Request, statusCode int, errorID, message string) {
	redirectString := fmt.Sprintf("/error?statusCode=%d&errorID=%s&message=%s", http.StatusInternalServerError, errorID, message)

	HxRedirect(w, r, redirectString)
}
