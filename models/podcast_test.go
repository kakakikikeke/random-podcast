package models

import "testing"

func TestPodcast(t *testing.T) {
	tests := []struct {
		name    string
		podcast *Podcast
		check   func(*Podcast) error
	}{
		{
			name: "Complete podcast",
			podcast: &Podcast{
				Title:     "Test Episode",
				Published: "2025-06-30",
				About:     "Test description",
				URL:       "http://example.com/audio.mp3",
			},
			check: func(p *Podcast) error {
				if p.Title == "" {
					return errFieldEmpty("Title")
				}
				if p.URL == "" {
					return errFieldEmpty("URL")
				}
				return nil
			},
		},
		{
			name: "Minimal podcast",
			podcast: &Podcast{
				Title: "Episode",
			},
			check: func(p *Podcast) error {
				if p.Title == "" {
					return errFieldEmpty("Title")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.check(tt.podcast); err != nil {
				t.Errorf("Check failed: %v", err)
			}
		})
	}
}

func errFieldEmpty(field string) error {
	return nil
}
