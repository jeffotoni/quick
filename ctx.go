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
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeffotoni/quick/rand"
)

const ACCUMULATED_CONTEXT_KEY = "__quick_context_data__"

type ContextBuilder struct {
	ctx *Ctx
}

// ContextDataCallback is called when SetContext is invoked
var ContextDataCallback func(*http.Request, map[string]any)

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

	handlers     []HandlerFunc   // list of handlers for this route
	handlerIndex int             // current position in handlers stack
	wroteHeader  bool            // valid writeResponse
	Context      context.Context // custom context
}

// SetCtx allows you to override the default context used in c.Ctx()
func (c *Ctx) SetCtx(ctx context.Context) {
	c.Context = ctx
}

// Ctx returns the active context for this request.
// It returns the custom context if SetCtx was used; otherwise defaults to Request.Context()
func (c *Ctx) Ctx() context.Context {
	if c.Context != nil {
		return c.Context
	}
	return c.Request.Context()
}

// responseWriter is a custom wrapper around http.ResponseWriter that prevents
// duplicate WriteHeader calls, which would otherwise cause "superfluous
// response.WriteHeader" warnings or errors.
//
// The wrapper tracks whether WriteHeader has already been called and silently
// ignores subsequent attempts. This is particularly useful in complex middleware
// chains or when handlers accidentally attempt multiple response writes.
//
// Fields:
//   - ResponseWriter: The underlying http.ResponseWriter being wrapped
//   - statusCode: The HTTP status code that was written
//   - wroteHeader: Flag indicating whether WriteHeader has been called
//
// Example Usage:
//
//	rw := &responseWriter{ResponseWriter: w}
//	rw.WriteHeader(200) // First call - executes normally
//	rw.WriteHeader(500) // Second call - silently ignored
type responseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

// WriteHeader writes the HTTP status code to the response.
//
// This method ensures that WriteHeader is only called once on the underlying
// http.ResponseWriter. If WriteHeader has already been called, subsequent calls
// are silently ignored to prevent "superfluous response.WriteHeader" errors.
//
// Example:
//
//	w.WriteHeader(200)  // Writes status 200
//	w.WriteHeader(500)  // Ignored - header already written
func (w *responseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return // Silently ignore duplicate calls
	}
	w.statusCode = code
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(code)
}

// Write writes the response body data to the client.
//
// This method ensures that if WriteHeader hasn't been explicitly called before
// writing the body, it automatically calls WriteHeader with http.StatusOK (200).
// This matches the standard http.ResponseWriter behavior where the first Write
// call implicitly sends a 200 status if no status was set.
//
// Example:
//
//	n, err := w.Write([]byte("Hello, World!"))
//	// If WriteHeader wasn't called, automatically sends 200 status first
func (w *responseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}

func (c *Ctx) HTML(name string, data interface{}, layouts ...string) error {
	if c.App == nil {
		return errors.New("App is nil")
	}

	cfg := c.App.GetConfig()
	//log.Println("DEBUG", "cfg.Views:", cfg.Views)

	if cfg.Views == nil {
		return errors.New("template engine not configured")
	}

	return cfg.Views.Render(c.Response, name, data, layouts...)
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

// Get retrieves a specific header value from the request.
//
// Parameters:
//   - key: The name of the header to retrieve.
//
// Returns:
//   - string: The value of the specified header, or an empty string if not found.
func (c *Ctx) Get(key string) string {
	return c.Request.Header.Get(key)
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

// GetHeaders retrieves all headers from the incoming HTTP request.
//
// This method provides direct access to the request headers, allowing
// middleware and handlers to inspect and modify header values.
//
// Example Usage:
//
//	q.Get("/", func(c *quick.Ctx) error {
//	    headers := c.GetHeaders()
//	    return c.Status(200).JSON(headers)
//	})
//
// Returns:
//   - http.Header: A map containing all request headers.
func (c *Ctx) GetHeaders() http.Header {
	return c.Request.Header
}

// RemoteIP extracts the client's IP address from the request.
//
// If the request's `RemoteAddr` contains a port (e.g., "192.168.1.100:54321"),
// this method extracts only the IP part. If extraction fails, it returns
// the full `RemoteAddr` as a fallback.
//
// Example Usage:
//
//	q.Get("/", func(c *quick.Ctx) error {
//	    return c.Status(200).SendString("Client IP: " + c.RemoteIP())
//	})
//
// Returns:
//   - string: The client's IP address. If extraction fails, returns `RemoteAddr`.
func (c *Ctx) RemoteIP() string {
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}

// IP is alias RemoteIP
func (c *Ctx) IP() string {
	return c.RemoteIP()
}

// Method retrieves the HTTP method of the current request.
//
// This method returns the HTTP method as a string, such as "GET", "POST", "PUT",
// "DELETE", etc. It is useful for middleware and route handlers to differentiate
// between request types.
//
// Example Usage:
//
//	q.Use(func(c *quick.Ctx) error {
//	    if c.Method() == "POST" {
//	        return c.Status(403).SendString("POST requests are not allowed")
//	    }
//	    return c.Next()
//	})
//
// Returns:
//   - string: The HTTP method (e.g., "GET", "POST", "PUT").
func (c *Ctx) Method() string {
	return c.Request.Method
}

// OriginalURI retrieves the original request URI sent by the client.
//
// This method returns the full unprocessed request URI, including the path and
// optional query string, exactly as it was sent by the client. It can be useful
// for logging, debugging, or routing decisions that depend on the raw URI.
//
// Example Usage:
//
//	q.Use(func(c *quick.Ctx) error {
//	    uri := c.OriginalURI()
//	    if strings.HasPrefix(uri, "/admin") {
//	        return c.Status(401).SendString("Unauthorized access to admin area")
//	    }
//	    return c.Next()
//	})
//
// Returns:
//   - string: The raw request URI (e.g., "/api/v1/resource?id=123").
func (c *Ctx) OriginalURI() string {
	return c.Request.RequestURI
}

// Path retrieves the URL path of the incoming HTTP request.
//
// This method extracts the path component from the request URL, which
// is useful for routing and request handling.
//
// Example Usage:
//
//	q.Get("/info", func(c *quick.Ctx) error {
//	    return c.Status(200).SendString("Requested Path: " + c.Path())
//	})
//
// Returns:
//   - string: The path component of the request URL (e.g., "/v1/user").
func (c *Ctx) Path() string {
	return c.Request.URL.Path
}

// Host returns the host name from the HTTP request.
//
// This method extracts the host from `c.Request.Host`. If the request includes
// a port number (e.g., "localhost:3000"), it returns the full host including
// the port.
//
// Example Usage:
//
//	q.Get("/", func(c *quick.Ctx) error {
//	    return c.Status(200).SendString("Host: " + c.Host())
//	})
//
// Returns:
//   - string: The host name from the request.
func (c *Ctx) Host() string {
	if c.Request == nil {
		return ""
	}
	return c.Request.Host
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
	if c.App.hasEmbed {
		filePath = strings.TrimPrefix(filePath, "./")

		if strings.HasSuffix(filePath, "/") || filePath == "" {
			filePath = filepath.Join(filePath, "index.html")
		} else if info, err := fs.Stat(c.App.embedFS, filePath); err == nil && info.IsDir() {
			filePath = filepath.Join(filePath, "index.html")
		}

		data, err := c.App.embedFS.ReadFile(filePath)
		if err != nil {
			return c.Status(StatusNotFound).SendString("File not found")
		}

		contentType := http.DetectContentType(data)
		c.Response.Header().Set("Content-Type", contentType)
		c.Response.Write(data)
		return nil
	}

	// Caso contrário, usa o disco como antes
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
		c.resStatus = StatusOK
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

func (c *Ctx) Del(key string) {
	c.Response.Header().Del(key)
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

// Add defines an HTTP header in the response.
//
// This function sets the specified HTTP response header to the provided value.
//
// Parameters:
//   - key: The name of the HTTP header to set.
//   - value: The value to assign to the header.
func (c *Ctx) Add(key, value string) {
	c.Response.Header().Add(key, value)
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
		fileContentType := detectUploadedFileContentType(handler, buf.Bytes())

		// Append file details
		uploadedFiles = append(uploadedFiles, &UploadedFile{
			File:      fileCopy,
			Multipart: handler,
			Info: FileInfo{
				Filename:    handler.Filename,
				Size:        handler.Size,
				ContentType: fileContentType,
				Bytes:       buf.Bytes(),
			},
		})
	}

	return uploadedFiles, nil
}

func detectUploadedFileContentType(handler *multipart.FileHeader, fileBytes []byte) string {
	filename := ""
	partContentType := ""
	if handler != nil {
		filename = handler.Filename
		partContentType = strings.TrimSpace(handler.Header.Get("Content-Type"))
	}

	ext := strings.ToLower(filepath.Ext(filename))
	extContentType := contentTypeByExtension(filename)

	sniffed := http.DetectContentType(fileBytes)
	sniffBase := baseMediaType(sniffed)
	partBase := baseMediaType(partContentType)

	if sniffBase != "" && sniffBase != "application/zip" && !isGenericMediaType(sniffBase) {
		return sniffed
	}
	if partContentType != "" && partBase != "" && partBase != "application/zip" && !isGenericMediaType(partBase) {
		return partContentType
	}

	if sniffBase == "application/zip" || partBase == "application/zip" {
		if kind := officeKindFromZip(fileBytes); kind != officeUnknown {
			if info, ok := officeExtTypes[ext]; ok && info.kind == kind {
				return info.mime
			}
			if fallback := defaultOfficeMime(kind); fallback != "" {
				return fallback
			}
		}
		return sniffed
	}

	if extContentType != "" && (isGenericMediaType(sniffBase) || isGenericMediaType(partBase)) {
		return extContentType
	}

	return sniffed
}

type officeKind uint8

const (
	officeUnknown officeKind = iota
	officeWord
	officeExcel
	officePowerPoint
)

type officeExtInfo struct {
	mime string
	kind officeKind
}

var officeExtTypes = map[string]officeExtInfo{
	// Word (legacy)
	".doc": {mime: "application/msword", kind: officeWord},
	".dot": {mime: "application/msword", kind: officeWord},
	// Word (OOXML)
	".docx": {mime: "application/vnd.openxmlformats-officedocument.wordprocessingml.document", kind: officeWord},
	".docm": {mime: "application/vnd.ms-word.document.macroEnabled.12", kind: officeWord},
	".dotx": {mime: "application/vnd.openxmlformats-officedocument.wordprocessingml.template", kind: officeWord},
	".dotm": {mime: "application/vnd.ms-word.template.macroEnabled.12", kind: officeWord},

	// Excel (legacy)
	".xls": {mime: "application/vnd.ms-excel", kind: officeExcel},
	".xlt": {mime: "application/vnd.ms-excel", kind: officeExcel},
	".xla": {mime: "application/vnd.ms-excel", kind: officeExcel},
	// Excel (OOXML)
	".xlsx": {mime: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", kind: officeExcel},
	".xlsm": {mime: "application/vnd.ms-excel.sheet.macroEnabled.12", kind: officeExcel},
	".xltx": {mime: "application/vnd.openxmlformats-officedocument.spreadsheetml.template", kind: officeExcel},
	".xltm": {mime: "application/vnd.ms-excel.template.macroEnabled.12", kind: officeExcel},
	".xlam": {mime: "application/vnd.ms-excel.addin.macroEnabled.12", kind: officeExcel},
	".xlsb": {mime: "application/vnd.ms-excel.sheet.binary.macroEnabled.12", kind: officeExcel},

	// PowerPoint (legacy)
	".ppt": {mime: "application/vnd.ms-powerpoint", kind: officePowerPoint},
	".pps": {mime: "application/vnd.ms-powerpoint", kind: officePowerPoint},
	".pot": {mime: "application/vnd.ms-powerpoint", kind: officePowerPoint},
	".ppa": {mime: "application/vnd.ms-powerpoint", kind: officePowerPoint},
	// PowerPoint (OOXML)
	".pptx": {mime: "application/vnd.openxmlformats-officedocument.presentationml.presentation", kind: officePowerPoint},
	".pptm": {mime: "application/vnd.ms-powerpoint.presentation.macroEnabled.12", kind: officePowerPoint},
	".potx": {mime: "application/vnd.openxmlformats-officedocument.presentationml.template", kind: officePowerPoint},
	".potm": {mime: "application/vnd.ms-powerpoint.template.macroEnabled.12", kind: officePowerPoint},
	".ppsx": {mime: "application/vnd.openxmlformats-officedocument.presentationml.slideshow", kind: officePowerPoint},
	".ppsm": {mime: "application/vnd.ms-powerpoint.slideshow.macroEnabled.12", kind: officePowerPoint},
	".sldx": {mime: "application/vnd.openxmlformats-officedocument.presentationml.slide", kind: officePowerPoint},
	".sldm": {mime: "application/vnd.ms-powerpoint.slide.macroEnabled.12", kind: officePowerPoint},
	".ppam": {mime: "application/vnd.ms-powerpoint.addin.macroEnabled.12", kind: officePowerPoint},
}

func officeKindFromZip(fileBytes []byte) officeKind {
	zr, err := zip.NewReader(bytes.NewReader(fileBytes), int64(len(fileBytes)))
	if err != nil {
		return officeUnknown
	}

	for _, f := range zr.File {
		name := f.Name
		switch {
		case strings.HasPrefix(name, "word/"):
			return officeWord
		case strings.HasPrefix(name, "xl/"):
			return officeExcel
		case strings.HasPrefix(name, "ppt/"):
			return officePowerPoint
		}
	}

	return officeUnknown
}

func defaultOfficeMime(kind officeKind) string {
	switch kind {
	case officeWord:
		return officeExtTypes[".docx"].mime
	case officeExcel:
		return officeExtTypes[".xlsx"].mime
	case officePowerPoint:
		return officeExtTypes[".pptx"].mime
	default:
		return ""
	}
}

func baseMediaType(contentType string) string {
	contentType = strings.TrimSpace(contentType)
	if contentType == "" {
		return ""
	}
	if idx := strings.IndexByte(contentType, ';'); idx >= 0 {
		contentType = contentType[:idx]
	}
	return strings.ToLower(strings.TrimSpace(contentType))
}

func isGenericMediaType(baseMediaType string) bool {
	switch baseMediaType {
	case "", "application/octet-stream", "binary/octet-stream":
		return true
	default:
		return false
	}
}

func contentTypeByExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return ""
	}

	if info, ok := officeExtTypes[ext]; ok {
		return info.mime
	}

	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		return mimeType
	}

	return ""
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

// Redirect sends an HTTP redirect response to the client.
//
// It sets the "Location" header to the given URL and returns a plain text
// message indicating the redirection. By default, it uses HTTP status 302 (Found),
// but an optional custom status code can be provided (e.g., 301, 307, 308).
//
// Example usage:
//
//	c.Redirect("https://quick.com")
//	c.Redirect("/new-path", 301)
func (c *Ctx) Redirect(location string, code ...int) error {
	status := StatusFound // 302
	if len(code) > 0 {
		status = code[0]
	}
	c.Set("Location", location)
	c.Status(status)
	return c.String("Redirecting to " + location)
}

// responseWritten checks if the response has already been written.
//
// This is useful to prevent duplicate writes or to decide whether to send
// default responses such as 404.
func (c *Ctx) responseWritten() bool {
	if rw, ok := c.Response.(interface{ Written() bool }); ok {
		return rw.Written()
	}
	// Fallback: assume false
	return false
}

// setHandlers sets the chain of handlers for the current request.
//
// This is used internally by the router to assign the matched middleware
// and final handler to the context.
func (c *Ctx) setHandlers(handlers []HandlerFunc) {
	c.handlers = handlers
	c.handlerIndex = -1
}

// Next executes the next handler in the chain.
//
// If no more handlers are available and no response has been written,
// it automatically sends a 404 response.
//
// Usage:
//
//	func myMiddleware(c *quick.Ctx) error {
//	    err := c.Next()
//	    return err
//	}
func (c *Ctx) Next() error {
	c.handlerIndex++
	if c.handlerIndex < len(c.handlers) {
		handler := c.handlers[c.handlerIndex]
		return handler(c)
	}

	// No more handlers — if no response has been written, send 404
	if !c.responseWritten() {
		NotFound(c.Response, c.Request)
	}

	return nil
}

// Logger creates a new context builder for adding key-value pairs.
//
// Usage:
//
//	c.Logger().Str("service", "user-service").Int("userID", 123)
func (c *Ctx) Logger() *ContextBuilder {
	return &ContextBuilder{ctx: c}
}

// SetContext creates a new context builder for adding key-value pairs.
//
// Usage:
//
//	c.SetContext().Str("service", "user-service").Int("userID", 123)
func (c *Ctx) SetContext() *ContextBuilder {
	return &ContextBuilder{ctx: c}
}

// Str adds a string value to the context and returns the builder for chaining.
//
// Parameters:
//   - key: Context key name
//   - value: String value to store
//
// Usage:
//
//	c.SetContext().Str("service", "user-service").Str("function", "createUser")
func (cb *ContextBuilder) Str(key, value string) *ContextBuilder {
	if key != "" && value != "" {
		existingData := cb.ctx.getAccumulatedData()
		existingData[key] = value
		cb.ctx.applyAccumulatedContext(existingData)
	}
	return cb
}

// Int adds an integer value to the context as a string and returns the builder for chaining.
//
// Parameters:
//   - key: Context key name
//   - value: Integer value to store (converted to string)
//
// Usage:
//
//	c.SetContext().Int("userID", 12345).Int("attempts", 3)
func (cb *ContextBuilder) Int(key string, value int) *ContextBuilder {
	if key != "" {
		existingData := cb.ctx.getAccumulatedData()
		existingData[key] = value
		cb.ctx.applyAccumulatedContext(existingData)
	}
	return cb
}

// Bool adds a boolean value to the context as a string and returns the builder for chaining.
//
// Parameters:
//   - key: Context key name
//   - value: Boolean value to store (converted to "true" or "false")
//
// Usage:
//
//	c.SetContext().Bool("authenticated", true).Bool("admin", false)
func (cb *ContextBuilder) Bool(key string, value bool) *ContextBuilder {
	if key != "" {
		existingData := cb.ctx.getAccumulatedData()
		existingData[key] = value
		cb.ctx.applyAccumulatedContext(existingData)
	}
	return cb
}

// SetTraceID sets a trace ID header and adds it to the context.
//
// Parameters:
//   - key: Header name (e.g., "X-Trace-ID")
//   - val: Trace ID value
//
// Returns a ContextBuilder for chaining additional context data.
//
// Usage:
//
//	c.SetTraceID("X-Trace-ID", traceID).Str("service", "user-service")
func (c *Ctx) SetTraceID(key, val string) *ContextBuilder {
	// Set trace ID header
	c.Set(key, val)
	existingData := c.getAccumulatedData()
	existingData[key] = val
	c.applyAccumulatedContext(existingData)
	return &ContextBuilder{ctx: c}
}

// getAccumulatedData retrieves all accumulated context data or creates a new map.
func (c *Ctx) getAccumulatedData() map[string]any {
	if existing := c.Request.Context().Value(ACCUMULATED_CONTEXT_KEY); existing != nil {
		if data, ok := existing.(map[string]any); ok {
			copy := make(map[string]any)
			for k, v := range data {
				copy[k] = v
			}
			return copy
		}
	}
	return make(map[string]any)
}

// applyAccumulatedContext applies all accumulated data to the request context.
func (c *Ctx) applyAccumulatedContext(allData map[string]any) {
	ctx := c.Request.Context()

	// Store all accumulated data in special key
	ctx = context.WithValue(ctx, ACCUMULATED_CONTEXT_KEY, allData)

	for key, value := range allData {
		if key != "" && value != "" {
			ctx = context.WithValue(ctx, key, value)
		}
	}
	c.Request = c.Request.WithContext(ctx)

	// Also update the request in logger middleware if it exists
	if rw, ok := c.Response.(interface{ SetRequest(*http.Request) }); ok {
		rw.SetRequest(c.Request)
	}
}

// GetAllContextData returns all accumulated context data.
//
// Usage:
//
//	allData := c.GetAllContextData()
//	fmt.Printf("Context: %+v", allData)
func (c *Ctx) GetAllContextData() map[string]any {
	return c.getAccumulatedData()
}

// GetTraceID retrieves or generates a trace ID for the current request.
//
// This method first attempts to get an existing trace ID from the request headers.
// If no trace ID is found, it automatically generates a new one using a random
// trace ID generator. This ensures that every request has a unique trace ID for
// distributed tracing purposes.
//
// Parameters:
//   - nameTraceID: The header name to look for the trace ID (e.g., "X-Trace-ID", "Trace-ID")
//
// Returns:
//   - traceID: The existing trace ID from headers, or a newly generated one if none exists
//
// This method is typically used at the beginning of request processing to ensure
// traceability across your application stack. The generated trace ID can then be
// used with SetTraceContext or forwarded to downstream services.
//
// Usage:
//
//	func myHandler(c *quick.Ctx) error {
//	    traceID := c.GetTraceID("X-Trace-ID")
//	    c.SetTraceContext("X-Trace-ID", traceID, "user-service", "createUser")
//
//	    // Continue with your handler logic
//	   return c.Status(200).String(traceID)
//	}
func (c *Ctx) GetTraceID(nameTraceID string) (traceID string) {
	// Generate TraceID for this request
	traceID = c.Get(nameTraceID)
	if traceID == "" {
		traceID = rand.TraceID()
	}
	return
}

// Flusher is an alias for http.Flusher, used to abstract the standard net/http flush interface
// within the Quick framework. This definition allows you to hide
// the direct dependency on the http package at points of use, keeping the code cleaner,
// semantic, and consistent with the framework's architecture.
//
// Flusher is primarily used in connections that support real-time writing,
// such as Server-Sent Events (SSE), where it is necessary to send parts of the response
// immediately to the client.
//
// Usage example:
// flusher, ok := c.Flusher()
// if ok {
// flusher.Flush()
// }
type Flusher = http.Flusher

// Flusher returns the underlying http.Flusher if available.
// This is useful for Server-Sent Events (SSE) and streaming responses.
//
// Returns:
//   - http.Flusher: The flusher instance, or nil if not supported.
//   - bool: true if flusher is available, false otherwise.
//
// Example Usage:
//
//	q.Get("/events", func(c *quick.Ctx) error {
//	    // Set SSE headers manually
//	    c.Set("Content-Type", "text/event-stream")
//	    c.Set("Cache-Control", "no-cache")
//	    c.Set("Connection", "keep-alive")
//
//	    flusher, ok := c.Flusher()
//	    if !ok {
//	        return c.Status(500).SendString("Streaming not supported")
//	    }
//
//	    for i := 0; i < 10; i++ {
//	        fmt.Fprintf(c.Response, "data: Message %d\n\n", i)
//	        flusher.Flush()
//	        time.Sleep(time.Second)
//	    }
//	    return nil
//	})
func (c *Ctx) Flusher() (http.Flusher, bool) {
	flusher, ok := c.Response.(http.Flusher)
	return flusher, ok
}

// Flush flushes any buffered data to the client.
// Returns an error if flushing is not supported by the ResponseWriter.
//
// This is particularly useful for Server-Sent Events (SSE) and streaming responses.
//
// Example Usage:
//
//	q.Get("/stream", func(c *quick.Ctx) error {
//	    c.Set("Content-Type", "text/event-stream")
//	    c.Set("Cache-Control", "no-cache")
//	    c.Set("Connection", "keep-alive")
//
//	    for i := 0; i < 5; i++ {
//	        fmt.Fprintf(c.Response, "data: %d\n\n", i)
//	        if err := c.Flush(); err != nil {
//	            return err
//	        }
//	        time.Sleep(time.Second)
//	    }
//	    return nil
//	})
func (c *Ctx) Flush() error {
	if flusher, ok := c.Response.(http.Flusher); ok {
		flusher.Flush()
		return nil
	}
	return errors.New("flushing not supported")
}

// Hijack implements http.Hijacker interface for WebSocket support.
//
// This method allows the responseWriter to support connection hijacking,
// which is essential for protocols like WebSockets that need direct access
// to the underlying TCP connection.
//
// Returns:
//   - net.Conn: The underlying network connection
//   - *bufio.ReadWriter: Buffered reader/writer for the connection
//   - error: An error if hijacking is not supported by the underlying ResponseWriter
//
// Example:
//
//	conn, bufrw, err := w.Hijack()
//	if err != nil {
//	    log.Fatal("WebSocket upgrade failed")
//	}
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, errors.New("hijacking not supported")
}

// Push implements http.Pusher interface for HTTP/2 Server Push support.
//
// This method enables HTTP/2 server push, allowing the server to proactively
// send resources to the client before they are requested. This can improve
// page load performance by reducing round trips.
//
// Parameters:
//   - target: The path of the resource to push (e.g., "/style.css")
//   - opts: Push options, can be nil for defaults
//
// Returns:
//   - error: An error if push is not supported or if the push fails
//
// Example:
//
//	err := w.Push("/static/app.js", nil)
//	if err != nil {
//	    log.Println("HTTP/2 push not available")
//	}
func (w *responseWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return errors.New("push not supported")
}

// Pusher returns the underlying http.Pusher if available for HTTP/2 server push.
// This is useful for proactively sending resources to the client before they are requested,
// reducing latency and improving page load performance.
//
// Returns:
//   - http.Pusher: The pusher instance, or nil if not supported.
//   - bool: true if pusher is available, false otherwise.
//
// Example Usage:
//
//	q.Get("/", func(c *quick.Ctx) error {
//	    pusher, ok := c.Pusher()
//	    if ok {
//	        pusher.Push("/static/style.css", nil)
//	        pusher.Push("/static/app.js", nil)
//	    }
//	    return c.Status(200).SendString("<html>...</html>")
//	})
func (c *Ctx) Pusher() (http.Pusher, bool) {
	pusher, ok := c.Response.(http.Pusher)
	return pusher, ok
}
