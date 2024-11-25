package main

import (
	"fmt"
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
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				app.logger.Info("GET /users")
				fmt.Fprintf(w, "GET /users\n")
			})
			// r.Get("/", app.listUsers)
			// r.Post("/", app.createUser)
			r.Route("/{id}", func(r chi.Router) {
				// r.Get("/", app.showUser)
				// r.Put("/", app.updateUser)
				// r.Delete("/", app.deleteUser)
			})
		})
		r.Route("/tournaments", func(r chi.Router) {
			// r.Get("/", app.listTournaments)
			// r.Post("/", app.createTournament)
			r.Route("/{id}", func(r chi.Router) {
				// r.Get("/", app.showTournament)
				// r.Put("/", app.updateTournament)
				// r.Delete("/", app.deleteTournament)
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
