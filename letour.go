package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

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

func getLinks() *Links {
	return &Links{
		{"Tour de France Stage 2 Highlights", "http://videocdn.sbs.com.au/u/video/SBS_Production/managed/2014/07/07/297830467878_1500K.mp4"},
		{"Tour de France Stage 1 Highlights", "http://videocdn.sbs.com.au/u/video/SBS_Production/managed/2014/07/06/297182787904_1500K.mp4"},
	}
}
