package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/kakakikikeke/random-podcast/service"
)

// PodcastHandler handles HTTP requests for podcasts
type PodcastHandler struct {
	service   *service.PodcastService
	indexTmpl *template.Template
}

// NewPodcastHandler creates a new PodcastHandler instance
func NewPodcastHandler(svc *service.PodcastService, indexTemplate *template.Template) *PodcastHandler {
	return &PodcastHandler{
		service:   svc,
		indexTmpl: indexTemplate,
	}
}

// ServeHTTP handles the main podcast request
func (ph *PodcastHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	podcast, err := ph.service.GetRandomPodcast()
	if err != nil {
		log.Printf("Error getting random podcast: %v", err)
		http.Error(w, "Failed to fetch podcast", http.StatusInternalServerError)
		return
	}

	if err := ph.indexTmpl.Execute(w, podcast); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}
