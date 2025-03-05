package main

import (
	"context"
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
	userCtx := r.Context().Value(auth.TournamentAPIRequesterKey).(model.Subject)
	ctc := creatingTournamentConsumer{
		ID: int(userCtx.ID),
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
	uuid := chi.URLParam(r, "uuid")
	formFileName := "poster"
	tCreater := r.Context().Value(auth.TournamentAPIRequesterKey).(model.Subject)
	// TODO: 10 MB limit is too big, should be 5 MB
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		fmt.Println("error parsing form:", err)
		app.badRequestResponse(w, r, err)
		return
	}

	file, _, err := r.FormFile(formFileName)
	if err != nil {
		e := fmt.Errorf("failed to read file: %w", err)
		app.badRequestResponse(w, r, e)
		return
	}

	metaData := r.MultipartForm.File[formFileName][0]
	defer file.Close()
	fmt.Println("file size:", metaData.Size)
	fmt.Println("file name:", metaData.Filename)

	// Determine content type
	contentType := metaData.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	cancelCtx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Generate a unique object name to prevent overwrites
	// Format: {creator_website}/poster/{poster-uuid}/{unix_nano}_{filename}
	uniqueObjectName := fmt.Sprintf("%s/poster/%s/%d_%s", tCreater.Website, uuid, time.Now().UnixNano(), metaData.Filename)
	uploadedFilePath, err := app.digitalOcean.UploadFile(
		cancelCtx,
		uniqueObjectName,
		file,
		metaData.Size,
		contentType,
	)
	if err != nil {
		fmt.Println("error uploading file:", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	t := repository.Tournament{
		PosterURL: uploadedFilePath,
	}
	// TODO: Update the tournament with the poster URL
	app.repository.Tournaments.UpdateByUUID(uuid, t)

	env := envelope{
		"upload_file": map[string]interface{}{
			"status":    "uploaded",
			"filename":  uniqueObjectName,
			"file_size": metaData.Size,
			"file_path": uploadedFilePath,
		},
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
