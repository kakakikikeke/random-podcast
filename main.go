package main

import (
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/mmcdole/gofeed"
)

var (
	index   = template.Must(template.ParseFiles("index.html"))
	feedURL = "https://kakakikikeke.com/podcast/feed"
)

type params struct {
	Title     string
	Published string
	ShowNote  template.HTML
	About     string
	URL       string
}

func main() {
	// 静的ファイル配信 (例: /static/style.css → static/style.css)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// favicon など root に配置したファイルも個別で対応
	http.Handle("/favicon.ico", http.FileServer(http.Dir("static")))
	http.Handle("/podcast_icon.png", http.FileServer(http.Dir("static")))
	// サーバ起動
	http.HandleFunc("/", handle)
	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	// コンテキスト作成、httpclient と logging などに使用
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	// feed 取得
	resp, err := client.Get(feedURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch feed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read feed", http.StatusInternalServerError)
		return
	}
	// feed パース
	fp := gofeed.NewParser()
	feed, err := fp.ParseString(string(bodyBytes))
	if err != nil || len(feed.Items) == 0 {
		http.Error(w, "Failed to parse feed", http.StatusInternalServerError)
		return
	}
	// feed からランダムに取得
	rand.Seed(time.Now().UnixNano())
	item := feed.Items[rand.Intn(len(feed.Items))]
	p := params{
		Title:     item.Title,
		Published: item.Published,
		URL:       item.Enclosures[0].URL,
	}
	// Descriptionの解析
	desc := item.Description
	r1 := regexp.MustCompile(`<p>(.*?)</p>`)
	if m := r1.FindStringSubmatch(desc); len(m) > 1 {
		p.About = m[1]
	}
	r2 := regexp.MustCompile(`<ul id="menu">.*?</ul>`)
	if m := r2.FindString(desc); m != "" {
		p.ShowNote = template.HTML(m)
	}
	if err := index.Execute(w, p); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}
