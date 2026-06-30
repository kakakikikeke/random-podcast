// Package repository provides data access for podcast feeds
package repository

import (
	"fmt"
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

// NewPodcastRepositoryWithClient creates a new PodcastRepository with custom HTTP client
func NewPodcastRepositoryWithClient(feedURL string, client *http.Client) *PodcastRepository {
	return &PodcastRepository{
		feedURL: feedURL,
		client:  client,
	}
}

// FetchFeed retrieves the podcast feed from the remote URL
func (pr *PodcastRepository) FetchFeed() ([]*gofeed.Item, error) {
	resp, err := pr.client.Get(pr.feedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed from %s: %w", pr.feedURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("feed server returned status %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseString(string(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}

	if len(feed.Items) == 0 {
		return nil, fmt.Errorf("feed contains no items")
	}

	return feed.Items, nil
}
