package main

import (
	"log"
	"os"

	"jackhenry.com/reddit-fetcher-service/db"
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
	db.Connect() // Call MongoDB connection

	r := gin.Default()
	fetcherService := fetcher.NewFetcherService()
	go fetcherService.StartFetching("golang")

	r.GET("/fetch", fetcherService.GetPosts)

	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatalf("Failed to start Fetcher service: %v", err)
	}
}
