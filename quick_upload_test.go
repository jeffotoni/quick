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
