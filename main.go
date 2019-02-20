package main

import (
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/mmcdole/gofeed"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

var (
	index = template.Must(template.ParseFiles("index.html"))
)

type params struct {
	Title     string
	Published string
	ShowNote  template.HTML
	About     string
	URL       string
}

func main() {
	http.HandleFunc("/", handle)
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	// コンテキスト作成、httpclient と logging などに使用
	ctx := appengine.NewContext(r)
	// feed 取得
	client := urlfetch.Client(ctx)
	resp, _ := client.Get("https://kakakikikeke.com/podcast/feed")
	defer resp.Body.Close()
	body := ""
	if resp.StatusCode == http.StatusOK {
		bytes, _ := ioutil.ReadAll(resp.Body)
		body = string(bytes)
	}
	// feed パース
	fp := gofeed.NewParser()
	feed, _ := fp.ParseString(body)
	items := feed.Items
	// feed からランダムに取得
	rand.Seed(time.Now().Unix())
	item := items[rand.Intn(len(items))]
	// テンプレート用の変数を設定
	params := params{}
	// タイトル
	params.Title = item.Title
	// 公開日
	params.Published = item.Published
	// 正規表現を使って概要を取得、Submatch でカッコ内を取得
	p1 := `<p>(.*)</p>`
	r1 := regexp.MustCompile(p1)
	ret1 := r1.FindStringSubmatch(item.Description)
	params.About = ret1[1]
	// 正規表現を使って show note の HTML を取得
	p2 := `<ul id="menu">.*</ul>`
	r2 := regexp.MustCompile(p2)
	ret2 := r2.FindString(item.Description)
	params.ShowNote = template.HTML(ret2)
	// 音声 URL を取得
	params.URL = item.Enclosures[0].URL
	// HTML 出力
	index.Execute(w, params)
}
