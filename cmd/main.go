package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ctfrancia/buho/internal/adapters/repository"
	"github.com/ctfrancia/buho/internal/adapters/rest"
	"github.com/ctfrancia/buho/internal/core/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	env := os.Getenv("ENV")
	serverConfig, err := repository.NewConfig(env)
	if err != nil {
		os.Exit(1)
	}

	db, err := openDB(serverConfig)
	if err != nil {
		os.Exit(1)
	}

	// create secondary adapters
	authStore := repository.NewGormAuthRepository(db)

	// create primary services
	consumerService := services.NewConsumerService(authStore)

	routes := rest.NewRouter(consumerService)

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

func openDB(cfg *repository.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DB.DSN), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}

	// db.Debug() for debugging
	// db.AutoMigrate(&repository.Auth{}, &repository.Organizer{}, &repository.Location{}, &repository.Tournament{})

	return db, nil
}
