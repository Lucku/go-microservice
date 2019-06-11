package toysapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertOrder(t *testing.T, input []CategoryEntry, expected ...CategoryEntry) {

	for i, entry := range expected {
		assert.Equal(t, entry, input[i], "Wrong order of 'sorted' entries")
	}
}

func testSort(t *testing.T, input []CategoryEntry, criteria ...string) {
	err := doSort(input, criteria)
	assert.NoError(t, err, "The sorting should have worked correctly")
}

func TestSortInvalidArg(t *testing.T) {

	request, _ := http.NewRequest("GET", "/links?sort=test", nil)
	response := httptest.NewRecorder()

	err := NewToysAPI().handleSortParam(nil, response, request)

	assert.Error(t, err, "Error due to invalid attribute name expected")
	assert.EqualError(t, err, "Sorting arguments have invalid format")
}

func TestSortOneArg(t *testing.T) {

	c1 := CategoryEntry{"e", "www.d.de"}
	c2 := CategoryEntry{"a", "www.c.de"}
	c3 := CategoryEntry{"d", "www.e.de"}
	c4 := CategoryEntry{"b", "www.b.de"}
	c5 := CategoryEntry{"c", "www.a.de"}
	c := []CategoryEntry{c1, c2, c3, c4, c5}

	testSort(t, c, "label:asc")
	assertOrder(t, c, c2, c4, c5, c3, c1)

	testSort(t, c, "label:desc")
	assertOrder(t, c, c1, c3, c5, c4, c2)

	testSort(t, c, "label")
	assertOrder(t, c, c1, c3, c5, c4, c2)

	testSort(t, c, "url:asc")
	assertOrder(t, c, c5, c4, c2, c1, c3)

	testSort(t, c, "url:desc")
	assertOrder(t, c, c3, c1, c2, c4, c5)

	testSort(t, c, "url")
	assertOrder(t, c, c3, c1, c2, c4, c5)
}

func TestSortTwoArgs(t *testing.T) {

	c1 := CategoryEntry{"3", "1"}
	c2 := CategoryEntry{"4", "5"}
	c3 := CategoryEntry{"5", "4"}
	c4 := CategoryEntry{"5", "1"}
	c5 := CategoryEntry{"2", "3"}
	c6 := CategoryEntry{"6", "4"}
	c7 := CategoryEntry{"3", "2"}
	c8 := CategoryEntry{"1", "7"}
	c := []CategoryEntry{c1, c2, c3, c4, c5, c6, c7, c8}

	testSort(t, c, "url", "label")
	assertOrder(t, c, c8, c2, c6, c3, c5, c7, c4, c1)

	testSort(t, c, "url:asc", "label")
	assertOrder(t, c, c4, c1, c7, c5, c6, c3, c2, c8)

	testSort(t, c, "url", "label:asc")
	assertOrder(t, c, c8, c2, c3, c6, c5, c7, c1, c4)

	testSort(t, c, "url:asc", "label:asc")
	assertOrder(t, c, c1, c4, c7, c5, c3, c6, c2, c8)

	testSort(t, c, "label", "url:desc")
	assertOrder(t, c, c6, c3, c4, c2, c7, c1, c5, c8)

	testSort(t, c, "label:asc", "url:asc")
	assertOrder(t, c, c8, c5, c1, c7, c2, c4, c3, c6)

	testSort(t, c, "label:desc", "url:asc")
	assertOrder(t, c, c6, c4, c3, c2, c1, c7, c5, c8)

	testSort(t, c, "label:desc", "url:desc")
	assertOrder(t, c, c6, c3, c4, c2, c7, c1, c5, c8)
}

func TestSortSameKeyTwice(t *testing.T) {

	c := []CategoryEntry{}

	err := doSort(c, []string{"url:desc", "url:asc"})

	assert.Error(t, err, "Error due to duplicate attribute name expected")

	err = doSort(c, []string{"label", "label"})

	assert.Error(t, err, "Error due to duplicate attribute name expected")
	assert.EqualError(t, err, "Invalid sorting criteria")
}

func TestSortInvalidOrder(t *testing.T) {

	c := []CategoryEntry{}

	err := doSort(c, []string{"url:test"})

	assert.Error(t, err, "Error due to duplicate attribute name expected")
	assert.EqualError(t, err, "Invalid sorting criteria")
}
