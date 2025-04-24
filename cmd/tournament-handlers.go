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
	"github.com/ctfrancia/buho/internal/validator"
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

	v := validator.New()
	v.Check(sd.After(time.Now()), "start_date", "start_date cannot be in the past")
	v.Check(ed.After(sd), "end_date", "end_date has to be after start_date")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	}

	t := repository.Tournament{
		Name:           ctr.Name,
		TournamentUUID: tournamentUUID,
		StartDate:      sd,
		EndDate:        ed,
		CreatorID:      uint(ctc.ID),
		PosterURL:      ctr.PosterURL,
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

func (app *application) uploadQRCode(w http.ResponseWriter, r *http.Request) {
	tCreater := r.Context().Value(auth.TournamentAPIRequesterKey).(map[string]any)
	fmt.Println("tCreater", tCreater["website"])
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Get the file from the form (ensure your HTML form uses "file" as the field name)
	file, _, err := r.FormFile("qrcode")
	if err != nil {
		e := fmt.Errorf("failed to read file: %w", err)
		app.badRequestResponse(w, r, e)
		return
	}

	metaData := r.MultipartForm.File["qrcode"][0]
	defer file.Close()

	// TODO: updload to digital ocean below
	env := envelope{
		"status":    "uploaded",
		"file_name": metaData.Filename,
		"file_size": metaData.Size,
		// "file_path": uploadPath,
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
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
	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), metaData.Filename)
	// Format: {creator_website}/poster/{poster-uuid}/{fileName}
	uniqueObjectName := fmt.Sprintf("%s/poster/%s/%s", tCreater.Website, uuid, fileName)
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
		"upload_file": map[string]any{
			"status":    "uploaded",
			"file_name": fileName,
			"file_size": metaData.Size,
			"file_path": uploadedFilePath,
		},
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTournamentPoster(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing uuid"))
		return
	}

	fileName := chi.URLParam(r, "file_name")
	if fileName == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing file_name"))
		return
	}

	consumerWebsite := r.Context().Value(auth.TournamentAPIRequesterKey).(model.Subject).Website
	photoPath := fmt.Sprintf("%s/poster/%s/%s", consumerWebsite, uuid, fileName)

	er := app.digitalOcean.DeleteFile(cancelCtx, photoPath)
	if er != nil {
		app.serverErrorResponse(w, r, er)
		return
	}

	err := app.repository.Tournaments.RemoveTournamentPosterURL(uuid)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	env := envelope{
		"remove_file": map[string]any{
			"status":    "deleted",
			"file_name": fileName,
		},
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) downloadTournamentPoster(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing uuid"))
		return
	}
}
