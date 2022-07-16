package sbs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type NextData struct {
	Props Props `json:"props"`
}

type Props struct {
	PageProps PageProps `json:"pageProps"`
}

type PageProps struct {
	PageContent PageContent `json:"pageContent"`
}

type PageContent struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	TypeName string `json:"__typename"`
	Title    string `json:"title"`
	Items    []Item `json:"items"`
}

type Item struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Route Route  `json:"route"`
}

type Route struct {
	RoutePaths []string `json:"routePaths"`
}

/*
{
            "__typename": "CardShelf",
            "id": "00000181-bac6-d3c6-a997-bef655430000",
            "blockType": "CardShelf",
            "blockTheme": "OD",
            "title": "Daily highlights",
            "maxItems": 4,
            "callToAction": {
              "route": null,
              "url": "https://www.sbs.com.au/ondemand/program/tour-de-france-2022",
              "__typename": "CallToAction",
              "text": "Visit SBS On Demand for more highlights"
            },
            "items": [
              {
                "type": "video",
                "title": "Daily Highlights Stage 13: Tour de France 2022",
                "linkUrl": null,
                "linkTarget": null,
                "locale": {
                  "languageCode": "en-AU",
                  "__typename": "Locale"
                },
*/

func (m *Item) Url() string {
	if m != nil {
		return fmt.Sprintf("https://www.sbs.com.au%v\n", m.Route.RoutePaths[0])
	}

	return ""
}

// AllVideos returns every media item we find on the SBS feed
func GetHighlights() []Item {
	// curl http://www.sbs.com.au/cyclingcentral/?cid=infocus
	// SBS.mpxWidget.setVideos

	// <script id="__NEXT_DATA__" type="application/json">

	// generate urls like:
	//		http://www.sbs.com.au/ondemand/video/single/713877571952/?source=drupal&vertical=cyclingcentral
	res, err := http.Get("https://www.sbs.com.au/sport/tour-de-france")
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	r := regexp.MustCompile(`<script id="__NEXT_DATA__" type="application/json">(.*)</script>`)
	subMatches := r.FindAllSubmatch(body, -1)

	if len(subMatches) != 1 {
		panic("didn't find the mega json array")
	}

	var nextData = NextData{}

	fmt.Printf("%v\n", string(subMatches[0][1]))

	err = json.Unmarshal(subMatches[0][1], &nextData)
	if err != nil {
		panic(err)
	}

	items := []Item{}

	blocks := nextData.Props.PageProps.PageContent.Blocks
	for _, block := range blocks {
		if block.TypeName == "CardShelf" && block.Title == "Daily highlights" {
			items = append(items, items...)
		}
	}

	return items
}
