package server

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {

	r := Routes()

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		assert.Contains(t, route, "links", "The route for the links api is not correctly mounted")
		return nil
	}

	err := chi.Walk(r, walkFunc)

	assert.NoError(t, err, "Routes should be properly setup and traversal not result in an error")
}
