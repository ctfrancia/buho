package main

import (
	"net/http"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"version":     version,
			"environment": app.config.env,
		},
	}

	app.writeJSON(w, http.StatusOK, env, nil)
}
