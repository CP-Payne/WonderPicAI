package main

import (
	"log"
	"net/http"
	"os"

	gormadapter "github.com/CP-Payne/wonderpicai/internal/adapter/persistence/gorm"
	"github.com/CP-Payne/wonderpicai/internal/app"
	appconfig "github.com/CP-Payne/wonderpicai/internal/config"
	authhandler "github.com/CP-Payne/wonderpicai/internal/handler/http"
	"github.com/CP-Payne/wonderpicai/internal/service"
)

func main() {

	appconfig.LoadConfig()
	cfg := appconfig.Cfg

	gormadapter.ConnectDatabase(cfg.Database.DSN)
	db := gormadapter.DB

	// tokenProvider, err := jwtadapter.NewJWTTokenProvider(
	// 	cfg.JWT.SecretKey,
	// 	cfg.JWT.Issuer,
	// 	cfg.JWT.ExpiryMinutes,
	// ) // MODIFIED
	// if err != nil {
	// 	log.Fatalf("Failed to initialize token provider: %v", err)
	// }

	userRepo := gormadapter.NewGormUserRepository(db)

	authSvc := service.NewAuthService(userRepo)

	authHndlr := authhandler.NewAuthHandler(authSvc)

	router := app.NewRouter(authHndlr)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on http://localhost:%s\n", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
