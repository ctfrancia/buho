package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
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
			app.invalidCredentialsResponse(w, r)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Authorization token format is invalid", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Validate the token (for example, check if it's a predefined valid token)
		apiRequester, err := isValidToken(app.config.auth.publicKeyPath, token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), auth.TournamentAPIRequesterKey, apiRequester)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FIXME: This function should be returning a type not an interface
func isValidToken(publicKeyPath, token string) (map[string]interface{}, error) {
	// Load public key
	publicKeyFile, err := os.ReadFile(publicKeyPath)
	if err != nil {
		// Handle error
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyFile)
	if err != nil {
		// Handle error
	}
	// Verify token
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Validate the signing method is RS256
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		if claims["sub"] == nil {
			return nil, fmt.Errorf("invalid token or claims")
		}
		return claims["sub"].(map[string]interface{}), nil
	} else {
		return nil, fmt.Errorf("invalid token or claims")
	}
}
