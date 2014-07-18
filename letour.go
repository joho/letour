package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/bugsnag/bugsnag-go"
	"github.com/joho/prohttphandler"
)

func main() {
	portNumber := os.Getenv("PORT")
	if portNumber == "" {
		portNumber = "5000"
	}
	listenOn := fmt.Sprintf(":%v", portNumber)

	appHandler := prohttphandler.New("public")

	appHandler.ExactMatchFunc("/", func(w http.ResponseWriter, r *http.Request) {
		links := getLinks()
		t, err := template.ParseFiles("views/index.tmpl.html")
		if err == nil {
			t.Execute(w, links)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	bugsnagApiKey := os.Getenv("BUGSNAG_API_KEY")
	handler := bugsnagHttpHandler(bugsnagApiKey, appHandler)

	log.Fatal(http.ListenAndServe(listenOn, handler))
}

func bugsnagHttpHandler(bugsnagApiKey string, handler http.Handler) http.Handler {

	if bugsnagApiKey == "" {
		// env doesn't support error handling
		return handler
	} else {
		bugsnag.Configure(bugsnag.Configuration{
			APIKey: bugsnagApiKey,
		})
		return bugsnag.Handler(handler)
	}
}

type Feed struct {
	Entries []Entry `json:"entries"`
}

func (f *Feed) Highlights() []Entry {
	highlights := []Entry{}
	for _, entry := range f.Entries {
		if entry.IsHighlight() {
			highlights = append(highlights, entry)
		}
	}
	return highlights
}

type Entry struct {
	Title string  `json:"title"`
	Media []Media `json:"media$content"`
}

func (e *Entry) IsHighlight() bool {
	titleRegexps := []string{
		"(?i)tour de france.+stage \\d+.+highlights",
		"(?i)tour de france highlights.+stage \\d+",
	}

	titleMatch := false
	for _, titleRegexp := range titleRegexps {
		if match, _ := regexp.MatchString(titleRegexp, e.Title); match {
			titleMatch = true
			break
		}
	}

	return titleMatch
}

func (e *Entry) HighBitrateMedia() *Media {
	for _, media := range e.Media {
		if media.IsHighBitrate() {
			return &media
		}
	}
	return nil
}

func (e *Entry) VideoUrl() string {
	return e.HighBitrateMedia().DownloadUrl
}

type Media struct {
	DownloadUrl string `json:"plfile$downloadUrl"`
}

func (m *Media) IsHighBitrate() bool {
	return strings.Contains(m.DownloadUrl, "1500K")
}

func getLinks() []Entry {
	feedUrl := "http://www.sbs.com.au/api/video_feed/f/Bgtm9B/sbs-section-sbstv/?range=1-100&byCategories=Sport/Cycling&form=json&defaultThumbnailAssetType=Thumbnail"

	res, err := http.Get(feedUrl)
	if err != nil {
		return nil
	}

	defer res.Body.Close()

	var feed Feed
	err = json.NewDecoder(res.Body).Decode(&feed)
	if err != nil {
		return nil
	}

	return feed.Highlights()
}
