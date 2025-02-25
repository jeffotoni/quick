// Uploading multiple files using Quick is simple,
// and much more minimalistic than the patterns we've
// seen in other frameworks.
//
// $ curl -v -X POST http://localhost:8080/upload -F "file=@quicktest.go"
package main

import (
	"fmt"

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
		// c.FormFileLimit("10MB")

		uploadedFile, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(Msg{
				Msg:   "Upload error",
				Error: err.Error(),
			})
		}

		// saving file to disk
		uploadedFile.Save("/tmp/uploads")

		// or Can do It

		// saving file to disk
		uploadedFile.Save("/tmp/uploads", "codeJeff.go")

		// accessing file upload objects

		fmt.Println("Name:", uploadedFile.FileName())
		fmt.Println("Size:", uploadedFile.Size())
		fmt.Println("Content-Type:", uploadedFile.ContentType())
		fmt.Println("Bytes:", len(uploadedFile.Bytes())) // remove the len to get the bytes

		// or
		fmt.Println("Name:", uploadedFile.Info.Filename)
		fmt.Println("Size:", uploadedFile.Info.Size)
		fmt.Println("Content-Type:", uploadedFile.Info.ContentType)
		fmt.Println("Bytes:", len(uploadedFile.Info.Bytes)) // remove the len to get the bytes

		//
		// Respond with success
		return c.Status(200).JSONIN(uploadedFile)
	})

	q.Listen("0.0.0.0:8080")
}
