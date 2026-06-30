package repository

import (
	"io"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)

// PodcastRepository handles data access for podcast feeds
type PodcastRepository struct {
	feedURL string
	client  *http.Client
}

// NewPodcastRepository creates a new PodcastRepository instance
func NewPodcastRepository(feedURL string) *PodcastRepository {
	return &PodcastRepository{
		feedURL: feedURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchFeed retrieves the podcast feed from the remote URL
func (pr *PodcastRepository) FetchFeed() ([]*gofeed.Item, error) {
	resp, err := pr.client.Get(pr.feedURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseString(string(bodyBytes))
	if err != nil || len(feed.Items) == 0 {
		return nil, err
	}

	return feed.Items, nil
}
