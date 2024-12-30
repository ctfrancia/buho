package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strings"
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
			r.Use(AuthorizationMiddleware)
			r.Post("/", app.createTournament)
			r.Patch("/{id}", app.updateTournament)
			r.Route("/{id}", func(r chi.Router) {
			})
		})
		r.Route("/auth", func(r chi.Router) {
			r.Post("/token", app.createAuthToken)
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
func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the 'Authorization' header (e.g., 'Bearer <token>')
		authHeader := r.Header.Get("Authorization")

		// Check if the 'Authorization' header is provided
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Check if the token starts with 'Bearer' (common format)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Authorization token format is invalid", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Validate the token (for example, check if it's a predefined valid token)
		if !isValidToken(token) {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// If the token is valid, pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}

// Example function to validate the token (you can replace this with your actual logic)
func isValidToken(token string) bool {
	// FIXME: this is for testing only
	return true
}
