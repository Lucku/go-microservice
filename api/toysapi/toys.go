package toysapi

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/lucku/otto-coding-challenge/api/mytoystestapi"
)

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
		http.Error(w, http.StatusText(503), 503)
		return
	}

	entries, err := handleParentParam(data, w, r)

	if err != nil {
		// LOG ERROR
		http.Error(w, http.StatusText(400), 400)
		return
	}

	if err := handleSortParam(entries, w, r); err != nil {
		// LOG ERROR
		http.Error(w, http.StatusText(400), 400)
		return
	}

	render.JSON(w, r, entries)
}

func handleParentParam(input *mytoystestapi.Response, w http.ResponseWriter, r *http.Request) ([]CategoryEntry, error) {

	parent := r.URL.Query().Get("parent")

	leafs := make([]CategoryEntry, 0)

	n := mytoystestapi.NavigationEntry{Children: input.NavigationEntries}

	if parent == "" {
		traverseCatalogue(&n, &leafs, parent, true)
	} else {
		traverseCatalogue(&n, &leafs, parent, false)
	}

	if len(leafs) == 0 {
		return nil, errors.New("Invalid parent argument")
	}

	return leafs, nil
}

func traverseCatalogue(entry *mytoystestapi.NavigationEntry, leafs *[]CategoryEntry, parent string, found bool) {

	if entry.TypeName == "link" {

		if found == true || parent == "" {
			*leafs = append(*leafs, CategoryEntry{entry.Label, entry.URL})
		}

		return
	}

	if len(entry.Children) > 0 {
		for _, subentry := range entry.Children {

			if found && entry.Label != "" {
				subentry.Label = entry.Label + " - " + subentry.Label
			}

			match := found

			if !match {
				match = entry.Label == parent
			}

			traverseCatalogue(&subentry, leafs, parent, match)
		}
	}
}
