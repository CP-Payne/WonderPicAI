package gorm

import (
	"log"
	"os"

	"github.com/CP-Payne/wonderpicai/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase(dsn string) {

	if dsn == "" {
		log.Fatal("Database DSN cannot be empty")
		os.Exit(1)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v\nDSN Used: %s", err, dsn)
		os.Exit(1)
	}

	log.Println("Database connection established.")

	err = DB.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}
	log.Println("Database schema migrated")
}

func GetDB() *gorm.DB {
	return DB
}
