package catalogueapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// APIURL is the url to the 'catalogue API' to be queried
const APIURL = "https://mytoysiostestcase1.herokuapp.com/api/navigation"

// Response is the root of the JSON message of the catalogue
type Response struct {
	NavigationEntries []NavigationEntry `json:"navigationEntries"`
}

// NavigationEntry is a composite struct that either contains children of itself or acts as a leave
// In case it is a leave node, it contains the URL attribute
type NavigationEntry struct {
	TypeName string            `json:"type"`
	Label    string            `json:"label"`
	URL      string            `json:"url,omitempty"`
	Children []NavigationEntry `json:"children"`
}

// Catalogue is an interface to the consumed API
type Catalogue interface {
	RequestCatalogue(string) (*Response, error)
}

// CatalogueImpl is a concrete implementation of the catalogue interface
type CatalogueImpl struct {
}

// RequestCatalogue queries the catalogue API and unmarshals its data into a Response struct
func (CatalogueImpl) RequestCatalogue(url string) (*Response, error) {

	apiKey := viper.GetString("apiKey")

	if apiKey == "" {
		return nil, errors.New("Error due to undefined api key in the config file")
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("x-api-key", apiKey)

	if err != nil {
		return nil, errors.Wrap(err, "Error creating the Http request")
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, errors.Wrap(err, "Error consuming the catalogue data")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.Wrap(err, "Error while processing catalogue response")
	}

	data := new(Response)

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, errors.Wrap(err, "Error due to invalid catalogue data format")
	}

	return data, nil
}
