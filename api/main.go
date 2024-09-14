// api/main.go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Initialize database connection
	initDatabase()
	defer dbpool.Close()

	router := gin.Default()

	// Apply middleware and routes
	router.GET("/digests", BasicAuthMiddleware(), getDigests)

	// Start the server on port 8080
	router.Run(":8080")
}
