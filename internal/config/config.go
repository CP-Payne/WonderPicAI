package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	ComfyLite ComfyLiteConfig
	Stripe    StripeConfig
}

type ServerConfig struct {
	AppEnv   string
	Port     string
	LogLevel string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
	DSN      string // Constructed DSN
}

type JWTConfig struct {
	SecretKey     string
	Issuer        string
	ExpiryMinutes int
}

// Image generation server config
type ComfyLiteConfig struct {
	Host string
	Port string
}

type StripeConfig struct {
	Secret string
}

var Cfg AppConfig

func LoadConfig() {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	if appEnv != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: .env file not found or could not be loaded, relying on system env vars:", err)
		}
	} else {
		log.Println("Running in production mode, not loading .env file.")
	}

	// --- Server Config ---
	Cfg.Server.AppEnv = getEnv("APP_ENV", "development")
	Cfg.Server.Port = getEnv("PORT", "8080")
	Cfg.Server.LogLevel = getEnv("LOG_LEVEL", "info")

	// --- Database Config ---
	Cfg.Database.Host = getEnv("DB_HOST", "localhost")
	Cfg.Database.Port = getEnv("DB_PORT", "5432")
	Cfg.Database.User = getEnv("DB_USER", "postgres")
	Cfg.Database.Password = getEnv("DB_PASSWORD", "secret")
	Cfg.Database.DBName = getEnv("DB_NAME", "mydatabase")
	Cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")
	Cfg.Database.TimeZone = getEnv("DB_TIMEZONE", "UTC")
	Cfg.Database.DSN = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		Cfg.Database.Host,
		Cfg.Database.Port,
		Cfg.Database.User,
		Cfg.Database.Password,
		Cfg.Database.DBName,
		Cfg.Database.SSLMode,
		Cfg.Database.TimeZone,
	)

	// --- JWT Config ---
	Cfg.JWT.SecretKey = getEnv("JWT_SECRET_KEY", "") // No sensible default for a secret
	Cfg.JWT.Issuer = getEnv("JWT_ISSUER", "your-app")
	expiryMinutesStr := getEnv("JWT_EXPIRY_MINUTES", "60")
	expiryMinutes, err := strconv.Atoi(expiryMinutesStr)
	if err != nil {
		log.Printf("Warning: Invalid JWT_EXPIRY_MINUTES value '%s', using default 60: %v", expiryMinutesStr, err)
		expiryMinutes = 60
	}
	Cfg.JWT.ExpiryMinutes = expiryMinutes

	// Validate critical configurations
	if Cfg.JWT.SecretKey == "" {
		log.Fatal("FATAL: JWT_SECRET_KEY is not set. Application cannot start.")
	}
	if len(Cfg.JWT.SecretKey) < 32 {
		log.Fatal("FATAL: JWT_SECRET_KEY is too short (must be at least 32 bytes). Application cannot start.")
	}

	// --- ComfyLite ---
	Cfg.ComfyLite.Host = getEnv("COMFYLITE_HOST", "127.0.0.1")
	Cfg.ComfyLite.Port = getEnv("COMFYLITE_PORT", "8081")

	// --- Stripe ---
	Cfg.Stripe.Secret = getEnv("STRIPE_SECRET", "")

	log.Println("Configuration loaded successfully. APP_ENV:", Cfg.Server.AppEnv)

}

// getEnv retrieves the environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if defaultValue == "" && strings.Contains(key, "SECRET") || strings.Contains(key, "PASSWORD") || strings.Contains(key, "KEY") {
		log.Printf("Warning: Environment variable %s is not set and has no default value. This might be an issue.", key)
	}
	return defaultValue
}
