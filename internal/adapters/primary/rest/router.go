package rest

import (
	"net/http"

	"github.com/ctfrancia/buho/internal/adapters/primary/rest/handlers"
	"github.com/ctfrancia/buho/internal/core/ports"
	"github.com/ctfrancia/buho/internal/core/ports/primary"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	// "go.uber.org/zap"
)

type Router struct {
	tournamentService ports.TournamentService
	userService       any
	authService       primary.AuthPort
}

func NewRouter(ts ports.TournamentService, us any, as primary.AuthPort) *chi.Mux {
	router := &Router{
		tournamentService: ts,
		userService:       us,
		authService:       as,
	}

	return router.setupRoutes()
}

func (r *Router) setupRoutes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	miscHandler := handlers.NewMiscHandler()
	tournamentHandler := handlers.NewTournamentHandler(r.tournamentService)
	authHandler := handlers.NewAuthHandler(r.authService)
	// userHandler := handlers.NewUserHandler(r.userService)

	mux.Route("/v1", func(r chi.Router) {
		mux.Get("/healthcheck", miscHandler.Healthcheck)
		// Tournaments
		mux.Route("/tournaments", func(r chi.Router) {
			// r.Get("/", tournamentHandler.ListTournaments)
			mux.Post("/", tournamentHandler.CreateTournament)
		})
		// Auth
		r.Route("/auth", func(r chi.Router) {
			mux.Post("/login", authHandler.Login)
			// mux.Post("/refresh", auth.refresh)
			// mux.Route("/new", func(r chi.Router) {
			// r.Post("/consumer", app.newApiConsumer)
			// })
		})
	})
	/*
		r.Route("/v1", func(r chi.Router) {
			r.Get("/healthcheck", h.healthcheck)
			r.Route("/players", func(r chi.Router) {
				r.Get("/", app.listPlayers)
				r.Post("/", app.createPlayer)
				r.Get("/email/{email}", app.showUserByEmail)
				r.Get("/search", app.searchUsers)
			})
			r.Route("/tournaments", func(r chi.Router) {
				r.Get("/", app.listTournaments)
				r.Use(app.authorizationMiddleware)
				r.Post("/new", app.createTournament)
				r.Route("/{uuid}", func(r chi.Router) {
					r.Get("/", app.getTournament)
					r.Route("/poster", func(r chi.Router) {
						r.Delete("/", app.deleteTournamentPoster)
						r.Post("/upload", app.uploadTournamentPoster)
					})
					r.Patch("/data", app.updateTournament)
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
	*/

	return mux
}

func routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	/*
		r.Route("/v1", func(r chi.Router) {
			r.Get("/healthcheck", h.healthcheck)
			r.Route("/players", func(r chi.Router) {
				r.Get("/", app.listPlayers)
				r.Post("/", app.createPlayer)
				r.Get("/email/{email}", app.showUserByEmail)
				r.Get("/search", app.searchUsers)
			})
			r.Route("/tournaments", func(r chi.Router) {
				r.Get("/", app.listTournaments)
				r.Use(app.authorizationMiddleware)
				r.Post("/new", app.createTournament)
				r.Route("/{uuid}", func(r chi.Router) {
					r.Get("/", app.getTournament)
					r.Route("/poster", func(r chi.Router) {
						r.Delete("/", app.deleteTournamentPoster)
						r.Post("/upload", app.uploadTournamentPoster)
					})
					r.Patch("/data", app.updateTournament)
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
	*/
	r.Get("/", h.healthcheck)

	// Print out all routes
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fields := []zap.Field{
			zap.String("method", method),
			zap.String("route", route),
		}
		app.logger.Info("route", fields...)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		fields := []zap.Field{
			zap.String("err", err.Error()),
		}
		app.logger.Error("Walk err", fields...)
	}

	return r
}
