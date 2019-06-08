package toysapi

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lucku/otto-coding-challenge/api/mytoystestapi"
)

/*
[
  {
    "label": "Sortiment - Alter - Baby & Kleinkind - 0-6 Monate",
    "url": "http:\/\/www.mytoys.de\/0-6-months\/"
  },
  ....
]
*/

// Response is the returned message of the API
type Response struct {
	CategoryEntries []CategoryEntry
}

// CategoryEntry is an entry in the catalogue, containing the category label and a url
type CategoryEntry struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

func GetLinks(w http.ResponseWriter, r *http.Request) {

	links, err := mytoystestapi.GetAllLinks()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("", links)
}
