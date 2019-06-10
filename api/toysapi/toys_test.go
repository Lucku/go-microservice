package toysapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lucku/otto-coding-challenge/api/catalogueapi"
	"github.com/stretchr/testify/assert"
)

/*
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}
*/

type CatalogueMock struct {
	content catalogueapi.Response
}

func (m CatalogueMock) RequestCatalogue() (*catalogueapi.Response, error) {
	return &m.content, nil
}

func createNavigationEntry(typename, label, url string, fields ...catalogueapi.NavigationEntry) catalogueapi.NavigationEntry {
	return catalogueapi.NavigationEntry{TypeName: typename, Label: label, URL: url, Children: fields}
}

func TestGetLinks(t *testing.T) {

	nLink := createNavigationEntry("link", "Testlink", "www.testlink.de")
	nNode := createNavigationEntry("node", "Testnode", "", nLink)
	nSec := createNavigationEntry("section", "Testsection", "", nNode)

	r := catalogueapi.Response{NavigationEntries: []catalogueapi.NavigationEntry{nSec}}

	m := CatalogueMock{r}

	mToys := ToysAPI{m}

	request, _ := http.NewRequest("GET", "/links", nil)
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(mToys.GetLinks)
	handler.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")

	expected := `[{"label":"Testsection - Testnode - Testlink","url":"www.testlink.de"}]`
	// Cut the trailing newline artificially added by the json encoder
	actual := strings.Trim(response.Body.String(), "\n")

	assert.Equal(t, expected, actual, "API Response is incorrect")
}

func TestGetLinksParent(t *testing.T) {

}

func TestGetLinksEmptyParent(t *testing.T) {

	request, _ := http.NewRequest("GET", "/links?parent", nil)
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(NewToysAPI().GetLinks)
	handler.ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code, "Bad request status expected")
}

func TestGetLinksWrongParent(t *testing.T) {

	request, _ := http.NewRequest("GET", "/links?parent=test", nil)
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(NewToysAPI().GetLinks)
	handler.ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code, "Bad request status expected")
}
