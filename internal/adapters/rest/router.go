package rest

import (
	"net/http"

	"github.com/ctfrancia/buho/internal/ports/primary"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Router struct {
	// tournamentService ports.TournamentService
	// userService       any
	authService primary.ConsumerServicePort
}

func NewRouter(as primary.ConsumerServicePort) *chi.Mux {
	router := &Router{
		// tournamentService: ts,
		// userService:       us,
		authService: as,
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
	consumerHandler := handlers.NewAuthHandler(r.authService)
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

	return mux
}
