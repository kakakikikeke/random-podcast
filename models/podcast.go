// Package models defines data models for the podcast application
package models

import "html/template"

// Podcast represents a single podcast episode
type Podcast struct {
	// Title is the episode title
	Title string

	// Published is the publication date
	Published string

	// ShowNote contains structured show notes as HTML
	ShowNote template.HTML

	// About is a brief description of the episode
	About string

	// URL is the audio file URL
	URL string
}
