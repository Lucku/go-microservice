package toysapi

import "github.com/go-chi/chi"

// Routes initializes the router for the toys API
func Routes() *chi.Mux {

	router := chi.NewRouter()
	t := NewToysAPI()

	router.Get("/", t.GetLinks)
	return router
}
