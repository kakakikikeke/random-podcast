package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds the application configuration
type Config struct {
	FeedURL           string
	Port              string
	HTTPClientTimeout time.Duration
}

// NewConfig creates a config from environment variables with defaults
func NewConfig() *Config {
	return &Config{
		FeedURL:           getEnv("PODCAST_FEED_URL", "https://kakakikikeke.com/podcast/feed"),
		Port:              normalizePort(getEnv("PORT", ":8080")),
		HTTPClientTimeout: getDurationEnv("HTTP_CLIENT_TIMEOUT", 10*time.Second),
	}
}

func normalizePort(port string) string {
	if port == "" {
		return ":8080"
	}

	if strings.HasPrefix(port, ":") {
		return port
	}

	return ":" + port
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	// Try to parse as seconds
	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second
	}

	// Try to parse as duration string (e.g., "10s", "1m")
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}

	return defaultValue
}
