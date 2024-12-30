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

	// FIXME: THIS CHECK BELOW IS CREATING A BUG WITH IMAGE SAVING!
	if !isImage(file) {
		app.badRequestResponse(w, r, fmt.Errorf("invalid file type"))
		return
	}

	defer file.Close()

	// TODO: BELOW NEEDS TO BE SEFT TO SFTP
	outFile, err := os.Create("./uploaded_image.jpg") // Save the file to the server (e.g., in the current directory with a fixed name)
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
func isImage(file io.Reader) bool {
	// Create a buffer to hold the first 512 bytes
	buffer := make([]byte, 512)

	// Read the first 512 bytes from the file (but don't consume the whole file)
	_, err := io.ReadFull(file, buffer)
	if err != nil {
		fmt.Println("Error reading file for MIME type detection:", err)
		return false
	}

	// Detect content type based on the first 512 bytes
	contentType := http.DetectContentType(buffer)
	fmt.Println(contentType)
	// Check if the MIME type is an image (JPEG or PNG)
	return contentType == "image/jpeg" || contentType == "image/png"
}

func (app *application) updateTournament(w http.ResponseWriter, r *http.Request) {
}
