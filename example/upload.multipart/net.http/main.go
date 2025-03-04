package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Limite de 10MB para o upload
	r.ParseMultipartForm(10 << 20)

	// Recover the file
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error receiving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Creates a local file to save the upload
	dst, err := os.Create(handler.Filename)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copies the file data to the destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error copying file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Upload successful: %s", handler.Filename)
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	fmt.Println("Server Run Port:8080")
	http.ListenAndServe(":8080", nil)
}
