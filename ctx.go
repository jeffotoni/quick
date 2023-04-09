package quick

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
)

type Ctx struct {
	Response  http.ResponseWriter
	Request   *http.Request
	resStatus int
	bodyByte  []byte
	JsonStr   string
	Headers   map[string][]string
	Params    map[string]string
	Query     map[string]string
}

func (c *Ctx) Bind(v interface{}) (err error) {
	return extractBind(c, v)
}

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

func (c *Ctx) Param(key string) string {
	val, ok := c.Params[key]
	if ok {
		return val
	}
	return ""
}

func (c *Ctx) Body() []byte {
	return c.bodyByte
}

func (c *Ctx) BodyString() string {
	return string(c.bodyByte)
}

func (c *Ctx) JSON(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	c.Response.Header().Set("Content-Type", ContentTypeAppJSON)
	return c.writeResponse(b)
}

func (c *Ctx) XML(v interface{}) error {
	b, err := xml.Marshal(v)
	if err != nil {
		return err
	}
	c.Response.Header().Set("Content-Type", ContentTypeTextXML)
	return c.writeResponse(b)
}

func (c *Ctx) writeResponse(b []byte) error {
	if c.resStatus != 0 {
		c.Response.WriteHeader(c.resStatus)
	}
	_, err := c.Response.Write(b)
	return err
}

func (c *Ctx) Byte(b []byte) (err error) {
	return c.writeResponse(b)
}

func (c *Ctx) Send(b []byte) (err error) {
	return c.writeResponse(b)
}

func (c *Ctx) SendString(s string) error {
	return c.writeResponse([]byte(s))

}

func (c *Ctx) String(s string) error {
	return c.writeResponse([]byte(s))
}

func (c *Ctx) SendFile(file []byte) error {
	_, err := c.Response.Write(file)
	return err
}

func (c *Ctx) Set(key, value string) {
	c.Response.Header().Set(key, value)
}

func (c *Ctx) Append(key, value string) {
	c.Response.Header().Add(key, value)
}

func (c *Ctx) Accepts(acceptType string) *Ctx {
	c.Response.Header().Set("Accept", acceptType)
	return c
}

func (c *Ctx) Status(status int) *Ctx {
	c.resStatus = status
	return c
}
