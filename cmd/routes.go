package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ctfrancia/buho/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
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
			r.Use(app.authorizationMiddleware)
			r.Post("/poster", app.uploadTournamentPoster)
			r.Patch("/{id}", app.updateTournament)
			r.Route("/{id}", func(r chi.Router) {
			})
		})
		r.Route("/auth", func(r chi.Router) {
			r.Post("/token", app.createAuthToken)
			r.Post("/new", app.newApiUser)
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
		apiRequester, err := isValidToken(token, app.config.auth.secretKey)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), auth.TournamentAPIRequesterKey, apiRequester)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Example function to validate the token (you can replace this with your actual logic)
func isValidToken(token, secret string) (string, error) {
	// Parse the token with the HMAC signing method
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", fmt.Errorf("Error parsing token: %w", err)
	}

	// If valid, extract and print the claims
	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		if claims["sub"] == nil {
			return "", fmt.Errorf("Invalid token or claims")
		}
		return claims["sub"].(string), nil
	} else {
		return "", fmt.Errorf("Invalid token or claims")
	}
}
