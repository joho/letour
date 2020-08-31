package sbs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type MediaItemFeed []MediaItem

type MediaItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ProgramName string `json:"programName"`
}

func (m *MediaItem) Url() string {
	return fmt.Sprintf("https://www.sbs.com.au/ondemand/video/single/%v/?source=drupal&vertical=cyclingcentral\n", m.ID)
}

// GetHighlights returns all the videos that _should_ be highlights
func GetHighlights() MediaItemFeed {
	highlightsFeed := MediaItemFeed{}

	titleRegexp := regexp.MustCompile(`(?i)\: tour de france 2020`)
	for _, mediaItem := range AllVideos() {
		if titleRegexp.MatchString(mediaItem.Title) {
			highlightsFeed = append(highlightsFeed, mediaItem)
		}
	}

	return highlightsFeed
}

// AllVideos returns every media item we find on the SBS feed
func AllVideos() MediaItemFeed {
	// curl http://www.sbs.com.au/cyclingcentral/?cid=infocus
	// SBS.mpxWidget.setVideos

	// generate urls like:
	//		http://www.sbs.com.au/ondemand/video/single/713877571952/?source=drupal&vertical=cyclingcentral
	res, err := http.Get("https://www.sbs.com.au/cyclingcentral/")
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	r := regexp.MustCompile(`SBS.mpxWidget.setVideos\(mpxBeanId, (\[.+\]), \d.+;`)
	subMatches := r.FindAllSubmatch(body, -1)
	
	if len(subMatches) == 0 {
		panic("didn't find the mega json array")
	}

	var fullFeed = MediaItemFeed{}

	for _, v := range subMatches {
		var feed MediaItemFeed
		err = json.Unmarshal(v[1], &feed)
		if err != nil {
			panic(err)
		}
		
		fullFeed = append(fullFeed, feed...)
	}
	
	return fullFeed
}
