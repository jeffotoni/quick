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

// form multipart/formdata
///

// FileName returns the uploaded file's name.
// The result will FileName() string
func (uf *UploadedFile) FileName() string {
	return uf.Info.Filename
}

// Size returns the size of the uploaded file in bytes.
// The result will Size() int64
func (uf *UploadedFile) Size() int64 {
	return uf.Info.Size
}

// ContentType returns the MIME type of the uploaded file.
// The result will ContentType() string
func (uf *UploadedFile) ContentType() string {
	return uf.Info.ContentType
}

// Bytes returns the raw bytes of the uploaded file.
// The result will Bytes() []byte
func (uf *UploadedFile) Bytes() []byte {
	return uf.Info.Bytes
}

// Save stores the uploaded file in the specified directory.
// The result will Save(destination string) error
func (uf *UploadedFile) Save(destination string, nameFile ...string) error {
	var fullPath string

	if len(nameFile) > 0 {
		// Join the files into a single string separated by "/"
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
// The result will SaveAll(files []*UploadedFile, destination string) error {
func SaveAll(files []*UploadedFile, destination string) error {
	for _, file := range files {
		if err := file.Save(destination); err != nil {
			return err
		}
	}
	return nil
}

// parseSize converts a human-readable size string (e.g., "10MB") to bytes.
// The result will parseSize(sizeStr string) (int64, error)
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
