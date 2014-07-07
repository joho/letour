package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		links := getLinks()
		t, err := template.ParseFiles("views/index.tmpl.html")
		if err == nil {
			t.Execute(w, links)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}).Methods("GET")

	r.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "public/styles.css")
	}).Methods("GET")

	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		http.ServeFile(w, r, "public/favicon.ico")
	}).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3000", nil))
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
