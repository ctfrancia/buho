package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ctfrancia/buho/internal/adapters/repository"
	"github.com/ctfrancia/buho/internal/adapters/rest"
	"github.com/ctfrancia/buho/internal/core/services"
)

func main() {
	env := os.Getenv("ENV")
	serverConfig, err := repository.NewServerConfig(env)
	if err != nil {
		log.Fatal(err)
	}

	// create database
	db, err := repository.NewDatabase(serverConfig.DB)
	if err != nil {
		log.Fatal(err)
	}

	// create secondary adapters
	// authStore := repository.NewGormAuthRepository(db.DB)
	// playersStore := repository.NewGormPlayerRepository(db.DB)
	// tournamentStore := repository.NewGormTournamentRepository(db.DB)
	// tournamentStore := repository.NewGormTournamentRepository(db.DB)
	matchStore := repository.NewGormMatchRepository(db.DB)

	// create primary services
	hcs := services.NewHealthCheckService()
	// consumerService := services.NewConsumerService(authStore)
	// tournamentService := services.NewTournamentService(tournamentStore)
	apiClientService := services.NewApiClientService(authStore)
	matchService := services.NewMatchService(matchStore)

	routes := rest.NewRouter(hcs, apiClientService, matchService)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes,
	}

	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println("error starting server: ", err)
		os.Exit(1)
	}
}
