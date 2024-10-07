package fetcher

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	//"jackhenry.com/reddit-fetcher-service/db"

	"github.com/gin-gonic/gin"
	"jackhenry.com/reddit-fetcher-service/db"
)

type Post struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Score  int    `json:"score"`
}

type FetcherService struct {
	Posts      []Post
	mu         sync.Mutex
	oauthToken string
}

func NewFetcherService() *FetcherService {
	return &FetcherService{
		Posts: make([]Post, 0),
	}
}

// StartFetching fetches posts from a subreddit
func (f *FetcherService) StartFetching(subreddit string) {
	delay := 2 * time.Second
	f.refreshOAuthToken()

	for {
		delay = f.fetchPosts(subreddit, delay)
		time.Sleep(delay)
	}
}

func (f *FetcherService) fetchPosts(subreddit string, delay time.Duration) time.Duration {
	url := "https://oauth.reddit.com/r/" + subreddit + "/new.json"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+f.oauthToken)
	req.Header.Add("User-Agent", "reddit-fetcher-service")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		return delay
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		log.Printf("OAuth token expired, refreshing token")
		f.refreshOAuthToken()
		return delay
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: received non-200 status code: %v", resp.StatusCode)
		return delay
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
		return delay
	}

	f.mu.Lock()
	for _, child := range response.Data.Children {
		post := child.Data
		f.Posts = append(f.Posts, post)
		log.Printf("Post: %v", post)
		db.InsertPost(post)
	}
	f.mu.Unlock()

	_, remaining, reset := f.extractRateLimitHeaders(resp)
	if remaining <= 1 {
		delay = time.Duration(reset+1) * time.Second
		log.Printf("Rate limit reached, waiting for %d seconds to reset", reset)
	} else {
		delay = 2 * time.Second
		log.Printf("Rate limit remaining: %d, adjusting delay to %v", remaining, delay)
	}

	return delay
}

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

func (f *FetcherService) extractRateLimitHeaders(resp *http.Response) (used, remaining, reset int) {
	used, _ = strconv.Atoi(resp.Header.Get("X-Ratelimit-Used"))
	remaining, _ = strconv.Atoi(resp.Header.Get("X-Ratelimit-Remaining"))
	reset, _ = strconv.Atoi(resp.Header.Get("X-Ratelimit-Reset"))
	log.Printf("Rate Limit: Used: %d, Remaining: %d, Reset: %d", used, remaining, reset)
	return used, remaining, reset
}

func (f *FetcherService) GetPosts(c *gin.Context) {
	posts, err := db.FindAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}
