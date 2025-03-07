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

// This function is named ExampleUploadedFile_FileName()
// it with the Examples type.
func ExampleUploadedFile_FileName() {
	// Creating an UploadedFile object simulating an uploaded file
	uploadedFile := &UploadedFile{
		Info: FileInfo{
			Filename: "quick.txt",
		},
	}

	// Retrieving the filename
	fmt.Println(uploadedFile.FileName())

	// Out put: quick.txt
}

// This function is named ExampleUploadedFile_Size()
// it with the Examples type.
func ExampleUploadedFile_Size() {
	// Creating an UploadedFile object simulating an uploaded file of 2048 bytes
	uploadedFile := &UploadedFile{
		Info: FileInfo{
			Size: 2048,
		},
	}

	// Retrieving the file size
	fmt.Println(uploadedFile.Size())

	// Out put: 2048
}

// This function is named ExampleUploadedFile_ContentType()
// it with the Examples type.
func ExampleUploadedFile_ContentType() {
	// Creating an UploadedFile object simulating an uploaded PNG file
	uploadedFile := &UploadedFile{
		Info: FileInfo{
			ContentType: "image/png",
		},
	}

	// Retrieving the content type
	fmt.Println(uploadedFile.ContentType())

	// Out put: image/png
}

// This function is named ExampleUploadedFile_Bytes()
// it with the Examples type.
func ExampleUploadedFile_Bytes() {
	// Creating an UploadedFile object with content as bytes
	uploadedFile := &UploadedFile{
		Info: FileInfo{
			Bytes: []byte("Hello, Quick!"),
		},
	}

	// Converting bytes to a string to display the content
	fmt.Println(string(uploadedFile.Bytes()))

	// Out put: Hello, Quick!
}

// This function is named ExampleUploadedFile_Save()
// it with the Examples type.
func ExampleUploadedFile_Save() {
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

	// Out put: File saved successfully!
}

// This function is named ExampleSaveAll()
// it with the Examples type.
func ExampleSaveAll() {
	// Creating a list of files to simulate multiple uploads
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

	// Creating a temporary directory for testing
	tempDir := "test_uploads"

	// Saving all files to the temporary directory
	err := SaveAll(files, tempDir)

	// Handling errors
	if err != nil {
		fmt.Println("Error saving files:", err)
	} else {
		fmt.Println("All files saved successfully!")
	}

	// Out put:All files saved successfully!
}
