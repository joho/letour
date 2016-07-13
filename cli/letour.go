package main

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

func main() {
	// curl http://www.sbs.com.au/cyclingcentral/?cid=infocus
	// SBS.mpxWidget.setVideos

	// generate urls like:
	//		http://www.sbs.com.au/ondemand/video/single/713877571952/?source=drupal&vertical=cyclingcentral
	res, err := http.Get("http://www.sbs.com.au/cyclingcentral/?cid=infocus")
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	r := regexp.MustCompile(`SBS.mpxWidget.setVideos\(mpxBeanId, (\[.+\]), \d.+;`)
	subMatches := r.FindSubmatch(body)
	if len(subMatches) != 2 {
		panic("didn't find the mega json array")
	}

	// fmt.Println(string(subMatches[1]))

	var feed MediaItemFeed
	err = json.Unmarshal(subMatches[1], &feed)
	if err != nil {
		panic(err)
	}

	for _, mediaItem := range feed {
		if mediaItem.ProgramName == "Tour De France: Highlights" {
			fmt.Printf("%v\n", mediaItem.Title)
			fmt.Printf("http://www.sbs.com.au/ondemand/video/single/%v/?source=drupal&vertical=cyclingcentral\n", mediaItem.ID)
		}
	}
}
