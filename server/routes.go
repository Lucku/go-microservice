package server

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/lucku/otto-coding-challenge/api/toysapi"
)

// Routes initializes the global routing and sets the middlewares
func Routes() *chi.Mux {

	router := chi.NewRouter()
	router.Use(
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		middleware.Timeout(time.Second*10),
	)

	router.Route("/", func(r chi.Router) {
		r.Mount("/links", toysapi.Routes())
	})

	return router
}
