package catalogueapi

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	viper.Set("apikey", "hz7JPdKK069Ui1TRxxd1k8BQcocSVDkj219DVzzD")
}

func TestGetCatalogue(t *testing.T) {

	c := CatalogueImpl{}
	resp, err := c.RequestCatalogue()

	assert.NoError(t, err, "The API could not be called")
	assert.NotNil(t, resp, "Response from API is empty")
	assert.NotEmpty(t, resp.NavigationEntries, "The API outputs contents are empty")
}
