package service

import (
	"fmt"
	"html/template"
	"math/rand"
	"regexp"
	"time"

	"github.com/kakakikikeke/random-podcast/models"
	"github.com/mmcdole/gofeed"
)

// PodcastRepository defines the interface for podcast data access
type PodcastRepository interface {
	FetchFeed() ([]*gofeed.Item, error)
}

// PodcastService handles business logic for podcasts
type PodcastService struct {
	repo          PodcastRepository
	rng           *rand.Rand
	aboutRegex    *regexp.Regexp
	showNoteRegex *regexp.Regexp
}

// NewPodcastService creates a new PodcastService instance
func NewPodcastService(repo PodcastRepository) *PodcastService {
	return &PodcastService{
		repo:          repo,
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
		aboutRegex:    regexp.MustCompile(`<p>(.*?)</p>`),
		showNoteRegex: regexp.MustCompile(`<ul id="menu">.*?</ul>`),
	}
}

// GetRandomPodcast fetches a random podcast episode from the feed
func (ps *PodcastService) GetRandomPodcast() (*models.Podcast, error) {
	items, err := ps.repo.FetchFeed()
	if err != nil {
		return nil, fmt.Errorf("service: %w", err)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no podcast items available")
	}

	// Select a random item
	item := items[ps.rng.Intn(len(items))]

	// Parse and return podcast
	podcast := ps.parseItem(item)
	return podcast, nil
}

// parseItem converts a gofeed.Item to a models.Podcast
func (ps *PodcastService) parseItem(item *gofeed.Item) *models.Podcast {
	podcast := &models.Podcast{
		Title:     item.Title,
		Published: item.Published,
	}

	if len(item.Enclosures) > 0 {
		podcast.URL = item.Enclosures[0].URL
	}

	// Parse description
	desc := item.Description
	podcast.About = ps.extractAbout(desc)
	podcast.ShowNote = ps.extractShowNote(desc)

	return podcast
}

// extractAbout extracts the "About" section from the description
func (ps *PodcastService) extractAbout(desc string) string {
	if m := ps.aboutRegex.FindStringSubmatch(desc); len(m) > 1 {
		return m[1]
	}
	return ""
}

// extractShowNote extracts the show notes menu from the description
func (ps *PodcastService) extractShowNote(desc string) template.HTML {
	if m := ps.showNoteRegex.FindString(desc); m != "" {
		return template.HTML(m)
	}
	return ""
}
