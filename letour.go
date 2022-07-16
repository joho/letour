package main

import (
	"flag"
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
	debug := flag.Bool("debug", false, "spit out debugging info on what SBS serves")
	flag.Parse()
	if *debug {
		debugAPI()
	} else {
		serveWebsite()
	}
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
			fmt.Println("200 OK")
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	var handler http.Handler = appHandler

	bugsnagAPIKey := os.Getenv("BUGSNAG_API_KEY")
	if bugsnagAPIKey != "" {
		bugsnag.Configure(bugsnag.Configuration{
			APIKey: bugsnagAPIKey,
		})
		handler = bugsnag.Handler(appHandler)
	}

	fmt.Printf("Listening on %v\n", listenOn)
	log.Fatal(http.ListenAndServe(listenOn, handler))
}

func debugAPI() {
	fmt.Println("Matching Highlights")
	for _, video := range sbs.GetHighlights() {
		fmt.Printf("%#v\n", video)
	}
}
