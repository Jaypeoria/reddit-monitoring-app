package main

import (
	"log"
	"os"

	"jackhenry.com/reddit-statistics-service/db"
	"jackhenry.com/reddit-statistics-service/statistics"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file for local development
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	// Connect to MongoDB
	db.Connect()

	// Create an instance of the Statistics service
	var statisticsService statistics.IStatisticsService = statistics.NewStatisticsService()

	r := gin.Default()
	r.GET("/stats", statisticsService.GetStats)

	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatalf("Failed to start Statistics service: %v", err)
	}
}
