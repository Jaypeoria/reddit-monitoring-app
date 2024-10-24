package oauth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOAuthService is a mock implementation of IOAuthService
type MockOAuthService struct {
	mock.Mock
}

// GetOAuthToken mocks the behavior of OAuthService
func (m *MockOAuthService) GetOAuthToken(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{"access_token": "mock_token"})
}

func TestGetOAuthToken(t *testing.T) {
	router := gin.Default()

	// Initialize the mock service
	mockService := new(MockOAuthService)

	// Set expectations
	mockService.On("GetOAuthToken", mock.Anything).Return()

	// Setup route using the mock
	router.GET("/auth", mockService.GetOAuthToken)

	// Create a new HTTP GET request
	req, _ := http.NewRequest("GET", "/auth", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert that the status code is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Check that the access token is returned
	assert.Contains(t, w.Body.String(), "mock_token")

	// Assert that the mock was called
	mockService.AssertExpectations(t)
}
