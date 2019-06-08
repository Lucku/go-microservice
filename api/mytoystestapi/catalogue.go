package mytoystestapi

import (
	"github.com/spf13/viper"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// TODO read from some config
const apiURL = "https://mytoysiostestcase1.herokuapp.com/api/navigation"

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

// GetCatalogue queries the catalogue API and marshals its data
func GetCatalogue() (*Response, error) {

	apiKey := viper.GetString("apiKey")

	if apiKey == "" {
		return nil, errors.New("Error due to undefined api key in the config file")
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("x-api-key", apiKey)

	if err != nil {
		log.Fatal(err)
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