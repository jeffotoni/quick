// Package quick provides a high-performance HTTP framework for building web applications in Go.
//
// Quick is designed to be lightweight and efficient, offering a simplified API for handling
// HTTP requests, file uploads, middleware, and routing.
//
// Features:
//   - Route management with support for grouped routes.
//   - Middleware support for request processing.
//   - File handling capabilities, including uploads and size validation.
//   - High-performance request handling using Go’s standard `net/http` package.
package quick

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jeffotoni/quick/internal/concat"
)

// FileName returns the name of the uploaded file.
//
// This function retrieves the original filename as provided during the upload.
//
// Returns:
//   - string: The name of the uploaded file.
func (uf *UploadedFile) FileName() string {
	return uf.Info.Filename
}

// Size returns the size of the uploaded file in bytes.
//
// This function retrieves the file size based on the uploaded file metadata.
//
// Returns:
//   - int64: The size of the file in bytes.
func (uf *UploadedFile) Size() int64 {
	return uf.Info.Size
}

// ContentType returns the MIME type of the uploaded file.
//
// This function retrieves the detected content type of the uploaded file.
//
// Returns:
//   - string: The MIME type of the file.
func (uf *UploadedFile) ContentType() string {
	return uf.Info.ContentType
}

// Bytes returns the raw content of the uploaded file as a byte slice.
//
// This function allows access to the file data in memory before saving it.
//
// Returns:
//   - []byte: The raw bytes of the uploaded file.
func (uf *UploadedFile) Bytes() []byte {
	return uf.Info.Bytes
}

// Save stores the uploaded file in the specified directory.
//
// If a filename is provided, the file will be saved with that name.
// Otherwise, the original filename is used.
//
// Parameters:
//   - destination (string): The directory where the file should be saved.
//   - nameFile (optional []string): Custom filename (optional).
//
// Returns:
//   - error: An error if the save operation fails.
//
// Example Usage:
//
//	err := uploadedFile.Save("./uploads", "newfile.png")
//	if err != nil {
//	    log.Fatal("Failed to save file:", err)
//	}
func (uf *UploadedFile) Save(destination string, nameFile ...string) error {
	var fullPath string

	// Determine the full file path based on the given filename or the original name
	if len(nameFile) > 0 {
		fullPath = filepath.Join(destination, nameFile[0])
	} else {
		if len(uf.Info.Bytes) == 0 {
			return errors.New("no file available to save")
		}
		fullPath = concat.String(destination, "/", uf.Info.Filename)
	}

	// Ensure the destination directory exists
	if err := os.MkdirAll(destination, os.ModePerm); err != nil {
		return errors.New("failed to create destination directory")
	}

	// Create the file on disk
	dst, err := os.Create(fullPath)
	if err != nil {
		return errors.New("failed to create file on disk")
	}
	defer dst.Close()

	// Write the file content from memory
	_, err = dst.Write(uf.Info.Bytes)
	if err != nil {
		return errors.New("failed to save file")
	}

	return nil
}

// SaveAll saves all uploaded files to the specified directory.
//
// This function iterates through a slice of uploaded files and saves them
// using the `Save()` method.
//
// Parameters:
//   - files ([]*UploadedFile): Slice of uploaded files.
//   - destination (string): The directory where the files should be saved.
//
// Returns:
//   - error: An error if any file fails to save.
//
// Example Usage:
//
//	files := []*UploadedFile{file1, file2}
//	err := SaveAll(files, "./uploads")
//	if err != nil {
//	    log.Fatal("Failed to save all files:", err)
//	}
func SaveAll(files []*UploadedFile, destination string) error {
	for _, file := range files {
		if err := file.Save(destination); err != nil {
			return err
		}
	}
	return nil
}

// parseSize converts a human-readable size string (e.g., "10MB") into bytes.
//
// This function parses strings like "10MB", "500KB", and converts them into
// their corresponding byte values.
//
// Parameters:
//   - sizeStr (string): A string representing the file size (e.g., "10MB").
//
// Returns:
//   - int64: The size in bytes.
//   - error: An error if the string format is invalid.
//
// Example Usage:
//
//	bytes, err := parseSize("10MB")
//	if err != nil {
//	    log.Fatal("Invalid size format:", err)
//	}
//	fmt.Println("Size in bytes:", bytes) // Output: 10485760 (10MB)
func parseSize(sizeStr string) (int64, error) {
	// Normalize string (trim spaces and convert to lowercase)
	sizeStr = strings.TrimSpace(strings.ToLower(sizeStr))

	// Regular expression to validate format (e.g., "10mb", "200kb", "2gb")
	re := regexp.MustCompile(`^(\d+)(b|kb|mb|gb|tb)$`)
	matches := re.FindStringSubmatch(sizeStr)

	if len(matches) != 3 {
		return 0, errors.New("invalid size format")
	}

	// Convert the numeric part to an integer
	value, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, errors.New("invalid size number")
	}

	// Define multipliers for different size units
	unitMultipliers := map[string]int64{
		"b":  1,
		"kb": 1024,
		"mb": 1024 * 1024,
		"gb": 1024 * 1024 * 1024,
		"tb": 1024 * 1024 * 1024 * 1024,
	}

	// Multiply the value by the corresponding unit multiplier
	multiplier, exists := unitMultipliers[matches[2]]
	if !exists {
		return 0, errors.New("unknown size unit")
	}

	return value * multiplier, nil
}
