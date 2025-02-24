// Uploading multiple files using Quick is simple,
// and much more minimalistic than the patterns we've
// seen in other frameworks.
//
// $ curl -v -X POST http://localhost:8080/upload -F "file=@quicktest.go"
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/jeffotoni/quick"
)

type Msg struct {
	Msg   string `json:"msg"`
	Error string `json:"error"`
}

func main() {
	// start Quick
	q := quick.New()

	q.Post("/upload", func(c *quick.Ctx) error {

		// Default Size Upload 1MB
		// Set upload limit (10MB)
		c.FormFileLimit("10MB")

		uploadedFile, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(Msg{
				Msg:   "Upload error",
				Error: err.Error(),
			})
		}

		// handling manually, getting size, copy etc.
		file := uploadedFile.Multipart
		fileUp, err := file.Open()
		if err != nil {
			return c.Status(500).SendString("Erro open file multipart.FileHeader")
		}
		defer fileUp.Close()

		// Read file in bytes
		fileBytes, err := io.ReadAll(fileUp)
		if err != nil {
			return c.Status(500).SendString("Erro when reading bytes from the file")
		}

		// Read file in bytes
		fmt.Println("Size:", len(fileBytes))

		// save file path
		dir := "/tmp/upload_multipart/"

		// Creates the directory if it does not exist
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return c.Status(500).SendString("Error creating directory")
		}

		// Creates a local file to save the upload
		pathFile := dir + file.Filename
		dst, err := os.Create(pathFile)
		if err != nil {
			return c.Status(500).SendString("Error saving file")
		}
		defer dst.Close()

		// Copies the file data to the destination
		_, err = io.Copy(dst, fileUp)
		if err != nil {
			return c.Status(500).SendString("Error copying file")
		}

		// Respond with success
		return c.Status(200).JSONIN(uploadedFile)
	})

	q.Listen("0.0.0.0:8080")
}
