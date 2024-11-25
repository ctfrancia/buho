package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *application) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheck)
		r.Route("/users", func(r chi.Router) {
			r.Get("/", app.listUsers)
			r.Post("/", app.createUser)
			r.Get("/email/{email}", app.showUserByEmail)
			r.Get("/search", app.searchUsers)
		})
		r.Route("/tournaments", func(r chi.Router) {
			r.Post("/", app.createTournament)
			r.Patch("/{id}", app.updateTournament)
			r.Route("/{id}", func(r chi.Router) {
			})
		})
	})

	// Print out all routes
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		app.logger.Info("route", "method", method, "route", route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		app.logger.Error("Logging err", "err", err)
	}

	return r
}
