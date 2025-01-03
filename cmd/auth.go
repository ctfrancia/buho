package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ctfrancia/buho/internal/auth"
	"github.com/ctfrancia/buho/internal/model"
	"github.com/ctfrancia/buho/internal/repository"
	"github.com/ctfrancia/buho/internal/validator"
	"gorm.io/gorm"
)

func (app *application) createAuthToken(w http.ResponseWriter, r *http.Request) {
	// mashal the request body into a struct
	var requestBody model.CreateAuthTokenRequest
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

func (app *application) newApiUser(w http.ResponseWriter, r *http.Request) {
	// mashal the request body into a struct
	var requestBody model.NewAPIUserRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the request body that the required fields are present
	v := validator.New()
	v.Check(requestBody.Email != "", "email", "must be provided")
	v.Check(requestBody.FirstName != "", "first_name", "must be provided")
	v.Check(requestBody.LastName != "", "last_name", "must be provided")
	v.Check(requestBody.Website != "", "website", "must be provided")
	v.Check(validator.Matches(requestBody.Email, validator.EmailRX), "email", "must be a valid email address")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user := &repository.AuthModel{
		Email: requestBody.Email,
	}

	// check if user is in DB and password if the user is in the DB then return a 409
	err := app.repository.Auth.SelectByEmail(user)
	if err != gorm.ErrRecordNotFound {
		app.conflictResponse(w, r)
		return
	}

	pw, err := auth.CreateSecretKey(auth.PasswordGeneratorDefaultLength)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Assign the password to the user
	user.Password = pw

	err = app.repository.Auth.Create(*user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
