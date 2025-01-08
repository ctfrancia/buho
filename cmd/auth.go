package main

import (
	"github.com/ctfrancia/buho/internal/auth"
	"github.com/ctfrancia/buho/internal/model"
	"github.com/ctfrancia/buho/internal/repository"
	"github.com/ctfrancia/buho/internal/validator"

	"encoding/json"
	"gorm.io/gorm"
	"net/http"
)

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	// mashal the request body into a struct
	var requestBody model.CreateAuthTokenRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(requestBody.Email != "", "email", "must be provided")
	v.Check(requestBody.Password != "", "password", "must be provided")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	var user repository.AuthModel
	user.Email = requestBody.Email
	err = app.repository.Auth.SelectByEmail(&user)
	if err == gorm.ErrRecordNotFound {
		// TODO: Add a log message here
		app.invalidCredentialsResponse(w, r)
		return
	}

	// Compare the password from the request to the password in the DB
	match, err := auth.CompareHashAndPassword(user.Password, requestBody.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	tokenString, err := app.auth.CreateJWT(user)
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

func (app *application) newApiConsumer(w http.ResponseWriter, r *http.Request) {
	// mashal the request body into a struct
	var requestBody model.NewAPIConsumerRequest

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

	authModelUser := &repository.AuthModel{
		Email:     requestBody.Email,
		FirstName: requestBody.FirstName,
		LastName:  requestBody.LastName,
		Website:   requestBody.Website,
	}

	// check if user is in DB and password if the user is in the DB then return a 409
	err := app.repository.Auth.SelectByEmail(authModelUser)
	if err != gorm.ErrRecordNotFound {
		app.conflictResponse(w, r)
		return
	}

	generatedPW, err := auth.CreateSecretKey(auth.PasswordGeneratorDefaultLength)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Hash the password
	encodedHash, err := auth.Hash(generatedPW)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Assign the argon2 hash to the user password
	authModelUser.Password = encodedHash

	// Create the user in DB
	err = app.repository.Auth.Create(authModelUser)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	/// Return the user with the generated password
	authModelUser.Password = generatedPW
	err = app.writeJSON(w, http.StatusCreated, envelope{"consumer": authModelUser}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
