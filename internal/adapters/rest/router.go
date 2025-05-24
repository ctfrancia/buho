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
	Logger             ports.Logger
}

func NewRouter(hch ports.HealthCheckService, as ports.AuthService, l ports.Logger) *chi.Mux {
	router := &Router{
		HealthCheckHandler: handlers.NewHealthCheckHandler(hch),
		AuthHandler:        handlers.NewAuthHandler(as),
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

	// Print out all routes
	/*
		walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			fields := []zap.Field{
				zap.String("method", method),
				zap.String("route", route),
			}
			r.Logger.Info("route", fields...)
			return nil
		}

		if err := chi.Walk(mux, walkFunc); err != nil {
			fields := []zap.Field{
				zap.String("err", err.Error()),
			}
			r.Logger.Error("Walk err", fields...)
		}
	*/

	return mux
}
