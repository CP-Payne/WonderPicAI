package gorm

import (
	"log"
	"os"
	"time"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase(dsn, appEnv, logLevel string, appLogger *zap.Logger) {

	if dsn == "" {
		appLogger.Fatal("Database DSN cannot be empty")
		os.Exit(1)
	}

	newGormLogger := InitializeGormLogger(appEnv, logLevel)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newGormLogger,
	})

	if err != nil {
		appLogger.Fatal("Failed to connect to database", zap.String("DSN", dsn), zap.Error(err))
		os.Exit(1)
	}

	appLogger.Info("Database connection established.")

	if appEnv != "production" {

		err = DB.Migrator().DropTable(&domain.User{})
		if err != nil {
			appLogger.Error("Failed to drop table users", zap.Error(err))
		}

		err = DB.Migrator().DropTable(&domain.Prompt{})
		if err != nil {
			appLogger.Error("Failed to drop table prompts", zap.Error(err))
		}
		err = DB.Migrator().DropTable(&domain.Image{})
		if err != nil {
			appLogger.Error("Failed to drop table images", zap.Error(err))
		}
	}

	err = DB.AutoMigrate(&domain.User{})
	if err != nil {
		appLogger.Fatal("Failed to auto-migrate database schema", zap.Error(err))
	}

	err = DB.AutoMigrate(&domain.Prompt{})
	if err != nil {
		appLogger.Fatal("Failed to auto-migrate database schema", zap.Error(err))
	}

	err = DB.AutoMigrate(&domain.Image{})
	if err != nil {
		appLogger.Fatal("Failed to auto-migrate database schema", zap.Error(err))
	}

	appLogger.Info("Database schema migrated")
}

func GetDB() *gorm.DB {
	return DB
}

func InitializeGormLogger(appEnv, logLevel string) gormlogger.Interface {
	gormLogLevel := gormlogger.Warn
	if appEnv != "production" && logLevel == "debug" {
		gormLogLevel = gormlogger.Info
	}

	newGormLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormLogLevel,
			IgnoreRecordNotFoundError: false,
			Colorful:                  (appEnv != "production"),
		},
	)

	return newGormLogger
}
