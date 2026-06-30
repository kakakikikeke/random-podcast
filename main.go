package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/kakakikikeke/random-podcast/config"
	"github.com/kakakikikeke/random-podcast/handler"
	"github.com/kakakikikeke/random-podcast/repository"
	"github.com/kakakikikeke/random-podcast/service"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Initialize template
	indexTmpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatalf("failed to parse index.html: %v", err)
	}

	// Initialize layers
	repo := repository.NewPodcastRepository(cfg.FeedURL)
	svc := service.NewPodcastService(repo)
	podcastHandler := handler.NewPodcastHandler(svc, indexTmpl)

	// Static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/favicon.ico", http.FileServer(http.Dir("static")))
	http.Handle("/podcast_icon.png", http.FileServer(http.Dir("static")))

	// Main handler
	http.Handle("/", podcastHandler)

	log.Printf("Listening on %s...", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, nil))
}
