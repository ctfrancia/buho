package main

import (
	"encoding/json"
	"net/http"
)

type CreateAuthTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type CreateAuthTokenResponse struct {
	Token string `json:"token"`
}

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
	// mashal the request body into a struct
	var requestBody CreateAuthTokenRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// check if user is in DB and password is correct (this is a mock)
	if requestBody.Email != "foo" || requestBody.Password != "bar" {
		app.invalidCredentialsResponse(w, r)
		return
	}

	tokenString, err := app.auth.CreateJWT(requestBody.Email)
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
