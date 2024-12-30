package main

import (
	"fmt"
	"io"
	// "mime/multipart"
	"net/http"
	"os"
)

func (app *application) createTournament(w http.ResponseWriter, r *http.Request) {
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

	// TODO: this should be moved to the middleware!! And should be attached in the request context
	user := r.Header.Get("x-api-key")
	fmt.Println("user", user)

	metaData := r.MultipartForm.File["tournament"][0]
	fmt.Println("123", metaData.Filename)
	defer file.Close()

	// TODO: BELOW NEEDS TO BE SEFT TO SFTP
	outFile, err := os.Create(metaData.Filename) // Save the file to the server (e.g., in the current directory with a fixed name)
	if err != nil {
		http.Error(w, "Unable to save the image", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Copy the content of the uploaded file to the new file
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	env := envelope{
		"status":   "uploaded",
		"filename": "tournament.jpg",
	}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateTournament(w http.ResponseWriter, r *http.Request) {
}
