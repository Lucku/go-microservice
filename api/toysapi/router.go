package toysapi

import "github.com/go-chi/chi"

func Routes() *chi.Mux {

	router := chi.NewRouter()
	router.Get("/", GetLinks)
	return router
}
