package main

import (
	"log"
	"net/http"

	gormadapter "github.com/CP-Payne/wonderpicai/internal/adapter/persistence/gorm"
	"github.com/CP-Payne/wonderpicai/internal/app"
	appconfig "github.com/CP-Payne/wonderpicai/internal/config"
	allHandlers "github.com/CP-Payne/wonderpicai/internal/handler/http"
	applogger "github.com/CP-Payne/wonderpicai/internal/logger"
	"github.com/CP-Payne/wonderpicai/internal/service"
	"go.uber.org/zap"
)

func main() {

	appconfig.LoadConfig()
	cfg := appconfig.Cfg

	logger, err := applogger.New(cfg.Server.LogLevel, cfg.Server.AppEnv)
	if err != nil {
		log.Fatalf("Failed to initialize application logger: %v", err)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	gormadapter.ConnectDatabase(cfg.Database.DSN, cfg.Server.AppEnv, cfg.Server.LogLevel, logger)
	db := gormadapter.DB

	// tokenProvider, err := jwtadapter.NewJWTTokenProvider(
	// 	cfg.JWT.SecretKey,
	// 	cfg.JWT.Issuer,
	// 	cfg.JWT.ExpiryMinutes,
	// logger
	// ) // MODIFIED
	// if err != nil {
	// 	log.Fatalf("Failed to initialize token provider: %v", err)
	// }

	userRepo := gormadapter.NewGormUserRepository(db)

	authSvc := service.NewAuthService(userRepo, logger)

	apiHandlers := allHandlers.NewApiHandlers(authSvc, logger)

	router := app.NewRouter(apiHandlers)

	logger.Info("Server starting",
		zap.String("address", "http://localhost:"+cfg.Server.Port),
		zap.String("app_env", cfg.Server.AppEnv),
	)

	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
