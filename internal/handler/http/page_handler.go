package http

import (
	"log"
	"net/http"

	"github.com/CP-Payne/wonderpicai/web/template"
)

type PageHandler struct {
}

func NewPageHandler() *PageHandler {
	return &PageHandler{}
}

func (h *PageHandler) ServeHomePage(w http.ResponseWriter, r *http.Request) {
	// Render templ component

	component := template.IndexPage("wonderpicai")
	err := component.Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to render index page template: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
