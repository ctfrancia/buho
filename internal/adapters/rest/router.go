package rest

import (
	"github.com/ctfrancia/buho/internal/adapters/rest/handlers"
	"github.com/ctfrancia/buho/internal/core/ports"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	HealthCheckHandler ports.HealthCheckHandler
	AuthHandler        ports.AuthHandler
	TournamentHandler  ports.TournamentHandler
	MatchHandler       ports.MatchHandler
	Logger             ports.Logger
}

func NewRouter(hch ports.HealthCheckService, as ports.AuthService, ms ports.MatchService, l ports.Logger) *chi.Mux {
	httpResponses := handlers.NewHandlerResponse(l)

	router := &Router{
		HealthCheckHandler: handlers.NewHealthCheckHandler(hch),
		AuthHandler:        handlers.NewAuthHandler(as, httpResponses),
		MatchHandler:       handlers.NewMatchHandler(ms),
		Logger:             l,
	}

	return router.setupRoutes()
}

func (r *Router) setupRoutes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Route("/v1", func(router chi.Router) {
		mux.Get("/healthcheck", r.HealthCheckHandler.Handle)
		// Tournaments
		mux.Route("/tournaments", func(r chi.Router) {
			// r.Get("/", tournamentHandler.ListTournaments)
			//mux.Post("/", tournamentHandler.CreateTournament)
		})
		// Auth
		router.Route("/auth", func(r chi.Router) {
			// mux.Post("/login", authHandler.Login)
			// mux.Post("/refresh", auth.refresh)
			// mux.Route("/new", func(r chi.Router) {
			// r.Post("/consumer", app.newApiConsumer)
			// })
		})
	})

	return mux
}
