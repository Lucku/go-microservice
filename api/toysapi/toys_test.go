package toysapi

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"

	"github.com/lucku/otto-coding-challenge/api/catalogueapi"
	"github.com/stretchr/testify/assert"
)

type CatalogueMock struct {
	content catalogueapi.Response
}

func (m CatalogueMock) RequestCatalogue(url string) (*catalogueapi.Response, error) {
	return &m.content, nil
}

type ErroneousCatalogue struct {
}

func (ErroneousCatalogue) RequestCatalogue(url string) (*catalogueapi.Response, error) {
	return nil, errors.New("Error")
}

func createNavigationEntry(typename, label, url string, fields ...catalogueapi.NavigationEntry) catalogueapi.NavigationEntry {
	return catalogueapi.NavigationEntry{TypeName: typename, Label: label, URL: url, Children: fields}
}

func init() {
	viper.Set("apiKey", "hz7JPdKK069Ui1TRxxd1k8BQcocSVDkj219DVzzD")
}

func sendRequest(h http.HandlerFunc, requestStr string) *httptest.ResponseRecorder {
	request, _ := http.NewRequest("GET", requestStr, nil)
	response := httptest.NewRecorder()
	h.ServeHTTP(response, request)
	return response
}

func TestGetLinksTwoSections(t *testing.T) {

	nLink1 := createNavigationEntry("link", "link1", "www.link1.de")
	nLink2 := createNavigationEntry("link", "link2", "www.link2.de")
	nNode1 := createNavigationEntry("node", "node1", "", nLink1)
	nNode2 := createNavigationEntry("node", "node2", "", nLink2)
	nSec1 := createNavigationEntry("section", "section1", "", nNode1)
	nSec2 := createNavigationEntry("section", "section2", "", nNode2)

	r := catalogueapi.Response{NavigationEntries: []catalogueapi.NavigationEntry{nSec1, nSec2}}

	m := CatalogueMock{r}

	mToys := ToysAPI{m}

	response := sendRequest(http.HandlerFunc(mToys.GetLinks), "/links")
	assert.Equal(t, 200, response.Code, "OK response is expected")

	expected := `[{"label":"section1 - node1 - link1","url":"www.link1.de"},` +
		`{"label":"section2 - node2 - link2","url":"www.link2.de"}]`

	assert.JSONEq(t, expected, response.Body.String(), "API response is incorrect")
}

func TestGetLinksTwoEntries(t *testing.T) {

	nLink1 := createNavigationEntry("link", "link1", "www.link1.de")
	nLink2 := createNavigationEntry("link", "link2", "www.link2.de")
	nNode := createNavigationEntry("node", "node1", "", nLink1, nLink2)
	nSec := createNavigationEntry("section", "section1", "", nNode)

	r := catalogueapi.Response{NavigationEntries: []catalogueapi.NavigationEntry{nSec}}

	m := CatalogueMock{r}

	mToys := ToysAPI{m}

	response := sendRequest(http.HandlerFunc(mToys.GetLinks), "/links")
	assert.Equal(t, 200, response.Code, "OK response is expected")

	expected := `[{"label":"section1 - node1 - link1","url":"www.link1.de"},` +
		`{"label":"section1 - node1 - link2","url":"www.link2.de"}]`

	assert.JSONEq(t, expected, response.Body.String(), "API Response is incorrect")
}

func TestGetLinksParent(t *testing.T) {

	nLink1 := createNavigationEntry("link", "link1", "www.link1.de")
	nLink2 := createNavigationEntry("link", "link2", "www.link2.de")
	nNode1 := createNavigationEntry("node", "node1", "", nLink1)
	nNode2 := createNavigationEntry("node", "node2", "", nLink2)
	nSec1 := createNavigationEntry("section", "section1", "", nNode1)
	nSec2 := createNavigationEntry("section", "section2", "", nNode2)

	r := catalogueapi.Response{NavigationEntries: []catalogueapi.NavigationEntry{nSec1, nSec2}}
	m := CatalogueMock{r}
	mToys := ToysAPI{m}

	response := sendRequest(http.HandlerFunc(mToys.GetLinks), "/links?parent=section1")
	assert.Equal(t, 200, response.Code, "OK response is expected")
	expected := `[{"label":"node1 - link1","url":"www.link1.de"}]`
	assert.JSONEq(t, expected, response.Body.String(), "API Response is incorrect")

	response = sendRequest(http.HandlerFunc(mToys.GetLinks), "/links?parent=section2")
	assert.Equal(t, 200, response.Code, "OK response is expected")
	expected = `[{"label":"node2 - link2","url":"www.link2.de"}]`
	assert.JSONEq(t, expected, response.Body.String(), "API Response is incorrect")

	response = sendRequest(http.HandlerFunc(mToys.GetLinks), "/links?parent=node1")
	assert.Equal(t, 200, response.Code, "OK response is expected")
	expected = `[{"label":"link1","url":"www.link1.de"}]`
	assert.JSONEq(t, expected, response.Body.String(), "API Response is incorrect")

	response = sendRequest(http.HandlerFunc(mToys.GetLinks), "/links?parent=node2")
	assert.Equal(t, 200, response.Code, "OK response is expected")
	expected = `[{"label":"link2","url":"www.link2.de"}]`
	assert.JSONEq(t, expected, response.Body.String(), "API Response is incorrect")
}

func TestGetLinksInvalidParent(t *testing.T) {
	response := sendRequest(http.HandlerFunc(NewToysAPI().GetLinks), "/links?parent=test")
	assert.Equal(t, http.StatusNotFound, response.Code, "NotFound status expected")
}

func TestGetLinksQueryAPIError(t *testing.T) {
	e := ErroneousCatalogue{}
	response := sendRequest(http.HandlerFunc(ToysAPI{e}.GetLinks), "/links?parent=test")
	assert.Equal(t, http.StatusInternalServerError, response.Code, "InternalServerError status expected")
}

func TestHandleSortTooManyArgs(t *testing.T) {

	nLink := createNavigationEntry("link", "link", "www.link.de")
	nNode := createNavigationEntry("node", "node", "", nLink)
	nSec := createNavigationEntry("section", "section1", "", nNode)

	r := catalogueapi.Response{NavigationEntries: []catalogueapi.NavigationEntry{nSec}}
	m := CatalogueMock{r}
	mToys := ToysAPI{m}

	response := sendRequest(http.HandlerFunc(mToys.GetLinks), "/links?sort=one,two,three")

	assert.Equal(t, http.StatusNotFound, response.Code, "NotFound status expected")
}

func TestGetLinksParentAndSort(t *testing.T) {

	nLink1 := createNavigationEntry("link", "link1", "www.link1.de")
	nLink2 := createNavigationEntry("link", "link2", "www.link2.de")
	nLink3 := createNavigationEntry("link", "link3", "www.link3.de")
	nLink4 := createNavigationEntry("link", "link4", "www.link4.de")
	nNode1 := createNavigationEntry("node", "node1", "", nLink1, nLink3)
	nNode2 := createNavigationEntry("node", "node2", "", nLink4, nLink2)
	nSec := createNavigationEntry("section", "section1", "", nNode1, nNode2)

	r := catalogueapi.Response{NavigationEntries: []catalogueapi.NavigationEntry{nSec}}
	m := CatalogueMock{r}
	mToys := ToysAPI{m}

	response := sendRequest(http.HandlerFunc(mToys.GetLinks), "/links?sort=label")

	expected := `[{"label":"section1 - node2 - link4","url":"www.link4.de"},` +
		`{"label":"section1 - node2 - link2","url":"www.link2.de"},` +
		`{"label":"section1 - node1 - link3","url":"www.link3.de"},` +
		`{"label":"section1 - node1 - link1","url":"www.link1.de"}]`

	assert.JSONEq(t, expected, response.Body.String(), "API Response is incorrect")

	response = sendRequest(http.HandlerFunc(mToys.GetLinks), "/links?parent=node1&sort=url:asc")

	expected = `[{"label":"link1","url":"www.link1.de"},` +
		`{"label":"link3","url":"www.link3.de"}]`

	assert.JSONEq(t, expected, response.Body.String(), "API Response is incorrect")

	response = sendRequest(http.HandlerFunc(mToys.GetLinks), "/links?parent=node2&sort=label,url")

	expected = `[{"label":"link4","url":"www.link4.de"},` +
		`{"label":"link2","url":"www.link2.de"}]`

	assert.JSONEq(t, expected, response.Body.String(), "API Response is incorrect")

	response = sendRequest(http.HandlerFunc(mToys.GetLinks), "/links?parent=section1&sort=url:asc")

	expected = `[{"label":"node1 - link1","url":"www.link1.de"},` +
		`{"label":"node2 - link2","url":"www.link2.de"},` +
		`{"label":"node1 - link3","url":"www.link3.de"},` +
		`{"label":"node2 - link4","url":"www.link4.de"}]`

	assert.JSONEq(t, expected, response.Body.String(), "API Response is incorrect")
}
