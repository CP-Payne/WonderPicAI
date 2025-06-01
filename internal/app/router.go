package app

import (
	"net/http"

	allHandlers "github.com/CP-Payne/wonderpicai/internal/handler/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handlers *allHandlers.ApiHandlers) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	staticServer := StaticFSHandler()
	r.Mount("/static/", staticServer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Get("/", handlers.LandingHandler.ShowLandingPage)

	r.Route("/auth", func(r chi.Router) {
		r.Get("/login", handlers.AuthHandler.ShowLoginPage)
		r.Get("/signup", handlers.AuthHandler.ShowSignupPage)

		r.Post("/signup", handlers.AuthHandler.HandleSignup)
		r.Post("/login", handlers.AuthHandler.HandleLogin)
	})

	return r
}
