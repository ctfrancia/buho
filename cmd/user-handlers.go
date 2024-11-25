package main

import (
	"net/http"

	"github.com/ctfrancia/buho/internal/validator"
	"github.com/go-chi/chi/v5"
)

const (
	searchBy    = "by"
	searchValue = "value"
)

var validSearchParams = []string{"id", "email", "last_name", "first_name"}

func (app *application) showUserByEmail(w http.ResponseWriter, r *http.Request) {
	// Extract the email parameter from the URL
	email := chi.URLParam(r, "email")

	v := validator.New()
	emailValidFormat := validator.Matches(email, validator.EmailRX)

	v.Check(emailValidFormat, "email", "invalid email format")
	v.Check(email != "", "email", "email is required")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	app.logger.Info("GET /users/email", "email", email)
}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("POST /users")
}

func (app *application) listUsers(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("GET /users")
}

func (app *application) searchUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	v := validator.New()
	query := queryParams.Get(searchBy)
	value := queryParams.Get(searchValue)

	v.Check(query != "", searchBy, "type is required in query and value of 'id', 'email', 'last_name', 'first_name' is allowed")
	v.Check(value != "", searchValue, "value is required")
	v.Check(v.In(query, validSearchParams...), "by", "type is not valid: [id, email, last_name, first_name] are allowed")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	app.logger.Info("GET /users/search", "query", query, "value", value)
}
