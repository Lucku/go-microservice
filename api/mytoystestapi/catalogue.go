package mytoystestapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// TODO read from some config
const apiKey = "hz7JPdKK069Ui1TRxxd1k8BQcocSVDkj219DVzzD"
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

// GetAllLinks consumes the catalogue API and processes its results by filtering out all link entries in
// the JSON tree of entries
func GetAllLinks() ([]NavigationEntry, error) {

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

	var data Response

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, errors.Wrap(err, "Error due to invalid catalogue data format")
	}

	leafs := make([]NavigationEntry, 0)

	for _, e := range data.NavigationEntries {
		traverseLinks(e, &leafs)
	}

	return leafs, nil
}

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
