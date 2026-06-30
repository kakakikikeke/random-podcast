package models

import "html/template"

// Podcast represents a single podcast episode
type Podcast struct {
	Title     string
	Published string
	ShowNote  template.HTML
	About     string
	URL       string
}

// RawFeedItem represents raw feed item data for parsing
type RawFeedItem struct {
	Title       string
	Published   string
	Description string
	AudioURL    string
}
