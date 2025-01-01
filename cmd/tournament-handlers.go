package main

import (
	"fmt"
	"github.com/ctfrancia/buho/internal/auth"
	"net/http"
)

func (app *application) uploadTournamentPoster(w http.ResponseWriter, r *http.Request) {
	tCreater := r.Context().Value(auth.TournamentAPIRequesterKey).(string)
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
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

	uploadPath, err := app.sftp.UploadFile(file, metaData.Filename, tCreater)
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

func (app *application) updateTournament(w http.ResponseWriter, r *http.Request) {
}
