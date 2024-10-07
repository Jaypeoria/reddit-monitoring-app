package statistics

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetStats(t *testing.T) {
	router := gin.Default()

	statisticsService := NewStatisticsService()
	router.GET("/stats", statisticsService.GetStats)

	req, _ := http.NewRequest("GET", "/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", w.Code)
	}

	// Check if the response contains the expected fields
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil || response["most_upvoted_post"] == nil || response["user_post_counts"] == nil {
		t.Errorf("Expected statistics data, got %v", response)
	}
}
