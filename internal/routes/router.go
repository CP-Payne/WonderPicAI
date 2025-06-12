package routes

import (
	"net/http"

	allHandlers "github.com/CP-Payne/wonderpicai/internal/handler/http"
	"github.com/CP-Payne/wonderpicai/internal/middleware"
	"github.com/CP-Payne/wonderpicai/internal/port"
	"github.com/CP-Payne/wonderpicai/internal/service"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func NewRouter(handlers *allHandlers.ApiHandlers, logger *zap.Logger, tokenService port.TokenService, walletService service.WalletService) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(middleware.CustomRecoverer(logger))
	// r.Use(middleware.Recoverer)
	r.Use(chimiddleware.StripSlashes)

	staticServer := StaticFSHandler()
	r.Mount("/static/", staticServer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Get("/", handlers.LandingHandler.ShowLandingPage)
	r.Get("/error", handlers.ErrorHandler.ServeGenericErrorPage)

	r.Post("/gen/update", handlers.GenHandler.HandleImageCompletionWebhook)

	r.Route("/gen", func(r chi.Router) {
		r.Use(middleware.WithAuth(logger, tokenService))

		r.Group(func(r chi.Router) {
			r.Use(middleware.WithCredits(logger, walletService))
			r.Get("/", handlers.GenHandler.ShowGenPage)
			r.Post("/", handlers.GenHandler.HandleGenerationCreate)
		})

		r.Get("/image/{id}/status", handlers.GenHandler.HandleImageStatus)
		r.Delete("/image/{id}", handlers.GenHandler.HandleImageDelete)
		r.Delete("/image/failed", handlers.GenHandler.HandleFailedImagesDelete)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Get("/login", handlers.AuthHandler.ShowLoginPage)
		r.Get("/signup", handlers.AuthHandler.ShowSignupPage)

		r.Post("/logout", handlers.AuthHandler.HandleLogout)
		r.Post("/signup", handlers.AuthHandler.HandleSignup)
		r.Post("/login", handlers.AuthHandler.HandleLogin)
	})

	return r
}
