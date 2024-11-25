package main

import (
	"net/http"

	"github.com/ctfrancia/buho/internal/model"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	data := model.HealthCheck{Status: "available", Version: "1.0.0", Environment: app.config.env}

	app.writeJSON(w, http.StatusOK, envelope{"info": data}, nil)
}
