package main

import (
	"log"
	"os"

	"jackhenry.com/reddit-fetcher-service/db" // MongoDB helper
	"jackhenry.com/reddit-fetcher-service/fetcher"

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

	// Create an instance of the Fetcher service
	var fetcherService fetcher.IFetcherService = fetcher.NewFetcherService()

	r := gin.Default()
	go fetcherService.StartFetching("golang") // Change this to your subreddit

	r.GET("/fetch", fetcherService.GetPosts)

	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatalf("Failed to start Fetcher service: %v", err)
	}
}
