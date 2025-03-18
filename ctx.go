// Ctx represents the context of an HTTP request and response.
//
// It provides access to the request, response, headers, query parameters,
// body, and other necessary attributes for handling HTTP requests.
//
// Fields:
//   - Response: The HTTP response writer.
//   - Request: The HTTP request object.
//   - resStatus: The HTTP response status code.
//   - MoreRequests: Counter for additional requests in a batch processing scenario.
//   - bodyByte: The raw body content as a byte slice.
//   - JsonStr: The raw body content as a string.
//   - Headers: A map containing all request headers.
//   - Params: A map containing URL parameters (e.g., /users/:id → id).
//   - Query: A map containing query parameters (e.g., ?name=John).
//   - uploadFileSize: The maximum allowed upload file size in bytes.
//   - App: A reference to the Quick application instance.
package quick

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Ctx represents the context of an HTTP request and response.
type Ctx struct {
	Response       http.ResponseWriter // HTTP response writer to send responses
	Request        *http.Request       // Incoming HTTP request object
	resStatus      int                 // HTTP response status code
	MoreRequests   int                 // Counter for batch processing requests
	bodyByte       []byte              // Raw request body as byte slice
	JsonStr        string              // Request body as a string (for JSON handling)
	Headers        map[string][]string // Map of request headers
	Params         map[string]string   // Map of URL parameters
	Query          map[string]string   // Map of query parameters
	uploadFileSize int64               // Maximum allowed file upload size (bytes)
	App            *Quick              // Reference to the Quick application instance
}

// SetStatus sets the HTTP response status code.
//
// Parameters:
//   - status: The HTTP status code to be set.
func (c *Ctx) SetStatus(status int) {
	c.resStatus = status
}

// UploadedFile holds details of an uploaded file.
//
// Fields:
//   - File: The uploaded file as a multipart.File.
//   - Multipart: The file header containing metadata about the uploaded file.
//   - Info: Additional information about the file, including filename, size, and content type.
type UploadedFile struct {
	File      multipart.File
	Multipart *multipart.FileHeader
	Info      FileInfo
}

// FileInfo contains metadata about an uploaded file.
//
// Fields:
//   - Filename: The original name of the uploaded file.
//   - Size: The file size in bytes.
//   - ContentType: The MIME type of the file (e.g., "image/png").
//   - Bytes: The raw file content as a byte slice.
type FileInfo struct {
	Filename    string
	Size        int64
	ContentType string
	Bytes       []byte
}

// GetHeader retrieves a specific header value from the request.
//
// Parameters:
//   - key: The name of the header to retrieve.
//
// Returns:
//   - string: The value of the specified header, or an empty string if not found.
func (c *Ctx) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

// GetHeaders returns all request headers.
//
// Returns:
//   - http.Header: A map containing all request headers.
func (c *Ctx) GetHeaders() http.Header {
	return c.Request.Header
}

// RemoteIP retrieves the client's IP address from the request.
//
// Returns:
//   - string: The client's IP address. If extraction fails, it returns the original RemoteAddr.
func (c *Ctx) RemoteIP() string {
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}

// Method returns the HTTP method of the request.
//
// Returns:
//   - string: The HTTP method (e.g., "GET", "POST").
func (c *Ctx) Method() string {
	return c.Request.Method
}

// Path returns the URL path of the request.
//
// Returns:
//   - string: The path component of the request URL.
func (c *Ctx) Path() string {
	return c.Request.URL.Path
}

// QueryParam retrieves a query parameter value from the URL.
//
// Parameters:
//   - key: The name of the query parameter to retrieve.
//
// Returns:
//   - string: The value of the specified query parameter, or an empty string if not found.
func (c *Ctx) QueryParam(key string) string {
	return c.Request.URL.Query().Get(key)
}

// QueryParam retrieves a query parameter value from the URL.
//
// Parameters:
//   - key: The name of the query parameter to retrieve.
//
// Returns:
//   - string: The value of the specified query parameter, or an empty string if not found.
func (c *Ctx) GetReqHeadersAll() map[string][]string {
	return c.Headers
}

// GetHeadersAll returns all HTTP response headers stored in the context.
//
// Returns:
//   - map[string][]string: A map containing all response headers with their values.
func (c *Ctx) GetHeadersAll() map[string][]string {
	return c.Headers
}

// File serves a specific file to the client.
//
// This function trims any trailing "/*" from the provided file path, checks if it is a directory,
// and serves "index.html" if applicable. If the file exists, it is sent as the response.
//
// Parameters:
//   - filePath: The path to the file to be served.
//
// Returns:
//   - error: Always returns nil, as `http.ServeFile` handles errors internally.
func (c *Ctx) File(filePath string) error {
	filePath = strings.TrimSuffix(filePath, "/*")

	if stat, err := os.Stat(filePath); err == nil && stat.IsDir() {
		filePath = filepath.Join(filePath, "index.html")
	}
	http.ServeFile(c.Response, c.Request, filePath)
	return nil
}

// Bind parses and binds the request body to a Go struct.
//
// This function extracts and maps the request body content to the given struct (v).
// It supports various content types and ensures proper deserialization.
//
// Parameters:
//   - v: A pointer to the structure where the request body will be bound.
//
// Returns:
//   - error: An error if parsing fails or if the structure is incompatible with the request body.
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

// Param retrieves the value of a URL parameter corresponding to the given key.
//
// This function searches for a parameter in the request's URL path and returns its value.
// If the parameter is not found, an empty string is returned.
//
// Parameters:
//   - key: The name of the URL parameter to retrieve.
//
// Returns:
//   - string: The value of the requested parameter or an empty string if not found.
func (c *Ctx) Param(key string) string {
	val, ok := c.Params[key]
	if ok {
		return val
	}
	return ""
}

// Body retrieves the request body as a byte slice.
//
// This function returns the raw request body as a slice of bytes ([]byte).
//
// Returns:
//   - []byte: The request body in its raw byte form.
func (c *Ctx) Body() []byte {
	return c.bodyByte
}

// BodyString retrieves the request body as a string.
//
// This function converts the request body from a byte slice into a string format.
//
// Returns:
//   - string: The request body as a string.
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

	buf := acquireJSONBuffer()
	defer releaseJSONBuffer(buf)

	// Exemplo com JSON:
	enc := json.NewEncoder(buf)
	enc.SetIndent(prefix, indent)

	if buf.Len() > 0 && buf.Bytes()[buf.Len()-1] == '\n' {
		buf.Truncate(buf.Len() - 1)
	}

	if err := enc.Encode(v); err != nil {
		return err
	}

	c.writeResponse(buf.Bytes())
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

// Byte writes a byte slice to the HTTP response.
//
// This function writes raw bytes to the response body using writeResponse().
//
// Parameters:
//   - b: The byte slice to be written.
//
// Returns:
//   - error: An error if the response write operation fails.
func (c *Ctx) Byte(b []byte) (err error) {
	return c.writeResponse(b)
}

// Send writes a byte slice to the HTTP response.
//
// This function writes raw bytes to the response body using writeResponse().
//
// Parameters:
//   - b: The byte slice to be written.
//
// Returns:
//   - error: An error if the response write operation fails.
func (c *Ctx) Send(b []byte) (err error) {
	return c.writeResponse(b)
}

// SendString writes a string to the HTTP response.
//
// This function converts the given string into a byte slice and writes it to the response body.
//
// Parameters:
//   - s: The string to be written.
//
// Returns:
//   - error: An error if the response write operation fails.
func (c *Ctx) SendString(s string) error {
	return c.writeResponse([]byte(s))
}

// String writes a string to the HTTP response.
//
// This function converts the given string into a byte slice and writes it to the response body.
//
// Parameters:
//   - s: The string to be written.
//
// Returns:
//   - error: An error if the response write operation fails.
func (c *Ctx) String(s string) error {
	return c.writeResponse([]byte(s))
}

// SendFile writes a file to the HTTP response as a byte slice.
//
// This function writes the provided byte slice (representing a file) to the response body.
//
// Parameters:
//   - file: The file content as a byte slice.
//
// Returns:
//   - error: An error if the response write operation fails.
func (c *Ctx) SendFile(file []byte) error {
	_, err := c.Response.Write(file)
	return err
}

// Set defines an HTTP header in the response.
//
// This function sets the specified HTTP response header to the provided value.
//
// Parameters:
//   - key: The name of the HTTP header to set.
//   - value: The value to assign to the header.
func (c *Ctx) Set(key, value string) {
	c.Response.Header().Set(key, value)
}

// Append adds a value to an HTTP response header.
//
// This function appends a new value to an existing HTTP response header.
//
// Parameters:
//   - key: The name of the HTTP header.
//   - value: The value to append to the header.
func (c *Ctx) Append(key, value string) {
	c.Response.Header().Add(key, value)
}

// Accepts sets the "Accept" header in the HTTP response.
//
// This function assigns a specific accept type to the HTTP response header "Accept."
//
// Parameters:
//   - acceptType: The MIME type to set in the "Accept" header.
//
// Returns:
//   - *Ctx: The current context instance for method chaining.
func (c *Ctx) Accepts(acceptType string) *Ctx {
	c.Response.Header().Set("Accept", acceptType)
	return c
}

// Status sets the HTTP status code of the response.
//
// This function assigns a specific HTTP status code to the response.
//
// Parameters:
//   - status: The HTTP status code to set.
//
// Returns:
//   - *Ctx: The current context instance for method chaining.
func (c *Ctx) Status(status int) *Ctx {
	c.resStatus = status
	return c
}

//MultipartForm

// FormFileLimit sets the maximum allowed upload size.
//
// This function configures the maximum file upload size for multipart form-data requests.
//
// Parameters:
//   - limit: A string representing the maximum file size (e.g., "10MB").
//
// Returns:
//   - error: An error if the limit value is invalid.
func (c *Ctx) FormFileLimit(limit string) error {
	size, err := parseSize(limit)
	if err != nil {
		return err
	}
	c.uploadFileSize = size
	return nil
}

// FormFile processes an uploaded file and returns its details.
//
// This function retrieves the first uploaded file for the specified form field.
//
// Parameters:
//   - fieldName: The name of the form field containing the uploaded file.
//
// Returns:
//   - *UploadedFile: A struct containing the uploaded file details.
//   - error: An error if no file is found or the retrieval fails.
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

// fileWrapper wraps a bytes.Reader and adds a Close() method.
//
// This struct implements io.ReadCloser, allowing it to be used as a multipart.File.
// It ensures that the file can be read multiple times without losing data.
//
// Fields:
//   - *bytes.Reader: A reader that holds the file content in memory.
type fileWrapper struct {
	*bytes.Reader
}

// Close satisfies the io.ReadCloser interface.
//
// This function does nothing since the file is stored in memory
// and does not require explicit closing.
//
// Returns:
//   - error: Always returns nil.
func (fw *fileWrapper) Close() error {
	return nil
}

// FormFiles retrieves all uploaded files for the given field name.
//
// This function extracts all files uploaded in a multipart form request.
//
// Parameters:
//   - fieldName: The name of the form field containing the uploaded files.
//
// Returns:
//   - []*UploadedFile: A slice containing details of the uploaded files.
//   - error: An error if no files are found or the retrieval fails.
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

// MultipartForm provides access to the raw multipart form data.
//
// This function parses and retrieves the multipart form from the request.
//
// Returns:
//   - *multipart.Form: A pointer to the multipart form data.
//   - error: An error if parsing fails.
func (c *Ctx) MultipartForm() (*multipart.Form, error) {
	if err := c.Request.ParseMultipartForm(c.uploadFileSize); err != nil {
		return nil, err
	}
	return c.Request.MultipartForm, nil
}

// FormValue retrieves a form value by its key.
//
// This function parses the form data and returns the value of the specified field.
//
// Parameters:
//   - key: The name of the form field.
//
// Returns:
//   - string: The value of the requested field, or an empty string if not found.
func (c *Ctx) FormValue(key string) string {
	// Checks if the Content-Type is multipart
	if c.Request.Header.Get("Content-Type") == "multipart/form-data" {
		_ = c.Request.ParseMultipartForm(c.uploadFileSize) // Force correct processing
	} else {
		_ = c.Request.ParseForm() // For application/x-www-form-urlencoded
	}
	return c.Request.FormValue(key)
}

// FormValues retrieves all form values as a map.
//
// This function parses the form data and returns all form values.
//
// Returns:
//   - map[string][]string: A map of form field names to their corresponding values.
func (c *Ctx) FormValues() map[string][]string {
	// Checks if the Content-Type is multipart
	if c.Request.Header.Get("Content-Type") == "multipart/form-data" {
		_ = c.Request.ParseMultipartForm(c.uploadFileSize) // Required to process multipart
	} else {
		_ = c.Request.ParseForm() // Processes application/x-www-form-urlencoded
	}
	return c.Request.Form
}
