package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/ctfrancia/buho/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheck)
		r.Route("/players", func(r chi.Router) {
			r.Get("/", app.listPlayers)
			r.Post("/", app.createPlayer)
			r.Get("/email/{email}", app.showUserByEmail)
			r.Get("/search", app.searchUsers)
		})
		r.Route("/tournaments", func(r chi.Router) {
			r.Use(app.authorizationMiddleware)
			r.Post("/new", app.createTournament)
			r.Route("/{uuid}", func(r chi.Router) {
				r.Post("/qr", app.uploadQRCode)
				// r.Get("/", app.showTournament)
				r.Post("/poster", app.uploadTournamentPoster)
				r.Patch("/update", app.updateTournament)
			})
		})
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", app.login)
			r.Post("/refresh", app.refresh)
			r.Route("/new", func(r chi.Router) {
				r.Post("/consumer", app.newApiConsumer)
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

// AuthorizationMiddleware checks if the user is authorized by validating the token
func (app *application) authorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.invalidCredentialsCustomResponse(w, r, "Authorization header is required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.invalidCredentialsCustomResponse(w, r, "Authorization token format is invalid")
			return
		}

		token := parts[1]

		consumer, err := auth.VerifyJWTWithED25519(token, app.config.auth.publicKeyPath)
		if err != nil {
			app.invalidCredentialsCustomResponse(w, r, "error verifying token")
			return
		}

		ctx := context.WithValue(r.Context(), auth.TournamentAPIRequesterKey, consumer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
