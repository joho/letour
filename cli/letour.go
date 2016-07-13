package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

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

	//  SBS.mpxWidget.setVideos(mpxBeanId, [{"id":"723969091538","idUrl":"http:\/\/data.media.theplatform.com\/media\/data\/Media\/723969091538","defaultThumbnail":"http:\/\/videocdn.sbs.com.au\/u\/video\/SBS\/managed\/images\/2016\/07\/13\/723969091538_07131207_image124125_large.jpg","pubDate":1468371600000,"availableDate":1468371600000,"expirationDate":1499734920000,"description":"Share the spirit of the Tour with OFX","title":"Share the spirit of the Tour with OFX","programName":"","content":[{"plfile$bitrate":128000,"plfile$contentType":"video","plfile$duration":63,"plfile$format":"MPEG4","plfile$height":0,"plfile$width":0,"plfile$assetTypes":["Public"],"plfile$downloadUrl":"http:\/\/videocdn.sbs.com.au\/u\/video\/SBS_Production\/managed\/2016\/07\/13\/723969091538_128K.mp4"},{"plfile$bitrate":1000000,"plfile$contentType":"video","plfile$duration":63,"plfile$format":"MPEG4","plfile$height":0,"plfile$width":0,"plfile$assetTypes":["Public"],"plfile$downloadUrl":"http:\/\/videocdn.sbs.com.au\/u\/video\/SBS_Production\/managed\/2016\/07\/13\/723969091538_1000K.mp4"},{"plfile$bitrate":512000,"plfile$contentType":"video","plfile$duration":63,"plfile$format":"MPEG4","plfile$height":0,"plfile$width":0,"plfile$assetTypes":["Public"],"plfile$downloadUrl":"http:\/\/videocdn.sbs.com.au\/u\/video\/SBS_Production\/managed\/2016\/07\/13\/723969091538_300K.mp4"},{"plfile$bitrate":1500000,"plfile$contentType":"video","plfile$duration":63,"plfile$format":"MPEG4","plfile$height":0,"plfile$width":0,"plfile$assetTypes":["Public"],"plfile$downloadUrl":"http:\/\/videocdn.sbs.com.au\/u\/video\/SBS_Production\/managed\/2016\/07\/13\/723969091538_1500K.mp4"}],"categories":[{"media$name":"Sport","media$scheme":"Genre","media$label":""},{"media$name":"Sport\/Cycling","media$scheme":"","media$label":""},{"media$name":"Section","media$scheme":"Section","media$label":""},{"media$name":"Section\/Clips","media$scheme":
	//
	// SNIP
	//
	// http:\/\/videocdn.sbs.com.au\/u\/video\/SBS\/managed\/images\/2016\/07\/11\/722350147591_07110707_image074057_facebook.jpg","Thumbnail Carousel Small":"http:\/\/videocdn.sbs.com.au\/u\/video\/SBS\/managed\/images\/2016\/07\/11\/722350147591_07110707_image074057_carouselsmall.jpg"},"episodeNumbermpx":"","keywords":["Tour de France 2016 Stage 8","Stage8","Tour de France 2016"],"language":"English","useType":"Interview","tx":"Monday 11 July 7:23am"}], 1, '/cyclingcentral/profiles/sbsdistribution/themes/global/images/on-demand-carousel/no-image.jpg', isPromo, 'http://www.sbs.com.au/api/video_feed/f/Bgtm9B/sbs-search?form=json&range=1-20&byCategories=Sport%2FCycling');

	r := regexp.MustCompile(`SBS.mpxWidget.setVideos\(mpxBeanId, (\[.+\]), \d.+;`)
	subMatches := r.FindSubmatch(body)
	if len(subMatches) != 2 {
		panic("didn't find the mega json array")
	}
	fmt.Printf("%v", string(subMatches[1]))
}
