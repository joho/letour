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
	portNumber := os.Getenv("PORT")
	if portNumber == "" {
		portNumber = "5000"
	}
	listenOn := fmt.Sprintf(":%v", portNumber)

	appHandler := prohttphandler.New("public")

	appHandler.ExactMatchFunc("/", func(w http.ResponseWriter, r *http.Request) {
		links := sbs.GetLinks()
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
