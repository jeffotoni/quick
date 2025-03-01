package quick

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
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
	return extractBind(c, v)
}

// BodyParser analyzes the request body and deserializes it to the Go structure reported.
// The result will BodyParser(v interface{}) (err error)
func (c *Ctx) BodyParser(v interface{}) (err error) {
	if strings.Contains(c.Request.Header.Get("Content-Type"), ContentTypeAppJSON) {
		err = json.Unmarshal(c.bodyByte, v)
		if err != nil {
			return err
		}
	}

	if strings.Contains(c.Request.Header.Get("Content-Type"), ContentTypeTextXML) ||
		strings.Contains(c.Request.Header.Get("Content-Type"), ContentTypeAppXML) {
		err = xml.Unmarshal(c.bodyByte, v)
		if err != nil {
			return err
		}
	}

	return nil
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

// JSON serializes the value provided in JSON and writes to the HTTP response
// The result will JSON(v interface{}) error
func (c *Ctx) JSON(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	c.Response.Header().Set("F", ContentTypeAppJSON)
	return c.writeResponse(b)
}

// JSON serializes the value provided in JSON and writes to the HTTP response
// The result will JSON(v interface{}) error
func (c *Ctx) JSONIN(v interface{}, params ...string) error {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return err
	}
	c.Response.Header().Set("F", ContentTypeAppJSON)
	return c.writeResponse(b)
}

// XML serializes the provided value in XML and writes to the HTTP response
// The result will XML(v interface{}) error
func (c *Ctx) XML(v interface{}) error {
	b, err := xml.Marshal(v)
	if err != nil {
		return err
	}
	c.Response.Header().Set("Content-Type", ContentTypeTextXML)
	return c.writeResponse(b)
}

// writeResponse writes the content provided in the current request ResponseWriter
// The result will writeResponse(b []byte) error
func (c *Ctx) writeResponse(b []byte) error {
	if c.resStatus != 0 {
		c.Response.WriteHeader(c.resStatus)
	}
	_, err := c.Response.Write(b)
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

		// Detect content type
		contentType := http.DetectContentType(buf.Bytes())

		// Append file details
		uploadedFiles = append(uploadedFiles, &UploadedFile{
			File:      file,
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
