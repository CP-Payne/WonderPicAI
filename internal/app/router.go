package app

import (
	"net/http"

	allHandlers "github.com/CP-Payne/wonderpicai/internal/handler/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func NewRouter(handlers *allHandlers.ApiHandlers, logger *zap.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(CustomRecoverer(logger))
	// r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	staticServer := StaticFSHandler()
	r.Mount("/static/", staticServer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Get("/", handlers.LandingHandler.ShowLandingPage)
	r.Get("/error", handlers.ErrorHandler.ServeGenericErrorPage)

	r.Get("/gen", handlers.GenHandler.ShowGenPage)
	r.Post("/gen", handlers.GenHandler.HandleGenerationCreate)

	r.Route("/auth", func(r chi.Router) {
		r.Get("/login", handlers.AuthHandler.ShowLoginPage)
		r.Get("/signup", handlers.AuthHandler.ShowSignupPage)

		r.Post("/signup", handlers.AuthHandler.HandleSignup)
		r.Post("/login", handlers.AuthHandler.HandleLogin)
	})

	return r
}
