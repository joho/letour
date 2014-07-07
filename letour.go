package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/joho/prohttphandler"
)

func main() {
	portNumber := os.Getenv("PORT")
	if portNumber == "" {
		portNumber = "5000"
	}
	listenOn := fmt.Sprintf(":%v", portNumber)

	handler := prohttphandler.New("public")

	handler.ExactMatchFunc("/", func(w http.ResponseWriter, r *http.Request) {
		links := getLinks()
		t, err := template.ParseFiles("views/index.tmpl.html")
		if err == nil {
			t.Execute(w, links)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Fatal(http.ListenAndServe(listenOn, handler))
}

type Links []Link
type Link struct {
	Title    string
	VideoUrl string
}

type Feed struct {
	Entries []Entry `"json":"entries"`
}

type Entry struct {
	Title string  `"json":"title"`
	Media []Media `"json":"media$content"`
}

func (e *Entry) IsHighlight() bool {
	match, _ := regexp.MatchString("(?i)tour de france.+(stage \\d+|prologue).+highlights", e.Title)
	return match
}

func (e *Entry) HighBitrateMedia() *Media {
	for _, media := range e.Media {
		if media.IsHighBitrate() {
			return &media
		}
	}
	return nil
}

type Media struct {
	DownloadUrl string `"json":"plfile$downloadUrl"`
}

func (m *Media) IsHighBitrate() bool {
	return strings.Contains(m.DownloadUrl, "1500K")
}

func getLinks() *Links {
	feedUrl := "http://www.sbs.com.au/api/video_feed/f/Bgtm9B/sbs-section-sbstv/?range=1-100&byCategories=Sport/Cycling&form=json&defaultThumbnailAssetType=Thumbnail"

	res, err := http.Get(feedUrl)
	if err != nil {
		return nil
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}

	var feed Feed
	err = json.Unmarshal(body, &feed)
	if err != nil {
		panic(err)
	}

	highlights := []Entry{}
	for _, entry := range feed.Entries {
		if entry.IsHighlight() {
			highlights = append(highlights, entry)
		}
	}

	fmt.Printf("%v", highlights)
	fmt.Printf("%v", highlights[0].Media)

	links := Links{}

	for _, entry := range highlights {
		if entry.IsHighlight() {
			links = append(links, Link{
				Title:    entry.Title,
				VideoUrl: "",
			})
		}
	}
	return &links
}
