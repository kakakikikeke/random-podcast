package config

import (
	"os"
	"testing"
	"time"
)

func TestNewConfig_Defaults(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("PODCAST_FEED_URL")
	os.Unsetenv("PORT")
	os.Unsetenv("HTTP_CLIENT_TIMEOUT")

	cfg := NewConfig()

	if cfg.FeedURL != "https://kakakikikeke.com/podcast/feed" {
		t.Errorf("Expected default FeedURL, got %s", cfg.FeedURL)
	}
	if cfg.Port != ":8080" {
		t.Errorf("Expected default Port, got %s", cfg.Port)
	}
	if cfg.HTTPClientTimeout != 10*time.Second {
		t.Errorf("Expected default HTTPClientTimeout 10s, got %v", cfg.HTTPClientTimeout)
	}
}

func TestNewConfig_FromEnvironment(t *testing.T) {
	// Set environment variables
	os.Setenv("PODCAST_FEED_URL", "https://custom.feed/rss")
	os.Setenv("PORT", ":9000")
	os.Setenv("HTTP_CLIENT_TIMEOUT", "5")
	defer func() {
		os.Unsetenv("PODCAST_FEED_URL")
		os.Unsetenv("PORT")
		os.Unsetenv("HTTP_CLIENT_TIMEOUT")
	}()

	cfg := NewConfig()

	if cfg.FeedURL != "https://custom.feed/rss" {
		t.Errorf("Expected custom FeedURL, got %s", cfg.FeedURL)
	}
	if cfg.Port != ":9000" {
		t.Errorf("Expected custom Port, got %s", cfg.Port)
	}
	if cfg.HTTPClientTimeout != 5*time.Second {
		t.Errorf("Expected custom HTTPClientTimeout 5s, got %v", cfg.HTTPClientTimeout)
	}
}

func TestNewConfig_DurationFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected time.Duration
	}{
		{"Seconds string", "30", 30 * time.Second},
		{"Duration string", "1m30s", 1*time.Minute + 30*time.Second},
		{"Short duration", "500ms", 500 * time.Millisecond},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("HTTP_CLIENT_TIMEOUT", tt.envValue)
			defer os.Unsetenv("HTTP_CLIENT_TIMEOUT")

			cfg := NewConfig()
			if cfg.HTTPClientTimeout != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, cfg.HTTPClientTimeout)
			}
		})
	}
}

func TestNewConfig_InvalidDurationFallback(t *testing.T) {
	os.Setenv("HTTP_CLIENT_TIMEOUT", "invalid")
	defer os.Unsetenv("HTTP_CLIENT_TIMEOUT")

	cfg := NewConfig()

	// Should fallback to default
	if cfg.HTTPClientTimeout != 10*time.Second {
		t.Errorf("Expected fallback to 10s, got %v", cfg.HTTPClientTimeout)
	}
}
