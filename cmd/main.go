package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ctfrancia/buho/internal/adapters/repository"
	"github.com/ctfrancia/buho/internal/adapters/rest"
	"github.com/ctfrancia/buho/internal/core/services"
)

func main() {
	env := os.Getenv("ENV")
	serverConfig, err := repository.NewConfig(env)
	if err != nil {
		os.Exit(1)
	}

	// create database
	db, err := repository.NewDatabase()
	if err != nil {
		os.Exit(1)
	}

	// create secondary adapters
	authStore := repository.NewGormAuthRepository(db)

	// create primary services
	hcs := services.NewHealthCheckService()
	// consumerService := services.NewConsumerService(authStore)

	routes := rest.NewRouter(hcs)

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
