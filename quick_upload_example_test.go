// Package quick provides a high-performance HTTP framework with built-in utilities
// for handling file uploads, multipart form data, and request processing.
//
// This file contains Example functions for the GoDoc documentation, demonstrating
// how to use various methods in the Quick framework, particularly related to
// handling uploaded files, saving them, setting file size limits, and parsing multipart forms.
package quick

import (
	"fmt"
)

// form multipart/formdata
// Msg is an auxiliary structure for the examples
type Msg struct {
	Msg   string `json:"msg"`
	Error string `json:"error"`
}

// This function is named ExampleUploadedFile_FileName()
// it with the Examples type.
func ExampleUploadedFile_FileName() {
	// Start a new Quick instance
	q := New()

	// Define a POST route for file upload
	q.Post("/upload", func(c *Ctx) error {
		// Set file upload size limit
		c.FormFileLimit("10MB")

		// Retrieve the uploaded file
		uploadedFile, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(Msg{
				Msg:   "Upload error",
				Error: err.Error(),
			})
		}

		// Print only the file name
		fmt.Println(uploadedFile.FileName())

		// Return JSON response with the file name
		return c.Status(200).JSON(map[string]string{
			"name": uploadedFile.FileName(),
		})
	})

	// Simulate an UploadedFile manually (standalone example)
	uploadedFile := &UploadedFile{
		Info: FileInfo{
			Filename: "quick.txt",
		},
	}

	// Print only the file name
	fmt.Println(uploadedFile.FileName())

	// Output: quick.txt
}

// This function is named ExampleUploadedFile_Size()
// it with the Examples type.
func ExampleUploadedFile_Size() {

	// Start a new Quick instance
	q := New()

	// Define a POST route for file upload
	q.Post("/upload", func(c *Ctx) error {
		// Set file upload size limit
		c.FormFileLimit("10MB")

		// Retrieve the uploaded file
		uploadedFile, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(Msg{
				Msg:   "Upload error",
				Error: err.Error(),
			})
		}

		// Print only the file name
		fmt.Println(uploadedFile.FileName())

		// Return JSON response with the file name
		return c.Status(200).JSON(map[string]string{
			"name": uploadedFile.FileName(),
		})
	})

	// Creating an UploadedFile object simulating an uploaded file of 2048 bytes
	uploadedFile := &UploadedFile{
		Info: FileInfo{
			Size: 2048,
		},
	}

	// Retrieving the file size
	fmt.Println(uploadedFile.Size())

	// Output: 2048
}

// This function is named ExampleUploadedFile_ContentType()
// it with the Examples type.
func ExampleUploadedFile_ContentType() {
	// Start a new Quick instance
	q := New()

	// Define a POST route for file upload
	q.Post("/upload", func(c *Ctx) error {
		// Set file upload size limit
		c.FormFileLimit("10MB")

		// Retrieve the uploaded file
		uploadedFile, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(Msg{
				Msg:   "Upload error",
				Error: err.Error(),
			})
		}

		// Print only the file name
		fmt.Println(uploadedFile.FileName())

		// Return JSON response with the file name
		return c.Status(200).JSON(map[string]string{
			"name": uploadedFile.FileName(),
		})
	})

	// Creating an UploadedFile object simulating an uploaded PNG file
	uploadedFile := &UploadedFile{
		Info: FileInfo{
			ContentType: "image/png",
		},
	}

	// Retrieving the content type
	fmt.Println(uploadedFile.ContentType())

	// Output: image/png
}

// This function is named ExampleUploadedFile_Bytes()
// it with the Examples type.
func ExampleUploadedFile_Bytes() {
	// Start a new Quick instance
	q := New()

	// Define a POST route for file upload
	q.Post("/upload", func(c *Ctx) error {
		// Set file upload size limit
		c.FormFileLimit("10MB")

		// Retrieve the uploaded file
		uploadedFile, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(Msg{
				Msg:   "Upload error",
				Error: err.Error(),
			})
		}

		// Print only the file name
		fmt.Println(uploadedFile.FileName())

		// Return JSON response with the file name
		return c.Status(200).JSON(map[string]string{
			"name": uploadedFile.FileName(),
		})
	})

	// Creating an UploadedFile object with content as bytes
	uploadedFile := &UploadedFile{
		Info: FileInfo{
			Bytes: []byte("Hello, Quick!"),
		},
	}

	// Converting bytes to a string to display the content
	fmt.Println(string(uploadedFile.Bytes()))

	// Output: Hello, Quick!
}

// This function is named ExampleUploadedFile_Save()
// it with the Examples type.
func ExampleUploadedFile_Save() {
	// Start a new Quick instance
	q := New()

	// Define a POST route for file upload
	q.Post("/upload", func(c *Ctx) error {
		// Set file upload size limit
		c.FormFileLimit("10MB")

		// Retrieve the uploaded file
		uploadedFile, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(Msg{
				Msg:   "Upload error",
				Error: err.Error(),
			})
		}

		// Print only the file name
		fmt.Println(uploadedFile.FileName())

		// Return JSON response with the file name
		return c.Status(200).JSON(map[string]string{
			"name": uploadedFile.FileName(),
		})
	})

	// Creating an UploadedFile object simulating an uploaded file
	uploadedFile := &UploadedFile{
		Info: FileInfo{
			Filename: "quick.txt",
			Bytes:    []byte("Hello, Quick!"),
		},
	}

	// Saving the file to the "uploads" directory
	err := uploadedFile.Save("uploads")

	// Checking for errors
	if err != nil {
		fmt.Println("Error saving file:", err)
	} else {
		fmt.Println("File saved successfully!")
	}

	// Output: File saved successfully!
}

// This function is named ExampleSaveAll()
// it with the Examples type.
func ExampleSaveAll() {
	q := New()

	// Define a POST route for uploading multiple files
	q.Post("/upload", func(c *Ctx) error {
		// Set file upload size limit
		c.FormFileLimit("10MB")

		// Retrieve multiple uploaded files
		files, err := c.FormFiles("file")
		if err != nil {
			return c.Status(400).JSON(Msg{
				Msg:   "Upload error",
				Error: err.Error(),
			})
		}

		// Save all files to a directory (e.g., "uploads")
		err = SaveAll(files, "uploads")
		if err != nil {
			return c.Status(500).JSON(Msg{
				Msg:   "Failed to save files",
				Error: err.Error(),
			})
		}

		// Return JSON response indicating success
		return c.Status(200).JSON(map[string]string{
			"message": "All files saved successfully",
		})
	})

	// Create a list of UploadedFile instances to simulate multiple uploads
	files := []*UploadedFile{
		{
			Info: FileInfo{
				Filename: "file1.txt",
				Bytes:    []byte("File 1 content"),
			},
		},
		{
			Info: FileInfo{
				Filename: "file2.txt",
				Bytes:    []byte("File 2 content"),
			},
		},
	}

	// Define the target directory for saving the files
	targetDir := "test_uploads"

	// Save all files using SaveAll function
	err := SaveAll(files, targetDir)

	// Handle result
	if err != nil {
		fmt.Println("Error saving files:", err)
	} else {
		fmt.Println("All files saved successfully!")
	}

	// Output: All files saved successfully!
}
