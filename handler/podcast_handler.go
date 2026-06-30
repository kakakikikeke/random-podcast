// Package handler provides HTTP request handling for the podcast application
package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/kakakikikeke/random-podcast/models"
)

// PodcastServicer defines the interface for podcast service
type PodcastServicer interface {
	GetRandomPodcast() (*models.Podcast, error)
}

// PodcastHandler handles HTTP requests for podcasts
type PodcastHandler struct {
	service   PodcastServicer
	indexTmpl *template.Template
}

// NewPodcastHandler creates a new PodcastHandler instance
func NewPodcastHandler(svc PodcastServicer, indexTemplate *template.Template) *PodcastHandler {
	return &PodcastHandler{
		service:   svc,
		indexTmpl: indexTemplate,
	}
}

// ServeHTTP handles the main podcast request
func (ph *PodcastHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	podcast, err := ph.service.GetRandomPodcast()
	if err != nil {
		log.Printf("error getting random podcast: %v", err)
		http.Error(w, "Failed to fetch podcast", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := ph.indexTmpl.Execute(w, podcast); err != nil {
		log.Printf("template rendering error: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}
