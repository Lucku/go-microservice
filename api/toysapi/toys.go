package toysapi

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
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

// GetLinks outputs all the links in an unmodified format
func GetLinks(w http.ResponseWriter, r *http.Request) {

	data, err := mytoystestapi.GetCatalogue()

	if err != nil {
		log.Fatal(err)
	}

	leafs := make([]CategoryEntry, 0)

	for _, e := range data.NavigationEntries {
		traverseCatalogue(e, &leafs)
	}

	render.JSON(w, r, leafs)
}

func traverseCatalogue(entry mytoystestapi.NavigationEntry, leafs *[]CategoryEntry) {

  if entry.TypeName == "link" {
    *leafs = append(*leafs, CategoryEntry{entry.Label, entry.URL})
  }

  if len(entry.Children) > 0 {
		for _, subentry := range entry.Children {

      modEntry := mytoystestapi.NavigationEntry{
        TypeName: subentry.TypeName, 
        Label: entry.Label + " - " + subentry.Label, 
        URL: subentry.URL,
        Children: subentry.Children}
      
      traverseCatalogue(modEntry, leafs)
    }

		//traverseCatalogue(entry.Children[len(entry.Children)-1], leafs)
	}
}

/*
func traverseLinks(entry NavigationEntry, leafs *[]NavigationEntry) {

	if entry.TypeName == "link" {
		*leafs = append(*leafs, entry)
	}

	if len(entry.Children) > 0 {
		for _, subentry := range entry.Children {
			traverseLinks(subentry, leafs)
		}
		traverseLinks(entry.Children[len(entry.Children)-1], leafs)
	}
}
*/