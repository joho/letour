package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/bugsnag/bugsnag-go"
	"github.com/joho/letour/sbs"
	"github.com/joho/prohttphandler"
)

func main() {
	serveWebsite()
}

func serveWebsite() {
	portNumber := os.Getenv("PORT")
	if portNumber == "" {
		portNumber = "5000"
	}
	listenOn := fmt.Sprintf(":%v", portNumber)

	appHandler := prohttphandler.New("public")

	appHandler.ExactMatchFunc("/", func(w http.ResponseWriter, r *http.Request) {
		links := sbs.GetHighlights()
		t, err := template.ParseFiles("views/index.tmpl.html")
		if err == nil {
			t.Execute(w, links)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	var handler http.Handler = appHandler

	bugsnagApiKey := os.Getenv("BUGSNAG_API_KEY")
	if bugsnagApiKey != "" {
		bugsnag.Configure(bugsnag.Configuration{
			APIKey: bugsnagApiKey,
		})
		handler = bugsnag.Handler(appHandler)
	}

	fmt.Printf("Listening on %v", listenOn)
	log.Fatal(http.ListenAndServe(listenOn, handler))
}
