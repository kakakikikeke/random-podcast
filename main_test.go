package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// テスト用のモックRSSフィード
const mockFeed = `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">
  <channel>
    <title>Mock Podcast</title>
    <item>
      <title>Test Episode</title>
      <pubDate>Mon, 17 Jun 2025 10:00:00 +0000</pubDate>
      <description>
        <![CDATA[
        <p>This is a test episode.</p>
        <ul id="menu"><li><a href="#">Note 1</a></li><li><a href="#">Note 2</a></li></ul>
        ]]>
      </description>
      <enclosure url="http://example.com/audio.mp3" type="audio/mpeg" />
    </item>
  </channel>
</rss>`

func TestHandle(t *testing.T) {
	// モックフィードサーバーを起動
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockFeed))
	}))
	defer server.Close()

	// グローバル変数 feedURL をモックサーバーに差し替え
	originalFeedURL := feedURL
	feedURL = server.URL
	defer func() { feedURL = originalFeedURL }()

	// リクエストを作成
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// 実行
	handle(rr, req)

	// 検証
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}

	body := rr.Body.String()

	// 各項目の存在を確認
	if !strings.Contains(body, "Test Episode") {
		t.Errorf("Response missing title: %s", body)
	}
	if !strings.Contains(body, "This is a test episode.") {
		t.Errorf("Response missing description")
	}
	if !strings.Contains(body, "<ul id=\"menu\">") {
		t.Errorf("Response missing show notes")
	}
	if !strings.Contains(body, "http://example.com/audio.mp3") {
		t.Errorf("Response missing audio URL")
	}
}
