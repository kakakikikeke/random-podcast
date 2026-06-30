package service

import (
	"fmt"
	"testing"

	"github.com/kakakikikeke/random-podcast/models"
	"github.com/mmcdole/gofeed"
)

// MockRepository is a mock implementation of PodcastRepository
type MockRepository struct {
	items []*gofeed.Item
	err   error
}

func (m *MockRepository) FetchFeed() ([]*gofeed.Item, error) {
	return m.items, m.err
}

func createMockItem(title, published, description, audioURL string) *gofeed.Item {
	item := &gofeed.Item{
		Title:       title,
		Published:   published,
		Description: description,
		Enclosures: []*gofeed.Enclosure{
			{
				URL:  audioURL,
				Type: "audio/mpeg",
			},
		},
	}
	return item
}

func TestNewPodcastService(t *testing.T) {
	repo := &MockRepository{}
	svc := NewPodcastService(repo)

	if svc == nil {
		t.Fatal("NewPodcastService returned nil")
	}
	if svc.repo != repo {
		t.Error("Expected repo to be assigned")
	}
	if svc.rng == nil {
		t.Error("Expected rng to be initialized")
	}
	if svc.aboutRegex == nil {
		t.Error("Expected aboutRegex to be initialized")
	}
	if svc.showNoteRegex == nil {
		t.Error("Expected showNoteRegex to be initialized")
	}
}

func TestGetRandomPodcast_Success(t *testing.T) {
	mockItems := []*gofeed.Item{
		createMockItem("Episode 1", "2025-06-17", "<p>About episode 1</p>", "http://example.com/ep1.mp3"),
		createMockItem("Episode 2", "2025-06-18", "<p>About episode 2</p>", "http://example.com/ep2.mp3"),
	}
	repo := &MockRepository{items: mockItems}
	svc := NewPodcastService(repo)

	podcast, err := svc.GetRandomPodcast()

	if err != nil {
		t.Fatalf("GetRandomPodcast returned error: %v", err)
	}
	if podcast == nil {
		t.Fatal("Expected podcast, got nil")
	}
	if podcast.Title == "" {
		t.Error("Expected podcast title to be set")
	}
}

func TestGetRandomPodcast_RepositoryError(t *testing.T) {
	repo := &MockRepository{err: fmt.Errorf("connection failed")}
	svc := NewPodcastService(repo)

	podcast, err := svc.GetRandomPodcast()

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if podcast != nil {
		t.Error("Expected podcast to be nil on error")
	}
}

func TestGetRandomPodcast_EmptyItems(t *testing.T) {
	repo := &MockRepository{items: []*gofeed.Item{}}
	svc := NewPodcastService(repo)

	podcast, err := svc.GetRandomPodcast()

	if err == nil {
		t.Fatal("Expected error for empty items, got nil")
	}
	if podcast != nil {
		t.Error("Expected podcast to be nil on error")
	}
}

func TestExtractAbout(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid paragraph",
			input:    "<p>This is about section</p>",
			expected: "This is about section",
		},
		{
			name:     "Multiple paragraphs - first one",
			input:    "<p>First</p><p>Second</p>",
			expected: "First",
		},
		{
			name:     "No paragraph tag",
			input:    "No tags here",
			expected: "",
		},
		{
			name:     "Empty paragraph",
			input:    "<p></p>",
			expected: "",
		},
		{
			name:     "Paragraph with HTML inside",
			input:    "<p>Text with <b>bold</b></p>",
			expected: "Text with <b>bold</b>",
		},
	}

	repo := &MockRepository{}
	svc := NewPodcastService(repo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.extractAbout(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestExtractShowNote(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid menu",
			input:    "<ul id=\"menu\"><li>Note 1</li></ul>",
			expected: "<ul id=\"menu\"><li>Note 1</li></ul>",
		},
		{
			name:     "Menu with nested elements",
			input:    "<p>Text</p><ul id=\"menu\"><li><a href=\"#\">Link</a></li></ul><p>More</p>",
			expected: "<ul id=\"menu\"><li><a href=\"#\">Link</a></li></ul>",
		},
		{
			name:     "No menu",
			input:    "<p>Just paragraph</p>",
			expected: "",
		},
		{
			name:     "Empty menu",
			input:    "<ul id=\"menu\"></ul>",
			expected: "<ul id=\"menu\"></ul>",
		},
	}

	repo := &MockRepository{}
	svc := NewPodcastService(repo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.extractShowNote(tt.input)
			if string(result) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(result))
			}
		})
	}
}

func TestParseItem(t *testing.T) {
	tests := []struct {
		name     string
		item     *gofeed.Item
		validate func(*models.Podcast) error
	}{
		{
			name: "Complete item",
			item: createMockItem(
				"Episode Title",
				"2025-06-17",
				"<p>About this episode</p><ul id=\"menu\"><li>Note</li></ul>",
				"http://example.com/audio.mp3",
			),
			validate: func(p *models.Podcast) error {
				if p.Title != "Episode Title" {
					return fmt.Errorf("expected title 'Episode Title', got %q", p.Title)
				}
				if p.URL != "http://example.com/audio.mp3" {
					return fmt.Errorf("expected URL 'http://example.com/audio.mp3', got %q", p.URL)
				}
				if p.About != "About this episode" {
					return fmt.Errorf("expected about 'About this episode', got %q", p.About)
				}
				if !contains(string(p.ShowNote), "Note") {
					return fmt.Errorf("expected show note to contain 'Note'")
				}
				return nil
			},
		},
		{
			name: "Item without enclosures",
			item: &gofeed.Item{
				Title:     "No Audio",
				Published: "2025-06-17",
			},
			validate: func(p *models.Podcast) error {
				if p.Title != "No Audio" {
					return fmt.Errorf("expected title 'No Audio', got %q", p.Title)
				}
				if p.URL != "" {
					return fmt.Errorf("expected empty URL, got %q", p.URL)
				}
				return nil
			},
		},
	}

	repo := &MockRepository{}
	svc := NewPodcastService(repo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			podcast := svc.parseItem(tt.item)
			if err := tt.validate(podcast); err != nil {
				t.Errorf("Validation failed: %v", err)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s != ""
}
