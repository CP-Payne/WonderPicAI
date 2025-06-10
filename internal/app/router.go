package app

import (
	"net/http"

	authhandler "github.com/CP-Payne/wonderpicai/internal/handler/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(authHandler *authhandler.AuthHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	staticServer := StaticFSHandler()
	r.Mount("/static/", staticServer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	pageHndlr := authhandler.NewPageHandler()

	r.Get("/", pageHndlr.ServeHomePage)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	return r
}
