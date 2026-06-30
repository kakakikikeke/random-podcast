package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kakakikikeke/random-podcast/models"
)

// MockPodcastService is a mock implementation of PodcastServicer
type MockPodcastService struct {
	podcast *models.Podcast
	err     error
}

func (m *MockPodcastService) GetRandomPodcast() (*models.Podcast, error) {
	return m.podcast, m.err
}

func TestNewPodcastHandler(t *testing.T) {
	tmpl := template.New("test")
	svc := &MockPodcastService{}
	handler := NewPodcastHandler(svc, tmpl)

	if handler == nil {
		t.Fatal("NewPodcastHandler returned nil")
	}
	if handler.service != svc {
		t.Error("Expected service to be assigned")
	}
	if handler.indexTmpl != tmpl {
		t.Error("Expected template to be assigned")
	}
}

func TestServeHTTP_Success(t *testing.T) {
	mockPodcast := &models.Podcast{
		Title:     "Test Episode",
		Published: "2025-06-17",
		About:     "Test about section",
		URL:       "http://example.com/audio.mp3",
	}

	svc := &MockPodcastService{podcast: mockPodcast}

	// Create a simple template
	tmpl, err := template.New("test").Parse(`
		<html>
			<body>
				<h1>{{.Title}}</h1>
				<p>{{.About}}</p>
				<audio src="{{.URL}}"></audio>
			</body>
		</html>
	`)
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	handler := NewPodcastHandler(svc, tmpl)

	// Create HTTP request and response recorder
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(rec, req)

	// Verify status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Verify content type
	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("Expected Content-Type to contain 'text/html', got %q", contentType)
	}

	// Verify response body contains expected content
	body := rec.Body.String()
	if !strings.Contains(body, "Test Episode") {
		t.Errorf("Expected response to contain 'Test Episode', got %q", body)
	}
	if !strings.Contains(body, "Test about section") {
		t.Errorf("Expected response to contain 'Test about section'")
	}
}

func TestServeHTTP_ServiceError(t *testing.T) {
	svc := &MockPodcastService{
		err: fmt.Errorf("service failed"),
	}

	tmpl, _ := template.New("test").Parse("<h1>Test</h1>")
	handler := NewPodcastHandler(svc, tmpl)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Verify error response
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Failed to fetch podcast") {
		t.Errorf("Expected error message in response")
	}
}

func TestServeHTTP_TemplateError(t *testing.T) {
	mockPodcast := &models.Podcast{
		Title: "Test Episode",
	}

	svc := &MockPodcastService{podcast: mockPodcast}

	// Create a template that will fail during execution
	tmpl, _ := template.New("test").Parse("{{.NonExistent}}")
	handler := NewPodcastHandler(svc, tmpl)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Verify error response for template error
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500 for template error, got %d", rec.Code)
	}
}

func TestServeHTTP_HeadersSet(t *testing.T) {
	mockPodcast := &models.Podcast{
		Title: "Test Episode",
	}

	svc := &MockPodcastService{podcast: mockPodcast}
	tmpl, _ := template.New("test").Parse("<h1>{{.Title}}</h1>")
	handler := NewPodcastHandler(svc, tmpl)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Verify headers
	contentType := rec.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/html; charset=utf-8', got %q", contentType)
	}
}
