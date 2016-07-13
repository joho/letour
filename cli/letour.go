package main

import (
	"fmt"

	"github.com/joho/letour/sbs"
)

func main() {
	feed := sbs.GetHighlights()
	for _, v := range feed {
		fmt.Printf("%v\n%v\n", v.Title, v.Url())
	}
}
