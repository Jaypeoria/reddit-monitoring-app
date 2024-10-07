package oauth

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// GetOAuthToken provides Reddit OAuth token
func GetOAuthToken(c *gin.Context) {
	clientID := os.Getenv("REDDIT_CLIENT_ID")
	clientSecret := os.Getenv("REDDIT_CLIENT_SECRET")

	// Prepare the Reddit OAuth token request
	reqBody := strings.NewReader("grant_type=client_credentials")
	req, err := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token", reqBody)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Set basic authentication header
	req.SetBasicAuth(clientID, clientSecret)
	// Set content-type to x-www-form-urlencoded
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error getting OAuth token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get token"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-200 response from Reddit API: %v", resp.Status)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve token", "status": resp.Status})
		return
	}

	// Parse the token response
	var tokenResp TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		log.Printf("Error decoding token response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode token response"})
		return
	}

	// Return the access token to the caller
	c.JSON(http.StatusOK, gin.H{"access_token": tokenResp.AccessToken})
}
