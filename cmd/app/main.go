package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CP-Payne/wonderpicai/internal/adapter/externalauth/googleprovider"
	"github.com/CP-Payne/wonderpicai/internal/adapter/generation/comfylite"
	"github.com/CP-Payne/wonderpicai/internal/adapter/paymentprovider/stripe"
	gormadapter "github.com/CP-Payne/wonderpicai/internal/adapter/persistence/gorm"
	"github.com/CP-Payne/wonderpicai/internal/adapter/tokenservice"
	appconfig "github.com/CP-Payne/wonderpicai/internal/config"
	allHandlers "github.com/CP-Payne/wonderpicai/internal/handler/http"
	applogger "github.com/CP-Payne/wonderpicai/internal/logger"
	"github.com/CP-Payne/wonderpicai/internal/routes"
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

	tokenService := tokenservice.NewTokenService(cfg.JWT.SecretKey, cfg.JWT.Issuer)
	genClient := comfylite.NewClient(logger, fmt.Sprintf("http://%s:%s", cfg.ComfyLite.Host, cfg.ComfyLite.Port), fmt.Sprintf("http://localhost:%s/gen/update", cfg.Server.Port))

	baseURL := "http://localhost:" + cfg.Server.Port
	successURL := baseURL + "/purchase/success"
	cancelURL := baseURL + "/purchase/cancel"

	stripeProvider := stripe.NewProvider(logger, cfg.Stripe.Secret, cfg.Stripe.VerificationSecret, successURL, cancelURL)
	googleAuthProvider := googleprovider.NewAuth(logger, cfg.GoogleAuth.ClientSecret)

	userRepo := gormadapter.NewGormUserRepository(db, logger)
	promptRepo := gormadapter.NewGormPromptRepository(db, logger)
	imageRepo := gormadapter.NewGormImageRepository(db, logger)
	walletRepo := gormadapter.NewGormWalletRepository(db, logger)

	walletSvc := service.NewWalletService(logger, walletRepo)
	authSvc := service.NewAuthService(userRepo, tokenService, logger, googleAuthProvider)
	genSvc := service.NewGenService(logger, genClient, promptRepo, imageRepo, walletSvc)
	purchaseSvc := service.NewPurchaseService(logger, walletSvc, stripeProvider, userRepo)

	apiHandlers := allHandlers.NewApiHandlers(authSvc, genSvc, purchaseSvc, logger)

	router := routes.NewRouter(apiHandlers, logger, tokenService, walletSvc)

	logger.Info("Server starting",
		zap.String("address", "http://0.0.0.0:"+cfg.Server.Port),
		zap.String("app_env", cfg.Server.AppEnv),
	)

	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
