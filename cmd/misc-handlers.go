package main

import (
	"fmt"
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

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createAuthToken(w http.ResponseWriter, r *http.Request) {
	// Get the IP address of the client
	ip := r.RemoteAddr

	// TODO: check database for Authroized IP addresses
	if ip == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing IP address"))
	}

	tokenString, err := app.auth.CreateJWT("USER_ID")
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	env := envelope{
		"token": tokenString,
	}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
