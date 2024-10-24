package statistics

import (
	"net/http"
	"sync"

	"jackhenry.com/reddit-statistics-service/db"

	"github.com/gin-gonic/gin"
)

// IStatisticsService defines an interface for statistics generation
type IStatisticsService interface {
	GetStats(c *gin.Context)
}

// StatisticsService implements the IStatisticsService interface
type StatisticsService struct {
	mu              sync.Mutex
	mostUpvotedPost map[string]interface{}
	userPostCounts  map[string]int
}

// NewStatisticsService creates a new instance of StatisticsService
func NewStatisticsService() IStatisticsService {
	return &StatisticsService{
		userPostCounts: make(map[string]int),
	}
}

// GetStats retrieves statistics from MongoDB and returns them via the API
func (s *StatisticsService) GetStats(c *gin.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	posts, err := db.FindAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	for _, post := range posts {
		author := post["author"].(string)
		score := post["score"].(int32)

		s.userPostCounts[author]++
		if s.mostUpvotedPost == nil || score > s.mostUpvotedPost["score"].(int32) {
			s.mostUpvotedPost = post
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"most_upvoted_post": s.mostUpvotedPost,
		"user_post_counts":  s.userPostCounts,
	})
}
