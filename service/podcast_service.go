// Package service provides business logic for podcast processing
package service

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"math/big"
	"regexp"
	"strings"

	"github.com/kakakikikeke/random-podcast/models"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
)

// PodcastRepository defines the interface for podcast data access
type PodcastRepository interface {
	FetchFeed() ([]*gofeed.Item, error)
}

// PodcastService handles business logic for podcasts
type PodcastService struct {
	repo          PodcastRepository
	randomInt     func(int) int
	aboutRegex    *regexp.Regexp
	showNoteRegex *regexp.Regexp
}

// NewPodcastService creates a new PodcastService instance
func NewPodcastService(repo PodcastRepository) *PodcastService {
	return &PodcastService{
		repo:          repo,
		randomInt:     secureRandomInt,
		aboutRegex:    regexp.MustCompile(`<p>(.*?)</p>`),
		showNoteRegex: regexp.MustCompile(`<ul id="menu">.*?</ul>`),
	}
}

func secureRandomInt(max int) int {
	if max <= 0 {
		return 0
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}
	return int(n.Int64())
}

// GetRandomPodcast fetches a random podcast episode from the feed
func (ps *PodcastService) GetRandomPodcast() (*models.Podcast, error) {
	items, err := ps.repo.FetchFeed()
	if err != nil {
		return nil, fmt.Errorf("service: %w", err)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no podcast items available")
	}

	// Select a random item using a cryptographically secure PRNG.
	item := items[ps.randomInt(len(items))]

	// Parse and return podcast
	podcast := ps.parseItem(item)
	return podcast, nil
}

// parseItem converts a gofeed.Item to a models.Podcast
func (ps *PodcastService) parseItem(item *gofeed.Item) *models.Podcast {
	podcast := &models.Podcast{
		Title:     item.Title,
		Published: item.Published,
	}

	if len(item.Enclosures) > 0 {
		podcast.URL = item.Enclosures[0].URL
	}

	// Parse description
	desc := item.Description
	podcast.About = ps.extractAbout(desc)
	podcast.ShowNote = ps.extractShowNote(desc)

	return podcast
}

// extractAbout extracts the "About" section from the description
func (ps *PodcastService) extractAbout(desc string) string {
	if m := ps.aboutRegex.FindStringSubmatch(desc); len(m) > 1 {
		return m[1]
	}
	return ""
}

// extractShowNote extracts and sanitizes the show notes menu from the description.
// The podcast feed is treated as trusted input, but we still strip unsafe tags/attributes
// before rendering as HTML in the template.
func (ps *PodcastService) extractShowNote(desc string) template.HTML {
	if m := ps.showNoteRegex.FindString(desc); m != "" {
		safe := sanitizeHTML(m)
		if safe != "" {
			// nosemgrep: go.lang.security.audit.xss.template-html-does-not-escape.unsafe-template-type
			return template.HTML(safe)
		}
	}
	return ""
}

func sanitizeHTML(raw string) string {
	nodes, err := html.ParseFragment(strings.NewReader(raw), nil)
	if err != nil {
		return ""
	}

	var b strings.Builder
	for _, n := range nodes {
		if sanitized := sanitizeNode(n); sanitized != nil {
			renderNode(sanitized, &b)
		}
	}
	return strings.TrimSpace(b.String())
}

func renderNode(n *html.Node, b *strings.Builder) {
	if n == nil {
		return
	}
	if n.Type == html.ElementNode && (n.Data == "html" || n.Data == "head" || n.Data == "body") {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			renderNode(c, b)
		}
		return
	}
	_ = html.Render(b, n)
}

func sanitizeNode(n *html.Node) *html.Node {
	if n == nil {
		return nil
	}

	if n.Type == html.ElementNode {
		switch n.Data {
		case "html", "head", "body":
			// Fragments are wrapped by the parser; allow the wrapper nodes.
		case "script", "iframe", "object", "embed", "svg", "math", "style", "link", "meta", "base":
			return nil
		case "a":
			filtered := n.Attr[:0]
			for _, attr := range n.Attr {
				switch {
				case strings.HasPrefix(strings.ToLower(attr.Key), "on"):
					continue
				case strings.EqualFold(attr.Key, "href") && strings.HasPrefix(strings.ToLower(attr.Val), "javascript:"):
					continue
				case !isAllowedAttribute(attr.Key):
					continue
				default:
					filtered = append(filtered, attr)
				}
			}
			n.Attr = filtered
		case "ul", "ol", "li", "p", "div", "span", "br", "strong", "b", "em", "i", "nav":
			// allow these tags without removing attributes.
		default:
			return nil
		}
	}

	for child := n.FirstChild; child != nil; {
		next := child.NextSibling
		if sanitized := sanitizeNode(child); sanitized == nil {
			n.RemoveChild(child)
			child = next
		} else {
			child = next
		}
	}
	return n
}

func isAllowedAttribute(key string) bool {
	switch strings.ToLower(key) {
	case "href", "title", "target", "rel", "id", "class":
		return true
	default:
		return false
	}
}
