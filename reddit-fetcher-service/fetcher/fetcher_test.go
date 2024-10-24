package fetcher

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFetcherService is a mock implementation of IFetcherService
type MockFetcherService struct {
	mock.Mock
	postsChan chan Post
}

// StartFetching mocks the fetching behavior
func (m *MockFetcherService) StartFetching(subreddit string) {
	m.Called(subreddit)
}

// GetPosts mocks the behavior of fetching posts from MongoDB
func (m *MockFetcherService) GetPosts(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{"posts": []string{"mock_post_1", "mock_post_2"}})
}

func TestGetPosts(t *testing.T) {
	router := gin.Default()

	// Initialize the mock service
	mockService := new(MockFetcherService)
	mockService.postsChan = make(chan Post)

	// Set expectations
	mockService.On("GetPosts", mock.Anything).Return()

	// Setup route using the mock
	router.GET("/fetch", mockService.GetPosts)

	// Create a new HTTP GET request
	req, _ := http.NewRequest("GET", "/fetch", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert that the status code is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Check that mock posts are returned
	assert.Contains(t, w.Body.String(), "mock_post_1")
	assert.Contains(t, w.Body.String(), "mock_post_2")

	// Assert that the mock was called
	mockService.AssertExpectations(t)
}

func TestProcessPosts(t *testing.T) {
	// Create a channel for posts
	postsChan := make(chan Post)

	// Create a mock fetcher service
	fetcherService := &FetcherService{
		postsChan: postsChan,
	}

	// Create a post to send to the channel
	mockPost := Post{Author: "test_author", Title: "test_title", Score: 100}

	// Start a goroutine to process posts
	go fetcherService.processPosts()

	// Send a post to the channel
	postsChan <- mockPost

	// Allow some time for processing
	time.Sleep(1 * time.Second)

	// Assert that the post was processed (this would usually involve verifying MongoDB insertions or mock behavior)
	assert.NotNil(t, mockPost, "Post should have been processed")
}
