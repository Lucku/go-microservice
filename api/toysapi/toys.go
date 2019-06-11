package toysapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/apex/log"

	"github.com/go-chi/render"
	"github.com/lucku/otto-coding-challenge/api/catalogueapi"
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

// ToysAPI is wrapping the functionality of the ToysAPI
type ToysAPI struct {
	catalogue catalogueapi.Catalogue
}

// NewToysAPI returns a new instance of the toys api
func NewToysAPI() *ToysAPI {
	return &ToysAPI{catalogueapi.CatalogueImpl{}}
}

// GetLinks outputs all the links in an unmodified format
func (t ToysAPI) GetLinks(w http.ResponseWriter, r *http.Request) {

	data, err := t.catalogue.RequestCatalogue(catalogueapi.APIURL)

	if err != nil {
		log.Errorf("Error retrieving data from Catalogue API: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	entries, err := t.handleParentParam(data, w, r)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err := t.handleSortParam(entries, w, r); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	render.JSON(w, r, entries)
}

func (t ToysAPI) handleParentParam(input *catalogueapi.Response, w http.ResponseWriter, r *http.Request) ([]CategoryEntry, error) {

	parent := r.URL.Query().Get("parent")

	leafs := make([]CategoryEntry, 0)

	n := catalogueapi.NavigationEntry{Children: input.NavigationEntries}

	if parent == "" {
		traverseCatalogue(&n, &leafs, parent, true)
	} else {
		traverseCatalogue(&n, &leafs, parent, false)
	}

	if len(leafs) == 0 {
		log.WithFields(log.Fields{
			"parameter": parent,
		}).Debugf("Invalid 'parent' parameter")

		return nil, errors.New("Invalid parent argument")
	}

	return leafs, nil
}

func (t ToysAPI) handleSortParam(input []CategoryEntry, w http.ResponseWriter, r *http.Request) error {

	sortParam := r.URL.Query().Get("sort")

	if sortParam == "" {
		return nil
	}

	sortArgs := strings.Split(sortParam, ",")

	if len(sortArgs) > 2 {
		log.WithFields(log.Fields{
			"numArgs": len(sortArgs),
		}).Debugf("Too many arguments for sorting")

		return errors.New("Invalid numbers of sorting arguments")
	}

	return doSort(input, sortArgs)
}

func traverseCatalogue(entry *catalogueapi.NavigationEntry, leafs *[]CategoryEntry, parent string, found bool) {

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
