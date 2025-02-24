// Uploading multiple files using Quick is simple,
// and much more minimalistic than the patterns we've
// seen in other frameworks.
//
// $ curl -v -X POST http://localhost:8080/upload -F "files=@quicktest.go" -F "files=@quick_test.go"
package main

import (
	"fmt"

	"github.com/jeffotoni/quick"
)

type Msg struct {
	Msg   string `json:"msg"`
	Error string `json:"error"`
}

// Simulated S3 upload function
func uploadToS3(bucket, name, contentType string, size int64, data []byte) error {
	// Simulating the upload process
	fmt.Printf("Uploading to S3...\n")
	fmt.Printf("Bucket: %s\n", bucket)
	fmt.Printf("File Name: %s\n", name)
	fmt.Printf("Content Type: %s\n", contentType)
	fmt.Printf("Size: %d bytes\n", size)
	fmt.Println("Upload successful!")
	return nil
}

func main() {
	// start Quick
	q := quick.New()

	q.Post("/upload", func(c *quick.Ctx) error {

		// Default Size Upload 1MB
		// Set upload limit (10MB)
		err := c.FormFileLimit("10MB")
		if err != nil {
			return c.Status(500).JSON(Msg{
				Msg:   "Error Limit",
				Error: err.Error(),
			})
		}

		// Retrieve multiple files
		uploadedFiles, err := c.FormFiles("files")
		if err != nil {
			return c.Status(400).JSON(Msg{
				Msg:   "Upload error",
				Error: err.Error(),
			})
		}

		fmt.Println("Received", len(uploadedFiles), "files")

		// Simulated S3 bucket name
		bucket := "my-awesome-bucket"

		// Iterate over uploaded files
		for _, file := range uploadedFiles {
			// Simulating upload to S3
			err := uploadToS3(
				bucket,
				file.Info.Filename,
				file.Info.ContentType,
				file.Info.Size,
				file.Info.Bytes,
			)
			if err != nil {

				return c.Status(500).JSON(Msg{
					Msg:   "Error uploading to S3",
					Error: err.Error(),
				})
			}
		}

		// uploadedFiles.SaveAll("/tmp/uploads")

		// Respond with success
		return c.Status(200).JSON(Msg{
			Msg: fmt.Sprintf("Uploaded %d files successfully to S3!", len(uploadedFiles)),
		})
	})

	q.Listen("0.0.0.0:8080")
}
