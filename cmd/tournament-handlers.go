package main

import (
	"encoding/json"
	"fmt"
	"github.com/ctfrancia/buho/internal/auth"
	"net/http"
)

type creatingTournamentEntity struct {
	Website string
	Email   string
	ID      int
}

type creatingTournamentRequest struct {
	// Name string `json:"name"`
}

func (app *application) createTournament(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(auth.TournamentAPIRequesterKey).(map[string]interface{})
	cte := creatingTournamentEntity{
		Website: user["website"].(string),
		Email:   user["email"].(string),
		ID:      int(user["id"].(float64)),
	}

	var ctr creatingTournamentRequest
	decoder := json.NewDecoder(r.Body)
	if decoder.Decode(&ctr) != nil {
		app.badRequestResponse(w, r, fmt.Errorf("failed to decode request body"))
		return
	}

	// TODO: START HERE FOR CREATING A TOURNAMENT
	// returned back the tournament id, that id will be used to upload the poster

	fmt.Println("----------", cte)
}

func (app *application) uploadTournamentPoster(w http.ResponseWriter, r *http.Request) {
	// when uploading a file, check the tournament ID and make sure that the id is the same as
	// the requester of the tournament website. ex: tournament ID 1 fetch tournmentByID(1) and check if the
	// creator is the same as the requester!
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
