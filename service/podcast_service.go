package service

import (
	"html/template"
	"math/rand"
	"regexp"
	"time"

	"github.com/kakakikikeke/random-podcast/models"
	"github.com/kakakikikeke/random-podcast/repository"
	"github.com/mmcdole/gofeed"
)

// PodcastService handles business logic for podcasts
type PodcastService struct {
	repo *repository.PodcastRepository
}

// NewPodcastService creates a new PodcastService instance
func NewPodcastService(repo *repository.PodcastRepository) *PodcastService {
	return &PodcastService{
		repo: repo,
	}
}

// GetRandomPodcast fetches a random podcast episode from the feed
func (ps *PodcastService) GetRandomPodcast() (*models.Podcast, error) {
	items, err := ps.repo.FetchFeed()
	if err != nil {
		return nil, err
	}

	// Select a random item
	rand.Seed(time.Now().UnixNano())
	item := items[rand.Intn(len(items))]

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
	r := regexp.MustCompile(`<p>(.*?)</p>`)
	if m := r.FindStringSubmatch(desc); len(m) > 1 {
		return m[1]
	}
	return ""
}

// extractShowNote extracts the show notes menu from the description
func (ps *PodcastService) extractShowNote(desc string) template.HTML {
	r := regexp.MustCompile(`<ul id="menu">.*?</ul>`)
	if m := r.FindString(desc); m != "" {
		return template.HTML(m)
	}
	return ""
}
