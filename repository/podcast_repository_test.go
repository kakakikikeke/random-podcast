package repository

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const mockFeed = `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">
  <channel>
    <title>Test Podcast</title>
    <item>
      <title>Episode 1</title>
      <pubDate>Mon, 17 Jun 2025 10:00:00 +0000</pubDate>
      <description><![CDATA[<p>Test about</p>]]></description>
      <enclosure url="http://example.com/audio1.mp3" type="audio/mpeg" />
    </item>
    <item>
      <title>Episode 2</title>
      <pubDate>Mon, 17 Jun 2025 11:00:00 +0000</pubDate>
      <description><![CDATA[<p>Another about</p>]]></description>
      <enclosure url="http://example.com/audio2.mp3" type="audio/mpeg" />
    </item>
  </channel>
</rss>`

func TestNewPodcastRepository(t *testing.T) {
	repo := NewPodcastRepository("https://example.com/feed")
	if repo == nil {
		t.Fatal("NewPodcastRepository returned nil")
	}
	if repo.feedURL != "https://example.com/feed" {
		t.Errorf("Expected feedURL to be 'https://example.com/feed', got %s", repo.feedURL)
	}
	if repo.client == nil {
		t.Fatal("Expected client to be initialized")
	}
}

func TestFetchFeed_Success(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockFeed))
	}))
	defer server.Close()

	repo := NewPodcastRepository(server.URL)
	items, err := repo.FetchFeed()

	if err != nil {
		t.Fatalf("FetchFeed returned error: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	if items[0].Title != "Episode 1" {
		t.Errorf("Expected first item title 'Episode 1', got %s", items[0].Title)
	}

	if items[1].Title != "Episode 2" {
		t.Errorf("Expected second item title 'Episode 2', got %s", items[1].Title)
	}
}

func TestFetchFeed_NetworkError(t *testing.T) {
	repo := NewPodcastRepository("http://invalid-domain-that-does-not-exist.test")
	_, err := repo.FetchFeed()

	if err == nil {
		t.Fatal("Expected error for invalid domain, got nil")
	}
}

func TestFetchFeed_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	repo := NewPodcastRepository(server.URL)
	_, err := repo.FetchFeed()

	if err == nil {
		t.Fatal("Expected error for HTTP 404, got nil")
	}
}

func TestFetchFeed_InvalidFeed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Invalid XML"))
	}))
	defer server.Close()

	repo := NewPodcastRepository(server.URL)
	_, err := repo.FetchFeed()

	if err == nil {
		t.Fatal("Expected error for invalid feed, got nil")
	}
}

func TestFetchFeed_EmptyFeed(t *testing.T) {
	emptyFeed := `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">
  <channel>
    <title>Empty Podcast</title>
  </channel>
</rss>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(emptyFeed))
	}))
	defer server.Close()

	repo := NewPodcastRepository(server.URL)
	_, err := repo.FetchFeed()

	if err == nil {
		t.Fatal("Expected error for empty feed, got nil")
	}
}

func TestNewPodcastRepositoryWithClient(t *testing.T) {
	customClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	repo := NewPodcastRepositoryWithClient("https://example.com/feed", customClient)

	if repo == nil {
		t.Fatal("NewPodcastRepositoryWithClient returned nil")
	}
	if repo.client != customClient {
		t.Error("Expected custom client to be assigned")
	}
}
