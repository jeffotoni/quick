package ctx

// signature
// func (c *Ctx) Accepts(offers ...string) string {}
// func (c *Ctx) AcceptsCharsets(offers ...string) string {}
// func (c *Ctx) AcceptsEncodings(offers ...string) string {}
// func (c *Ctx) AcceptsLanguages(offers ...string) string {}

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"encoding/xml"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"mime/multipart"
// 	"net"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"strings"
// )

//principal struct 

// type Ctx struct {
// 	Response       http.ResponseWriter
// 	Request        *http.Request
// 	resStatus      int
// 	MoreRequests   int
// 	bodyByte       []byte
// 	JsonStr        string
// 	Headers        map[string]string
// 	Params         map[string]string
// 	Query          map[string]string
// 	uploadFileSize int64
// 	App            *Quick
// }

//context manipulation

// func (c *Ctx) SetCtx(ctx context.Context) {
// 	*c = Ctx{
// 		Request: c.Request.WithContext(ctx),
// 		Response: c.Response,
// 	}
// }
// func (c *Ctx) Ctx() context.Context {
// 	return c.Request.Context()
// }

// reading requisition
//body acces

// func (c *Ctx) Body() []byte {
// 	if c.bodyByte == nil {
// 		body, _ := io.ReadAll(c.Request.Body)
// 		c.bodyByte = body
// 	}
// 	return c.bodyByte
// }
// func (c *Ctx) BodyString() string {
// 	if c.JsonStr == "" {
// 		c.JsonStr = string(c.Body())
// 	}
// 	return c.JsonStr
// }

//parameters querry and headers

// func (c *Ctx) Param(key string) string {
// 	return c.Params[key]
// }
// func (c *Ctx) QueryParam(key string) string {
// 	return c.Query[key]
// }
// func (c *Ctx) GetHeader(key string) string {
// 	return c.Request.Header.Get(key)
// }
// func (c *Ctx) GetHeadersAll() map[string]string {
// 	return c.Headers
// }


//client response

// func (c *Ctx) JSON(v interface{}) {
// 	b, _ := json.Marshal(v)
// 	c.writeResponse(b)
// }
// func (c *Ctx) JSONIN(v interface{}, prefix ...string) {
// 	b, _ := json.MarshalIndent(v, "", "  ")
// 	c.writeResponse(b)
// }
// func (c *Ctx) XML(v interface{}) {
// 	b, _ := xml.Marshal(v)
// 	c.writeResponse(b)
// }
// func (c *Ctx) SendString(s string) {
// 	c.writeResponse([]byte(s))
// }
// func (c *Ctx) SendFile(b []byte) {
// 	c.writeResponse(b)
// }

//send system files

// func (c *Ctx) File(filePath string) {
// data, _ := os.ReadFile(filePath)
// c.Response.Write(data)
// }

// header and status manipulation

// func (c *Ctx) SetStatus(status int) {
// 	c.resStatus = status
// }
// func (c *Ctx) Set(key, value string) {
// 	c.Response.Header().Set(key, value)
// }
// func (c *Ctx) Add(key, value string) {
// 	c.Response.Header().Add(key, value)
// }
// func (c *Ctx) Del(key string) {
// 	c.Response.Header().Del(key)
// }
// func (c *Ctx) Accepts(acceptType string) {
// 	c.Set("Accept", acceptType)
// }

// files upload limiter

// func (c *Ctx) FormFileLimit(limit string) {
// 	c.uploadFileSize = ParseFileSize(limit)
// }

//Get file from a field:

// func (c *Ctx) FormFile(fieldName string) (*UploadedFile, error) {
// 	err := c.Request.ParseMultipartForm(c.uploadFileSize)
// 	if err != nil {
// 		return nil, err
// 	}
// 	file, handler, err := c.Request.FormFile(fieldName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()
// 	buf := bytes.NewBuffer(nil)
// 	_, err = io.Copy(buf, file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &UploadedFile{
// 		File: buf.Bytes(),
// 		Info: &FileInfo{
// 			Filename: handler.Filename,
// 			MimeType: handler.Header.Get("Content-Type"),
// 			Size:     handler.Size,
// 		},
// 	}, nil
// }

//assistant structs for files

// type UploadedFile struct {
// 	File []byte
// 	Info *FileInfo
// }

// type FileInfo struct {
// 	Filename string
// 	MimeType string
// 	Size     int64
// }


// Deserialization (Bind / BodyParser)
// Detect type and deserialize:

// func (c *Ctx) Bind(v interface{}) error {
// 	return extractParamsBind(c, v)
// }

// func (c *Ctx) BodyParser(v interface{}) error {
// 	contentType := c.GetHeader("Content-Type")
// 	switch {
// 	case strings.Contains(contentType, "application/json"):
// 		return json.Unmarshal(c.Body(), v)
// 	case strings.Contains(contentType, "application/xml"):
// 		return xml.Unmarshal(c.Body(), v)
// 	default:
// 		return errors.New("unsupported content type")
// 	}
// }

//internal function for write response

// func (c *Ctx) writeResponse(b []byte) {
// 	c.Response.WriteHeader(c.resStatus)
// 	c.Response.Write(b)
// }

//-----------------------------------------------------------------------------------------------------

// Headers

// obtain every headers request
// func TestCtx_GetReqHeadersAll(t *testing.T) { ... }

// obtain every headers ctx
// func TestCtx_GetHeadersAll(t *testing.T) { ... }

// bind and BodyParser
// func TestCtx_ExampleBind(t *testing.T) { ... }

// parsing of body request
// func TestCtx_ExampleBodyParser(t *testing.T) { ... }

//url parameters
// func TestCtx_ExampleParam(t *testing.T) { ... }

//body request
// func TestCtx_ExampleBody(t *testing.T) { ... }

// body request as a string
// func TestCtx_ExampleBodyString(t *testing.T) { ... }

// return JSON on response 
// func TestCtx_ExampleJSON(t *testing.T) { ... }

// return JSON with header 
// func TestCtx_ExampleJSONIN(t *testing.T) { ... }

// return XML on response
// func TestCtx_ExampleXML(t *testing.T) { ... }

// write on response 
// func TestCtx_ExamplewriteResponse(t *testing.T) { ... }

// send bytes as a response 
// func TestCtx_ExampleByte(t *testing.T) { ... }

// send bytes with Sends()
// func TestCtx_ExampleSend(t *testing.T) { ... }

// send string with SendString()
// func TestCtx_ExampleSendString(t *testing.T) { ... }

// send file content with sendfile()
// func TestCtx_ExampleSendFile(t *testing.T) { ... }

// header manipulation
// set custom header 
// func TestCtx_ExampleSet(t *testing.T) { ... }

// add multiple values as a header 
// func TestCtx_ExampleAppend(t *testing.T) { ... }

// set header accept
// func TestCtx_ExampleAccepts(t *testing.T) { ... }

// status response HTTP
// altern http status manual
// func TestCtx_ExampleStatus(t *testing.T) { ... }


