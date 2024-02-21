package config

import (
	"go-crypto-market-api/internal/domain"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() (*gorm.DB, error) {
	var err error
	databaseURL := os.Getenv("DATABASE_URL")
	// Initialize database connection using GORM
	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
		return nil, err
	}

	DB.AutoMigrate(&domain.PriceHistory{})
	return DB, nil
}

func GetDB() *gorm.DB {
	return DB
}
