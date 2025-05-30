package main

import (
	"log"
	"net/http"
	"os"

	gormadapter "github.com/CP-Payne/wonderpicai/internal/adapter/persistence/gorm"
	"github.com/CP-Payne/wonderpicai/internal/app"
	authhandler "github.com/CP-Payne/wonderpicai/internal/handler/http"
	"github.com/CP-Payne/wonderpicai/internal/service"
)

func main() {

	gormadapter.ConnectDatabase()
	db := gormadapter.DB

	userRepo := gormadapter.NewGormUserRepository(db)

	authSvc := service.NewAuthService(userRepo)

	authHndlr := authhandler.NewAuthHandler(authSvc)

	router := app.NewRouter(authHndlr)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
