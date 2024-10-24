package main

import (
	"log"
	"os"

	"jackhenry.com/reddit-oauth-service/oauth"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file for local development
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	// Create an instance of the OAuth service
	var oauthService oauth.IOAuthService = &oauth.OAuthService{}

	r := gin.Default()
	r.GET("/auth", oauthService.GetOAuthToken)

	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatalf("Failed to start OAuth service: %v", err)
	}
}
