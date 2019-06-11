package catalogueapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetCatalogue(t *testing.T) {

	viper.Set("apikey", "test")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		assert.Equal(t, req.Header.Get("x-api-key"), "test")

		rw.Write([]byte(`{
			"navigationEntries": [{
				"type": "section",
				"label": "Sortiment"
				}]}`))
	}))

	defer server.Close()

	c := CatalogueImpl{}
	resp, err := c.RequestCatalogue(server.URL)

	assert.NoError(t, err, "The API should have been called correctly")
	assert.NotNil(t, resp, "Response from API shall not be empty")
	assert.NotEmpty(t, resp.NavigationEntries, "The API outputs contents should not be empty")
}

func TestGetCatalogueNoAPIKey(t *testing.T) {

	viper.Set("apiKey", "")

	c := CatalogueImpl{}
	_, err := c.RequestCatalogue(APIURL)

	assert.Error(t, err, "No error although API key missing")
	assert.EqualError(t, err, "Error due to undefined api key in the config file")
}

func TestGetCatalogueInvalidJSON(t *testing.T) {

	viper.Set("apikey", "hz7JPdKK069Ui1TRxxd1k8BQcocSVDkj219DVzzD")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"Invalid":json}`))
	}))

	defer server.Close()

	c := CatalogueImpl{}
	_, err := c.RequestCatalogue(server.URL)

	assert.Error(t, err, "No error due to invalid JSON response occurred")
	assert.Contains(t, err.Error(), "Error due to invalid catalogue data format")
}
