package quick

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Ctx struct {
	Response       http.ResponseWriter
	Request        *http.Request
	resStatus      int
	MoreRequests   int
	bodyByte       []byte
	JsonStr        string
	Headers        map[string][]string
	Params         map[string]string
	Query          map[string]string
	uploadFileSize int64 // Upload limit in bytes
}

func (c *Ctx) SetStatus(status int) {
	c.resStatus = status
}

// UploadedFile holds details of an uploaded file.
type UploadedFile struct {
	File      multipart.File
	Multipart *multipart.FileHeader
	Info      FileInfo
}

// FileInfo contains metadata of the uploaded file.
type FileInfo struct {
	Filename    string
	Size        int64
	ContentType string
	Bytes       []byte
}

// GetReqHeadersAll returns all the request headers
// The result will GetReqHeadersAll() map[string][]string
func (c *Ctx) GetReqHeadersAll() map[string][]string {
	return c.Headers
}

// GetHeadersAll returns all HTTP response headers stored in the context
// The result will GetHeadersAll() map[string][]string
func (c *Ctx) GetHeadersAll() map[string][]string {
	return c.Headers
}

// Http serveFile send specific file
// The result will File(filePath string)
func (c *Ctx) File(filePath string) error {
	filePath = strings.TrimSuffix(filePath, "/*")

	if stat, err := os.Stat(filePath); err == nil && stat.IsDir() {
		filePath = filepath.Join(filePath, "index.html")
	}
	http.ServeFile(c.Response, c.Request, filePath)
	return nil
}

// Bind analyzes and links the request body to a Go structure
// The result will Bind(v interface{}) (err error)
func (c *Ctx) Bind(v interface{}) (err error) {
	return extractParamsBind(c, v)
}

// BodyParser efficiently unmarshals the request body into the provided struct (v) based on the Content-Type header.
//
// Supported content-types:
// - application/json
// - application/xml, text/xml
//
// Parameters:
//   - v: The target structure to decode the request body into.
//
// Returns:
//   - error: An error if decoding fails or if the content-type is unsupported.
func (c *Ctx) BodyParser(v interface{}) error {
	contentType := strings.ToLower(c.Request.Header.Get("Content-Type"))

	switch {
	case strings.HasPrefix(contentType, ContentTypeAppJSON):
		return json.Unmarshal(c.bodyByte, v)

	case strings.Contains(contentType, ContentTypeAppXML),
		strings.Contains(contentType, ContentTypeTextXML):
		return xml.Unmarshal(c.bodyByte, v)

	default:
		return fmt.Errorf("unsupported content-type: %s", contentType)
	}
}

// Param returns the value of the URL parameter corresponding to the given key
// The result will Param(key string) string
func (c *Ctx) Param(key string) string {
	val, ok := c.Params[key]
	if ok {
		return val
	}
	return ""
}

// Body returns the request body as a byte slice ([]byte)
// The result will Body() []byte
func (c *Ctx) Body() []byte {
	return c.bodyByte
}

// BodyString returns the request body as a string
// The result will BodyString() string
func (c *Ctx) BodyString() string {
	return string(c.bodyByte)
}

// JSON encodes the provided interface (v) as JSON, sets the Content-Type header,
// and writes the response efficiently using buffer pooling.
//
// Parameters:
//   - v: The data structure to encode as JSON.
//
// Returns:
//   - error: An error if JSON encoding fails or if writing the response fails.
func (c *Ctx) JSON(v interface{}) error {
	buf := acquireJSONBuffer()
	defer releaseJSONBuffer(buf)

	if err := json.NewEncoder(buf).Encode(v); err != nil {
		return err
	}

	if buf.Len() > 0 && buf.Bytes()[buf.Len()-1] == '\n' {
		buf.Truncate(buf.Len() - 1)
	}

	c.writeResponse(buf.Bytes())
	return nil
}

// JSONIN encodes the given interface as JSON with indentation and writes it to the HTTP response.
// Allows optional parameters to define the indentation format.
//
// ATTENTION
// use only for debugging, very slow
//
// Parameters:
//   - v: The data structure to encode as JSON.
//   - params (optional): Defines the indentation settings.
//   - If params[0] is provided, it will be used as the prefix.
//   - If params[1] is provided, it will be used as the indentation string.
//
// Returns:
//   - error: An error if JSON encoding fails or if writing to the ResponseWriter fails.
func (c *Ctx) JSONIN(v interface{}, params ...string) error {
	// Default indentation settings
	prefix := ""
	indent := "  " // Default to 2 spaces

	// Override if parameters are provided
	if len(params) > 0 {
		prefix = params[0]
	}
	if len(params) > 1 {
		indent = params[1]
	}

	// Encode with the provided indentation settings
	b, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		return err
	}

	c.writeResponse(b)
	return nil
}

// XML serializes the given value to XML and writes it to the HTTP response.
// It avoids unnecessary memory allocations by using buffer pooling and ensures that no extra newline is appended.
//
// Parameters:
//   - v: The data structure to encode as XML.
//
// Returns:
//   - error: An error if XML encoding fails or if writing to the ResponseWriter fails.
func (c *Ctx) XML(v interface{}) error {
	buf := acquireXMLBuffer()
	defer releaseXMLBuffer(buf)

	if err := xml.NewEncoder(buf).Encode(v); err != nil {
		return err
	}

	if buf.Len() > 0 && buf.Bytes()[buf.Len()-1] == '\n' {
		buf.Truncate(buf.Len() - 1)
	}

	c.writeResponse(buf.Bytes())
	return nil
}

// writeResponse writes the provided byte content to the ResponseWriter.
//
// If a custom status code (resStatus) has been set, it writes the header before the body.
//
// Parameters:
//   - b: The byte slice to be written in the HTTP response.
//
// Returns:
//   - error: An error if writing to the ResponseWriter fails.
func (c *Ctx) writeResponse(b []byte) error {
	if c.Response == nil {
		return errors.New("nil response writer")
	}

	if c.resStatus == 0 {
		c.resStatus = http.StatusOK
	}

	c.Response.WriteHeader(c.resStatus)

	_, err := c.Response.Write(b)
	if flusher, ok := c.Response.(http.Flusher); ok {
		flusher.Flush()
	}

	return err
}

// Byte writes an array of bytes to the HTTP response, using writeResponse()
// The result will Byte(b []byte) (err error)
func (c *Ctx) Byte(b []byte) (err error) {
	return c.writeResponse(b)
}

// Send writes a byte array to the HTTP response, using writeResponse()
// The result will Send(b []byte) (err error)
func (c *Ctx) Send(b []byte) (err error) {
	return c.writeResponse(b)
}

// SendString writes a string in the HTTP response, converting it to an array of bytes and using writeResponse()
// The result will SendString(s string) error
func (c *Ctx) SendString(s string) error {
	return c.writeResponse([]byte(s))
}

// String escreve uma string na resposta HTTP, convertendo-a para um array de bytes e utilizando writeResponse()
// The result will String(s string) error
func (c *Ctx) String(s string) error {
	return c.writeResponse([]byte(s))
}

// SendFile writes a file in the HTTP response as an array of bytes
// The result will SendFile(file []byte) error
func (c *Ctx) SendFile(file []byte) error {
	_, err := c.Response.Write(file)
	return err
}

// Set defines an HTTP header in the response
// The result will Set(key, value string)
func (c *Ctx) Set(key, value string) {
	c.Response.Header().Set(key, value)
}

// Append adds a value to the HTTP header specified in the response
// The result will Append(key, value string)
func (c *Ctx) Append(key, value string) {
	c.Response.Header().Add(key, value)
}

// Accepts defines the HTTP header "Accept" in the response
// The result will Accepts(acceptType string) *Ctx
func (c *Ctx) Accepts(acceptType string) *Ctx {
	c.Response.Header().Set("Accept", acceptType)
	return c
}

// Status defines the HTTP status code of the response
// The result will Status(status int) *Ctx
func (c *Ctx) Status(status int) *Ctx {
	c.resStatus = status
	return c
}

//MultipartForm

// FormFileLimit sets the maximum allowed upload size.
func (c *Ctx) FormFileLimit(limit string) error {
	size, err := parseSize(limit)
	if err != nil {
		return err
	}
	c.uploadFileSize = size
	return nil
}

// FormFile processes an uploaded file and returns its details.
// The result will FormFile(fieldName string) (*UploadedFile, error)
func (c *Ctx) FormFile(fieldName string) (*UploadedFile, error) {
	files, err := c.FormFiles(fieldName)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, errors.New("no file uploaded")
	}

	return files[0], nil // Return the first file if multiple are uploaded
}

// fileWrapper, which wraps a bytes.Reader and adds the Close() method,
// allowing it to be treated as an io.ReadCloser.
// We ensure that the file can be read multiple times without losing data.
// fileWrapper supports multipart.File.
type fileWrapper struct {
	*bytes.Reader
}

// There is nothing to close as we are reading from memory
func (fw *fileWrapper) Close() error {
	return nil
}

// FormFiles processes an uploaded file and returns its details.
// The result will FormFiles(fieldName string) (*UploadedFile, error)
func (c *Ctx) FormFiles(fieldName string) ([]*UploadedFile, error) {
	if c.uploadFileSize == 0 {
		c.uploadFileSize = 1 << 20 // set default 1MB
	}

	// check request
	if c.Request == nil {
		return nil, errors.New("HTTP request is nil")
	}

	// check body
	if c.Request.Body == nil {
		return nil, errors.New("request body is nil")
	}

	// check if `Content-Type` this ok
	contentType := c.Request.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		return nil, errors.New("invalid content type, expected multipart/form-data")
	}

	// Parse multipart form with the defined limit
	if err := c.Request.ParseMultipartForm(c.uploadFileSize); err != nil {
		return nil, errors.New("failed to parse multipart form: " + err.Error())
	}

	// Debugging: Check if files exist
	if c.Request.MultipartForm == nil || c.Request.MultipartForm.File[fieldName] == nil {
		return nil, errors.New("no files found in the request")
	}

	// Retrieve all files for the given field name
	files := c.Request.MultipartForm.File[fieldName]
	if len(files) == 0 {
		return nil, errors.New("no files found for field: " + fieldName)
	}

	var uploadedFiles []*UploadedFile

	for _, handler := range files {
		// Open file
		file, err := handler.Open()
		if err != nil {
			return nil, errors.New("failed to open file: " + err.Error())
		}
		defer file.Close()

		// Read file content into memory
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, file); err != nil {
			return nil, errors.New("failed to read file into buffer")
		}

		// reset  multipart.File
		// Create a reusable copy of the file
		// that implements multipart.File correctly
		fileCopy := &fileWrapper{bytes.NewReader(buf.Bytes())}

		// Detect content type
		contentType := http.DetectContentType(buf.Bytes())

		// Append file details
		uploadedFiles = append(uploadedFiles, &UploadedFile{
			File:      fileCopy,
			Multipart: handler,
			Info: FileInfo{
				Filename:    handler.Filename,
				Size:        handler.Size,
				ContentType: contentType,
				Bytes:       buf.Bytes(),
			},
		})
	}

	return uploadedFiles, nil
}

// MultipartForm allows access to the raw multipart form data (for advanced users)
// The result will MultipartForm() (*multipart.Form, error)
func (c *Ctx) MultipartForm() (*multipart.Form, error) {
	if err := c.Request.ParseMultipartForm(c.uploadFileSize); err != nil {
		return nil, err
	}
	return c.Request.MultipartForm, nil
}

// FormValue retrieves a form value by key.
// It automatically calls ParseForm() before accessing the value.
// The result will FormValue(key string) string
func (c *Ctx) FormValue(key string) string {
	// Checks if the Content-Type is multipart
	if c.Request.Header.Get("Content-Type") == "multipart/form-data" {
		_ = c.Request.ParseMultipartForm(c.uploadFileSize) // Force correct processing
	} else {
		_ = c.Request.ParseForm() // For application/x-www-form-urlencoded
	}
	return c.Request.FormValue(key)
}

// FormValues returns all form values as a map.
// It automatically calls ParseForm() before accessing the values.
// The result will FormValues() map[string][]string
func (c *Ctx) FormValues() map[string][]string {
	// Checks if the Content-Type is multipart
	if c.Request.Header.Get("Content-Type") == "multipart/form-data" {
		_ = c.Request.ParseMultipartForm(c.uploadFileSize) // Required to process multipart
	} else {
		_ = c.Request.ParseForm() // Processes application/x-www-form-urlencoded
	}
	return c.Request.Form
}
