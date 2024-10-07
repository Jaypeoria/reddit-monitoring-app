package oauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetOAuthToken(t *testing.T) {
	router := gin.Default()
	router.GET("/auth", GetOAuthToken)

	req, _ := http.NewRequest("GET", "/auth", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", w.Code)
	}

	// Check if the response contains the expected fields
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil || response["access_token"] == "" {
		t.Errorf("Expected a valid access token, got %v", response)
	}
}
