package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/kakakikikeke/random-podcast/handler"
	"github.com/kakakikikeke/random-podcast/repository"
	"github.com/kakakikikeke/random-podcast/service"
)

const (
	feedURL = "https://kakakikikeke.com/podcast/feed"
	port    = ":8080"
)

func main() {
	// Initialize template
	indexTmpl := template.Must(template.ParseFiles("index.html"))

	// Initialize layers
	repo := repository.NewPodcastRepository(feedURL)
	svc := service.NewPodcastService(repo)
	podcastHandler := handler.NewPodcastHandler(svc, indexTmpl)

	// Static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/favicon.ico", http.FileServer(http.Dir("static")))
	http.Handle("/podcast_icon.png", http.FileServer(http.Dir("static")))

	// Main handler
	http.Handle("/", podcastHandler)

	log.Println("Listening on" + port + "...")
	log.Fatal(http.ListenAndServe(port, nil))
}
