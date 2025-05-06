package main

import (
	"net/http"

	"go.uber.org/zap"
)

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errs map[string]string) {
	env := envelope{"errors": errs}

	app.errorResponse(w, r, http.StatusUnprocessableEntity, env)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	fields := []zap.Field{
		zap.String("method", method),
		zap.String("uri", uri),
	}
	app.logger.Error(err.Error(), fields...)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) invalidCredentialsCustomResponse(w http.ResponseWriter, r *http.Request, message string) {
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "a record already exists with this email address"
	app.errorResponse(w, r, http.StatusConflict, message)
}
