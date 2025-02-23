package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ctfrancia/buho/internal/auth"
	"github.com/ctfrancia/buho/internal/model"
	"github.com/ctfrancia/buho/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type creatingTournamentConsumer struct {
	Website string
	Email   string
	ID      int
}

func (app *application) updateTournament(w http.ResponseWriter, r *http.Request) {
	// userCtx := r.Context().Value(auth.TournamentAPIRequesterKey).(map[string]interface{})
	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing uuid"))
		return
	}
	var t repository.Tournament

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.repository.Tournaments.UpdateByUUID(uuid, t)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Get the updated tournament to send back to the client TODO: I don't think this is necessary
	t, err = app.repository.Tournaments.GetByUUID(uuid)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"tournament": t,
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createTournament(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(auth.TournamentAPIRequesterKey).(map[string]interface{})
	ctc := creatingTournamentConsumer{
		// Website: userCtx["website"].(string),
		// Email:   userCtx["email"].(string),
		ID: int(userCtx["id"].(float64)),
	}

	var ctr model.CreateTournamentRequest
	err := json.NewDecoder(r.Body).Decode(&ctr)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tournamentUUID := uuid.New().String()
	// TODO: Add validation for the tournament dates
	// Build the tournament object for the DB
	sd, err := time.Parse(time.RFC3339, ctr.StartDate)
	if err != nil {
		fmt.Println("error parsing time", err)
		app.badRequestResponse(w, r, fmt.Errorf("failed to parse start date"))
		return
	}

	ed, err := time.Parse(time.RFC3339, ctr.EndDate)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("failed to parse end date"))
		return
	}

	t := repository.Tournament{
		Name:           ctr.Name,
		TournamentUUID: tournamentUUID,
		StartDate:      sd, // TODO: Can't be in the past
		EndDate:        ed, // TODO: End date must be before/in the past
		CreatorID:      uint(ctc.ID),
		PosterURL:      "",
	}

	err = app.repository.Tournaments.Create(&t)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"tournament": t,
	}

	err = app.writeJSON(w, http.StatusCreated, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) uploadTournamentPoster(w http.ResponseWriter, r *http.Request) {
	tCreater := r.Context().Value(auth.TournamentAPIRequesterKey).(map[string]interface{})
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		fmt.Println("error parsing form", err)
		app.badRequestResponse(w, r, err)
		return
	}
	// Get the file from the form (ensure your HTML form uses "file" as the field name)
	file, _, err := r.FormFile("tournament")
	if err != nil {
		e := fmt.Errorf("failed to read file: %w", err)
		app.badRequestResponse(w, r, e)
		return
	}

	metaData := r.MultipartForm.File["tournament"][0]
	defer file.Close()

	uploadPath, err := app.sftp.UploadFile(file, metaData.Filename, tCreater["website"].(string))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"status":    "uploaded",
		"filename":  metaData.Filename,
		"file_size": metaData.Size,
		"file_path": uploadPath,
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
