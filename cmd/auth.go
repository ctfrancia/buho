package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ctfrancia/buho/internal/model"
	"github.com/ctfrancia/buho/internal/repository"
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
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := &repository.AuthModel{
		Email: requestBody.Email,
	}
	// TODO: Create validator for this to make sure all necessary fields are present

	// check if user is in DB and password if the user is in the DB then return a 409
	err = app.repository.Auth.SelectByEmail(user)
	if err != gorm.ErrRecordNotFound {
		fmt.Println("111111111111111111111111", err)
		// app.conflictResponse(w, r)
		return
	}
	// START HERE
	fmt.Println("22222222222", user)
	// if the user is not in the DB then create the user and return a 201
	/*

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
	*/
}
