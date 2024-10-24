package statistics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStatisticsService is a mock implementation of IStatisticsService
type MockStatisticsService struct {
	mock.Mock
}

// GetStats mocks the behavior of fetching statistics
func (m *MockStatisticsService) GetStats(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{
		"most_upvoted_post": gin.H{
			"author": "mock_author",
			"title":  "mock_title",
			"score":  999,
		},
		"user_post_counts": gin.H{
			"mock_author": 10,
		},
	})
}

func TestGetStats(t *testing.T) {
	router := gin.Default()

	// Initialize the mock service
	mockService := new(MockStatisticsService)

	// Set expectations
	mockService.On("GetStats", mock.Anything).Return()

	// Setup route using the mock
	router.GET("/stats", mockService.GetStats)

	// Create a new HTTP GET request
	req, _ := http.NewRequest("GET", "/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert that the status code is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Check that mock statistics are returned
	assert.Contains(t, w.Body.String(), "mock_author")
	assert.Contains(t, w.Body.String(), "mock_title")
	assert.Contains(t, w.Body.String(), "999")

	// Assert that the mock was called
	mockService.AssertExpectations(t)
}
