package goquick

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
)

type Ctx struct {
	Response     http.ResponseWriter
	Request      *http.Request
	resStatus    int
	MoreRequests int
	bodyByte     []byte
	JsonStr      string
	Headers      map[string][]string
	Params       map[string]string
	Query        map[string]string
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

// func (c *Ctx) GetHeaders(key string, defaultValue ...string) (err error) {

// }

// func (c *Ctx) GetReqHeaders(key string, defaultValue ...string) (err error) {

// }

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
	c.Response.Header().Set("Content-Type", ContentTypeAppJSON)
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
