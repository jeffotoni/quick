// running various tests, such as fuzz, table driver and in the traditional way
// $ go test -run=^$ -fuzz=FuzzTestFormFile -fuzztime=5s
package quick

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// UploadedFileJSON represents only the serializable fields of UploadedFile
type UploadedFileJSON struct {
	Info struct {
		Filename    string `json:"Filename"`
		Size        int64  `json:"Size"`
		ContentType string `json:"ContentType"`
		Bytes       string `json:"Bytes"` // Base64 encoded
	} `json:"Info"`
}

// TestFormFile simulates a file upload request and validates the response.
// This function is named TestFormFile
// The result will TestFormFile(t *testing.T)
func TestFormFile(t *testing.T) {

	// start Quick
	q := New()

	q.Post("/upload", func(c *Ctx) error {
		uploadedFile, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(map[string]string{
				"error": err.Error(),
			})
		}

		return c.Status(200).JSONIN(uploadedFile)
	})

	ts := httptest.NewServer(q)
	defer ts.Close()

	// Call the HTTP abstraction function
	bodyBytes, _ := sendMultipartRequest(t, ts.URL, "testfile.txt", "Hello, Quick!")

	// Decode response
	var uploadedFile UploadedFileJSON
	err := json.Unmarshal(bodyBytes, &uploadedFile)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Decode Base64 to get original file bytes
	fileBytes, err := base64.StdEncoding.DecodeString(uploadedFile.Info.Bytes)
	if err != nil {
		t.Fatalf("failed to decode base64 file bytes: %v", err)
	}

	// Validate the uploaded file fields
	if uploadedFile.Info.Filename != "testfile.txt" {
		t.Errorf("expected filename %s, got %s", "testfile.txt", uploadedFile.Info.Filename)
	}

	if uploadedFile.Info.ContentType != "text/plain; charset=utf-8" {
		t.Errorf("expected content type %s, got %s", "text/plain; charset=utf-8", uploadedFile.Info.ContentType)
	}

	if uploadedFile.Info.Size != int64(13) {
		t.Errorf("expected size %d, got %d", 13, uploadedFile.Info.Size)
	}

	if string(fileBytes) != "Hello, Quick!" {
		t.Errorf("expected file content %q, got %q", "Hello, Quick!", string(fileBytes))
	}
}

// TestFormFileTableDriven simulates a file upload request and validates the response.
// This function is named TestFormFileTableDriven
// The result will TestFormFileTableDriven(t *testing.T)
func TestFormFileTableDriven(t *testing.T) {
	// Start Quick server
	q := New()
	q.Post("/upload", func(c *Ctx) error {
		uploadedFile, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(map[string]string{
				"error": err.Error(),
			})
		}
		return c.Status(200).JSONIN(uploadedFile)
	})

	ts := httptest.NewServer(q)
	defer ts.Close()

	// Define test cases
	tests := []struct {
		name          string
		fileName      string
		fileContent   string
		expectedSize  int64
		expectedType  string
		expectedError bool
	}{
		{
			name:         "Valid TXT file",
			fileName:     "testfile.txt",
			fileContent:  "Hello, Quick!",
			expectedSize: 13,
			expectedType: "text/plain; charset=utf-8",
		},
		{
			name:         "Valid TXT file",
			fileName:     "sample.txt",
			fileContent:  "Quick Framework Test",
			expectedSize: 22,
			expectedType: "text/plain; charset=utf-8",
		},
		{
			name:          "Empty file",
			fileName:      "empty.txt",
			fileContent:   "",
			expectedSize:  0,
			expectedType:  "",
			expectedError: true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.fileContent) == 0 {
				t.Skip("Skipping test: file content is empty")
			}

			// Call the HTTP abstraction function
			bodyBytes, _ := sendMultipartRequest(t, ts.URL, tt.fileName, tt.fileContent)

			// Debugging: Print raw response
			// fmt.Println("Raw Response Body:", string(bodyBytes))

			// Decode JSON response
			var uploadedFile UploadedFileJSON
			err := json.Unmarshal(bodyBytes, &uploadedFile)
			if err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			// Handle error cases
			if tt.expectedError {
				// Ensure the file size is 0 (empty file)
				if uploadedFile.Info.Size != 0 {
					t.Errorf("expected empty file, but got size %d", uploadedFile.Info.Size)
				}

				// Allow multiple valid `Content-Type` values for empty files
				allowedContentTypes := []string{"", "application/octet-stream", "text/plain; charset=utf-8"}
				validContentType := false
				for _, ct := range allowedContentTypes {
					if uploadedFile.Info.ContentType == ct {
						validContentType = true
						break
					}
				}

				// If `Content-Type` is not in the allowed list, fail the test
				if !validContentType {
					t.Errorf("expected empty content type, but got %s", uploadedFile.Info.ContentType)
				}

				return
			}

			// Decode Base64 file bytes
			fileBytes, err := base64.StdEncoding.DecodeString(uploadedFile.Info.Bytes)
			if err != nil {
				t.Fatalf("failed to decode base64 file bytes: %v", err)
			}

			// Debugging: Print the real content
			// fmt.Printf("Decoded File Content: [%q] (Length: %d)\n", string(fileBytes), len(fileBytes))

			expectedContent := tt.fileContent
			receivedContent := string(fileBytes)

			if receivedContent != expectedContent {
				t.Errorf("expected file content %q, but got %q", expectedContent, receivedContent)
			}

			// Validate fields
			if uploadedFile.Info.Filename != tt.fileName {
				t.Errorf("expected filename %s, got %s", tt.fileName, uploadedFile.Info.Filename)
			}

			if uploadedFile.Info.ContentType != tt.expectedType {
				t.Errorf("expected content type %s, got %s", tt.expectedType, uploadedFile.Info.ContentType)
			}

			if string(fileBytes) != tt.fileContent {
				t.Errorf("expected file content %q, got %q", tt.fileContent, string(fileBytes))
			}
		})
	}
}

// FuzzTestFormFile simulates a file upload request and validates the response.
// This function is named FuzzTestFormFile
// The result will FuzzTestFormFile(t *testing.T)
// $ go test -fuzz=FuzzTestFormFile -fuzztime=5s
func FuzzTestFormFile(f *testing.F) {
	// Seed initial test cases
	f.Add("testfile1.txt", "Hello, Quick!")
	f.Add("testfile2.txt", "Testing Fuzzing!")
	f.Add("emptyfile.txt", "")
	f.Add("longfile.txt", "This is a very long test string.")

	// Initialize Quick
	q := New()
	q.Post("/upload", func(c *Ctx) error {
		uploadedFile, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(map[string]string{
				"error": err.Error(),
			})
		}
		return c.Status(200).JSONIN(uploadedFile)
	})

	ts := httptest.NewServer(q)
	defer ts.Close()

	// Run Fuzzing
	f.Fuzz(func(t *testing.T, fileName string, fileContent string) {
		if len(fileContent) == 0 {
			t.Skip("Skipping test: file content is empty")
		}

		// Call the HTTP abstraction function
		bodyBytes, _ := sendMultipartRequest(t, ts.URL, fileName, fileContent)

		// Debugging: Print response
		// fmt.Printf("[%s] Raw Response Body: %q\n", fileName, string(bodyBytes))

		// Decode response
		var uploadedFile UploadedFileJSON
		err := json.Unmarshal(bodyBytes, &uploadedFile)
		if err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(uploadedFile.Info.Bytes) > 0 {
			// Decode Base64 file bytes
			fileBytes, err := base64.StdEncoding.DecodeString(uploadedFile.Info.Bytes)
			if err != nil {
				t.Fatalf("failed to decode base64 file bytes: %v", err)
			}

			// Validate content
			if string(fileBytes) != fileContent {
				t.Errorf("expected file content %q, but got %q", fileContent, string(uploadedFile.Info.Bytes))
			}
		}
	})
}

// sendMultipartRequest creates and sends a POST request with multipart/form-data
// using the provided test server URL, fileName and fileContent.
// It returns the response body as a byte slice or an error if any step fails.
// The result will sendMultipartRequest(t *testing.T, tsURL, fileName, fileContent string) ([]byte, error)
func sendMultipartRequest(t *testing.T, tsURL, fileName, fileContent string) ([]byte, error) {
	// Validate that fileContent is not empty
	if fileContent == "" {
		return nil, fmt.Errorf("file content cannot be empty")
	}

	// Create a new multipart/form-data request body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create the form file field
	formFile, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}

	// Write the file content to the form file field
	_, err = io.Copy(formFile, bytes.NewReader([]byte(fileContent)))
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %v", err)
	}

	// Close the writer to finalize the multipart body
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %v", err)
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", tsURL+"/upload", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set the correct Content-Type header
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read and return the response body
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return respBytes, nil
}

// test ensures that we can read the file multiple times without losing the original data.
// The result will TestQuick_UploadFileReset(t *testing.T)
func TestQuick_UploadFileReset(t *testing.T) {
	q := New()

	// Create an upload endpoint
	q.Post("/upload", func(c *Ctx) error {
		// Get the file from the request
		uploadedFile, err := c.FormFile("file")
		if err != nil {
			t.Fatalf("Error getting file: %v", err)
		}

		// Check if the file data is correct
		if uploadedFile.Info.Size == 0 {
			t.Errorf("Error: file size is zero")
		}
		if len(uploadedFile.Info.Bytes) == 0 {
			t.Errorf("Error: file was not stored correctly")
		}

		// Now let's test if the reset works
		file, err := uploadedFile.Multipart.Open()
		if err != nil {
			t.Fatalf("Error opening file again: %v", err)
		}

		// Create a buffer to read the data again
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, file); err != nil {
			t.Fatalf("Error copying file data after reset: %v", err)
		}

		// Check if the data read again matches the original
		if !bytes.Equal(buf.Bytes(), uploadedFile.Info.Bytes) {
			t.Errorf("Error: reset failed, data read is not equal to original")
		}

		return c.String("Upload successful")
	})

	// Simulate a file for upload
	fileContent := []byte("This is an upload test to check reset.")
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatalf("Error creating form file: %v", err)
	}
	part.Write(fileContent)
	writer.Close()

	// Create the upload request
	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a ResponseRecorder to capture the response
	rec := httptest.NewRecorder()

	// Run the request through Quick
	q.ServeHTTP(rec, req)

	// Check the response
	if rec.Code != http.StatusOK {
		t.Errorf("Error: unexpected response, code %d", rec.Code)
	}

	if rec.Body.String() != "Upload successful" {
		t.Errorf("Error: Unexpected response from server, received: %s", rec.Body.String())
	}
}

// TestFormFileLimit verifies that the FormFileLimit function correctly  processes file sizes.
func TestFormFileLimit(t *testing.T) {

	// We create a dummy context for the test
	c := &Ctx{}

	t.Run("Valid limit - 10MB", func(t *testing.T) {
		err := c.FormFileLimit("10MB")
		if err != nil {
			t.Errorf("Unexpected error for 10MB: %v", err)
		}
		if c.uploadFileSize != 10*1024*1024 {
			t.Errorf("Expected 10MB (%d bytes), but got %d bytes", 10*1024*1024, c.uploadFileSize)
		}
	})

	t.Run("Valid limit - 2GB", func(t *testing.T) {
		err := c.FormFileLimit("2GB")
		if err != nil {
			t.Errorf("Unexpected error for 2GB: %v", err)
		}
		if c.uploadFileSize != 2*1024*1024*1024 {
			t.Errorf("Expected 2GB (%d bytes), but got %d bytes", 2*1024*1024*1024, c.uploadFileSize)
		}
	})

	t.Run("Invalid format - Text without numbers", func(t *testing.T) {
		err := c.FormFileLimit("abc")
		if err == nil {
			t.Errorf("Expected error for invalid input, but no error occurred")
		}
	})

	t.Run("Invalid format - Unknown drive", func(t *testing.T) {
		err := c.FormFileLimit("10XY")
		if err == nil {
			t.Errorf("Expected error for unknown drive, but no error occurred")
		}
	})

	t.Run("Invalid format - Negative number", func(t *testing.T) {
		err := c.FormFileLimit("-5MB")
		if err == nil {
			t.Errorf("Expected error for negative number, but no error occurred")
		}
	})
}

func TestFormFile_Cover(t *testing.T) {
	t.Run("Error calling FormFiles", func(t *testing.T) {
		// Create a fake context that simulates an error in FormFiles
		c := &Ctx{
			Request: nil, // Simulates an internal error
		}

		_, err := c.FormFile("file")
		if err == nil {
			t.Errorf("Expected error calling FormFiles, but no error occurred")
		}
	})
	t.Run("No file uploaded", func(t *testing.T) {
		// We create a valid request without files
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.Close() // Closes the writer without adding files

		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		c := &Ctx{Request: req}

		_, err := c.FormFile("file")
		if err == nil || err.Error() != "no files found in the request" { // Updating the expected message
			t.Errorf("Expected error 'no files found in the request', but got: %v", err)
		}
	})

	t.Run("File uploaded successfully", func(t *testing.T) {
		// We create a valid request with a simulated file
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "testfile.txt")
		if err != nil {
			t.Fatalf("Error creating form file: %v", err)
		}
		part.Write([]byte("Hello, Quick!"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		c := &Ctx{Request: req}

		file, err := c.FormFile("file")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if file.Multipart.Filename != "testfile.txt" {
			t.Errorf("Expected 'testfile.txt', but got '%s'", file.Multipart.Filename)
		}
	})

	t.Run("Valid multipart form", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add a form field
		_ = writer.WriteField("name", "Quick")

		// Add a file to the form
		part, err := writer.CreateFormFile("file", "test.txt")
		if err != nil {
			t.Fatalf("Error creating file in form: %v", err)
		}
		part.Write([]byte("Hello, Quick!"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		c := &Ctx{Request: req, uploadFileSize: 10 * 1024 * 1024}

		form, err := c.MultipartForm()
		if err != nil {
			t.Fatalf("Unexpected error getting MultipartForm: %v", err)
		}

		if form == nil {
			t.Fatal("Expected MultipartForm, but got nil")
		}

		if form.Value["name"][0] != "Quick" {
			t.Errorf("Expected 'Quick', but got '%s'", form.Value["name"][0])
		}

		if len(form.File["file"]) == 0 {
			t.Errorf("Expected a file in the form, but none was found")
		}
	})

	t.Run("Error processing MultipartForm", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/upload", nil) // No body

		c := &Ctx{Request: req, uploadFileSize: 10 * 1024 * 1024}

		_, err := c.MultipartForm()
		if err == nil {
			t.Fatal("Expected an error processing MultipartForm, but none occurred")
		}
	})

	t.Run("Retrieve value from a form field", func(t *testing.T) {
		data := "username=jeffotoni"
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(data))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		c := &Ctx{Request: req}

		value := c.FormValue("username")
		if value != "jeffotoni" {
			t.Errorf("Expected 'jeffotoni', but got '%s'", value)
		}
	})

	t.Run("Returns empty string if field does not exist", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		c := &Ctx{Request: req}

		value := c.FormValue("password")
		if value != "" {
			t.Errorf("Expected empty value, but got '%s'", value)
		}
	})

	t.Run("Retrieves all form values", func(t *testing.T) {
		data := "field1=value1&field2=value2"
		req := httptest.NewRequest(http.MethodPost, "/form", bytes.NewBufferString(data))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		c := &Ctx{Request: req}

		values := c.FormValues()
		if len(values) != 2 {
			t.Errorf("Expected 2 values, but got %d", len(values))
		}

		if values["field1"][0] != "value1" {
			t.Errorf("Expected 'value1', but got '%s'", values["field1"][0])
		}

		if values["field2"][0] != "value2" {
			t.Errorf("Expected 'value2', but got '%s'", values["field2"][0])
		}
	})

	t.Run("Returns empty map if there is no data in the form", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/form", nil)
		c := &Ctx{Request: req}

		values := c.FormValues()
		if len(values) != 0 {
			t.Errorf("Expected empty map, but got %d values", len(values))
		}
	})
}

func TestUploadedFileMethods(t *testing.T) {
	mockFile := &UploadedFile{
		Info: FileInfo{
			Filename:    "testfile.txt",
			Size:        1024,
			ContentType: "text/plain",
			Bytes:       []byte("Hello, Quick!"),
		},
	}

	t.Run("FileName returns the correct name", func(t *testing.T) {
		if mockFile.FileName() != "testfile.txt" {
			t.Errorf("Expected 'testfile.txt', but got '%s'", mockFile.FileName())
		}
	})

	t.Run("Size returns the correct size", func(t *testing.T) {
		if mockFile.Size() != 1024 {
			t.Errorf("Expected 1024 bytes, but got %d", mockFile.Size())
		}
	})

	t.Run("ContentType returns the correct MIME type", func(t *testing.T) {
		if mockFile.ContentType() != "text/plain" {
			t.Errorf("Expected 'text/plain' but got '%s'", mockFile.ContentType())
		}
	})

	t.Run("Bytes returns the correct bytes", func(t *testing.T) {
		expectedBytes := []byte("Hello, Quick!")
		if !bytes.Equal(mockFile.Bytes(), expectedBytes) {
			t.Errorf("The bytes returned do not match what is expected")
		}
	})
}

func TestUploadedFileSave(t *testing.T) {
	t.Run("Save saves the file correctly", func(t *testing.T) {
		mockFile := &UploadedFile{
			Info: FileInfo{
				Filename: "testfile.txt",
				Bytes:    []byte("Hello, Quick!"),
			},
		}

		destination := t.TempDir() // Creating a temporary directory to save the file
		err := mockFile.Save(destination)
		if err != nil {
			t.Fatalf("Unexpected error saving file: %v", err)
		}

		// Check if the file was saved correctly
		savedFilePath := filepath.Join(destination, "testfile.txt")
		savedData, err := os.ReadFile(savedFilePath)
		if err != nil {
			t.Fatalf("Error reading saved file: %v", err)
		}

		if !bytes.Equal(savedData, mockFile.Bytes()) {
			t.Errorf("The saved file data does not match what is expected")
		}
	})

	t.Run("Save failed to save an empty file", func(t *testing.T) {
		mockFile := &UploadedFile{
			Info: FileInfo{
				Filename: "emptyfile.txt",
				Bytes:    []byte{}, // Empty file
			},
		}

		destination := t.TempDir()
		err := mockFile.Save(destination)
		if err == nil || err.Error() != "no file available to save" {
			t.Errorf("Expected error 'no file available to save', but got: %v", err)
		}
	})

	t.Run("Save failed when creating invalid target directory", func(t *testing.T) {
		mockFile := &UploadedFile{
			Info: FileInfo{
				Filename: "testfile.txt",
				Bytes:    []byte("Hello, Quick!"),
			},
		}

		invalidDestination := "/invalid/path"
		err := mockFile.Save(invalidDestination)
		if err == nil || err.Error() != "failed to create destination directory" {
			t.Errorf("Expected error 'failed to create destination directory', but got: %v", err)
		}
	})
}

func TestSaveAll(t *testing.T) {
	t.Run("SaveAll saves multiple files correctly", func(t *testing.T) {
		mockFiles := []*UploadedFile{
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

		destination := t.TempDir()
		err := SaveAll(mockFiles, destination)
		if err != nil {
			t.Fatalf("Unexpected error saving multiple files: %v", err)
		}

		// Check if the files were saved correctly
		for _, file := range mockFiles {
			savedFilePath := filepath.Join(destination, file.FileName())
			savedData, err := os.ReadFile(savedFilePath)
			if err != nil {
				t.Fatalf("Error reading saved file: %v", err)
			}

			if !bytes.Equal(savedData, file.Bytes()) {
				t.Errorf("The data in the saved file '%s' does not match what is expected", file.FileName())
			}
		}
	})

	t.Run("SaveAll fails if one of the files cannot be saved", func(t *testing.T) {
		mockFiles := []*UploadedFile{
			{
				Info: FileInfo{
					Filename: "file1.txt",
					Bytes:    []byte("File 1 content"),
				},
			},
			{
				Info: FileInfo{
					Filename: "emptyfile.txt",
					Bytes:    []byte{}, // Empty file that will cause an error
				},
			},
		}

		destination := t.TempDir()
		err := SaveAll(mockFiles, destination)
		if err == nil || err.Error() != "no file available to save" {
			t.Errorf("Expected error 'no file available to save', but got: %v", err)
		}
	})
}

func TestUploadedFileSaveAdditionalCoverage(t *testing.T) {
	t.Run("Save with custom filename", func(t *testing.T) {
		mockFile := &UploadedFile{
			Info: FileInfo{
				Filename: "default.txt",
				Bytes:    []byte("Hello, Quick!"),
			},
		}

		destination := t.TempDir()
		customFileName := "custom_name.txt"
		err := mockFile.Save(destination, customFileName)
		if err != nil {
			t.Fatalf("Unexpected error saving file with custom name: %v", err)
		}

		// Check if file was saved with custom name
		savedFilePath := filepath.Join(destination, customFileName)
		if _, err := os.Stat(savedFilePath); os.IsNotExist(err) {
			t.Errorf("Expected file '%s' to be created, but it does not exist", savedFilePath)
		}
	})

	t.Run("Save fails when os.Create fails", func(t *testing.T) {
		mockFile := &UploadedFile{
			Info: FileInfo{
				Filename: "testfile.txt",
				Bytes:    []byte("Hello, Quick!"),
			},
		}

		invalidDestination := "/invalid/path"
		err := mockFile.Save(invalidDestination)

		if err == nil {
			t.Fatalf("Expected an error, but got none")
		}

		expectedErrors := []string{
			"failed to create file on disk",
			"failed to create destination directory",
		}

		isValidError := false
		for _, expected := range expectedErrors {
			if err.Error() == expected {
				isValidError = true
				break
			}
		}

		if !isValidError {
			t.Errorf("Expected error to be one of %v, but got: %v", expectedErrors, err)
		}
	})

	t.Run("Save fails when writing file content fails", func(t *testing.T) {
		mockFile := &UploadedFile{
			Info: FileInfo{
				Filename: "testfile.txt",
				Bytes:    []byte("Hello, Quick!"),
			},
		}

		// Simulate an invalid file path by making a directory with the same name
		destination := t.TempDir()
		invalidFilePath := filepath.Join(destination, "testfile.txt")
		err := os.Mkdir(invalidFilePath, os.ModePerm) // Create a directory with file name
		if err != nil {
			t.Fatalf("Failed to create directory for test: %v", err)
		}

		err = mockFile.Save(destination)
		if err == nil || err.Error() != "failed to create file on disk" {
			t.Errorf("Expected error 'failed to create file on disk', but got: %v", err)
		}
	})
}
func TestFormFileLimitErrors(t *testing.T) {
	c := &Ctx{}

	t.Run("Invalid size number", func(t *testing.T) {
		err := c.FormFileLimit("MB") // Missing number before unit
		if err == nil || err.Error() != "invalid size format" {
			t.Errorf("Expected 'invalid size format' error, but got: %v", err)
		}
	})

	t.Run("Unknown size unit", func(t *testing.T) {
		err := c.FormFileLimit("10XY") // Unknown unit
		if err == nil || err.Error() != "invalid size format" {
			t.Errorf("Expected 'invalid size format' error, but got: %v", err)
		}
	})
}
