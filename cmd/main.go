package main

import (
	"log"
	"os"

	"go-crypto-market-api/internal/config"
	"go-crypto-market-api/internal/interfaces/routers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("PORT")

	config.InitDatabase()
	config.InitRedis()

	router := routers.SetupRouter()

	// Specify the port to run the HTTP server on
	if port == "" {
		port = "8080" // Default port
	}

	// Start the server
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
