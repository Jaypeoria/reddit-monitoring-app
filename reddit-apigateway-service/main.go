package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/posts", func(c *gin.Context) {
		resp, err := http.Get("http://reddit-fetcher-service:8082/fetch")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
			return
		}
		defer resp.Body.Close()

		var posts interface{}
		if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode posts"})
			return
		}

		c.JSON(http.StatusOK, posts)
	})

	r.GET("/stats", func(c *gin.Context) {
		resp, err := http.Get("http://reddit-statistics-service:8083/stats")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch statistics"})
			return
		}
		defer resp.Body.Close()

		var stats interface{}
		if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode statistics"})
			return
		}

		c.JSON(http.StatusOK, stats)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}
