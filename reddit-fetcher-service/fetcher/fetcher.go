package fetcher

import (
	"encoding/json"
	"log"
	"net/http"

	"time"

	"jackhenry.com/reddit-fetcher-service/db"

	"github.com/gin-gonic/gin"
)

// IFetcherService defines an interface for fetching Reddit posts
type IFetcherService interface {
	StartFetching(subreddit string)
	GetPosts(c *gin.Context)
}

// FetcherService implements the IFetcherService interface
type FetcherService struct {
	postsChan  chan Post // Channel to handle posts
	oauthToken string
}

// Post represents a Reddit post
type Post struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Score  int    `json:"score"`
}

// NewFetcherService creates a new instance of FetcherService
func NewFetcherService() IFetcherService {
	return &FetcherService{
		postsChan: make(chan Post), // Initialize the channel
	}
}

// StartFetching fetches posts from the specified subreddit and stores them in MongoDB
func (f *FetcherService) StartFetching(subreddit string) {
	// Start a goroutine to listen on the posts channel and process incoming posts
	go f.processPosts()

	// Start fetching posts
	f.refreshOAuthToken()
	for {
		f.fetchPosts(subreddit)
		time.Sleep(2 * time.Second)
	}
}

// fetchPosts fetches posts from Reddit API and sends them to the posts channel
func (f *FetcherService) fetchPosts(subreddit string) {
	url := "https://oauth.reddit.com/r/" + subreddit + "/new.json"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+f.oauthToken)
	req.Header.Add("User-Agent", "reddit-fetcher-service")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		log.Printf("OAuth token expired, refreshing token")
		f.refreshOAuthToken()
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: received non-200 status code: %v", resp.StatusCode)
		return
	}

	var response struct {
		Data struct {
			Children []struct {
				Data Post `json:"data"`
			} `json:"children"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("Error decoding response: %v", err)
		return
	}

	// Send the fetched posts to the posts channel for processing
	for _, child := range response.Data.Children {
		f.postsChan <- child.Data
	}
}

// processPosts listens on the posts channel and processes each post
func (f *FetcherService) processPosts() {
	for post := range f.postsChan {
		// Insert the post into MongoDB
		db.InsertPost(post)
		log.Printf("Processed post by %s with title '%s'", post.Author, post.Title)
	}
}

// refreshOAuthToken refreshes the OAuth token by calling the OAuth service
func (f *FetcherService) refreshOAuthToken() {
	resp, err := http.Get("http://reddit-oauth-service:8081/auth")
	if err != nil {
		log.Fatalf("Failed to get OAuth token: %v", err)
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		log.Fatalf("Error decoding OAuth token response: %v", err)
	}

	f.oauthToken = tokenResp.AccessToken
}

// GetPosts retrieves all posts stored in MongoDB and returns them via the API
func (f *FetcherService) GetPosts(c *gin.Context) {
	posts, err := db.FindAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}
