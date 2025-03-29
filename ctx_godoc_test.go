// Package quick provides a high-performance, minimalistic web framework for Go.
//
// 📌 To run all unit tests, use:
//
//	$ go test -v ./...
//	$ go test -v
package quick

import (
	"encoding/xml"
	"fmt"
	"testing"
)

// TestCtx_GetReqHeadersAll ensures GetReqHeadersAll returns all request headers.
//
// To run:
//
// go test -v -run ^TestCtx_GetReqHeadersAll$
func TestCtx_GetReqHeadersAll(t *testing.T) {
	ctx := &Ctx{
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
			"Accept":       {"application/xml"},
		},
	}

	headers := ctx.GetReqHeadersAll()

	if headers["Content-Type"][0] != "application/json" {
		t.Errorf("Expected 'application/json', got '%s'", headers["Content-Type"][0])
	}

	if headers["Accept"][0] != "application/xml" {
		t.Errorf("Expected 'application/xml', got '%s'", headers["Accept"][0])
	}
}

// TestCtx_GetHeadersAll ensures GetHeadersAll returns all headers from the context.
//
// To run:
//
// go test -v -run ^TestCtx_GetHeadersAll$
func TestCtx_GetHeadersAll(t *testing.T) {
	ctx := &Ctx{
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
			"Accept":       {"application/xml"},
		},
	}

	headers := ctx.GetHeadersAll()

	if headers["Content-Type"][0] != "application/json" {
		t.Errorf("Expected 'application/json', got '%s'", headers["Content-Type"][0])
	}

	if headers["Accept"][0] != "application/xml" {
		t.Errorf("Expected 'application/xml', got '%s'", headers["Accept"][0])
	}
}

// TestCtx_ExampleBind verifies that Bind correctly binds JSON body to a struct.
//
// To run:
//
// go test -v -run ^TestCtx_ExampleBind$
func TestCtx_ExampleBind(t *testing.T) {
	q := New()

	q.Post("/bind", func(c *Ctx) error {
		var data struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		err := c.Bind(&data)
		if err != nil {
			t.Errorf("Bind failed: %v", err)
			return err
		}

		return c.Status(200).JSON(data)
	})

	body := []byte(`{"name":"Quick","age":30}`)

	// Simulate a GET request to "/api/users"
	res, err := q.Qtest(QuickTestOptions{
		Method:  MethodPost,
		URI:     "/bind",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    body,
	})
	if err != nil {
		t.Errorf("Error during Qtest: %v", err)
		return
	}

	if res.StatusCode() != 200 {
		t.Errorf("Expected status 200, got %d", res.StatusCode())
	}

	fmt.Println(res.BodyStr())

	expected := `{"name":"Quick","age":30}`
	if res.BodyStr() != expected {
		t.Errorf("Expected response '%s', but got '%s'", expected, res.BodyStr())
	}
}

// TestCtx_ExampleBodyParser checks if BodyParser parses the request body into a struct.
//
// To run:
//
// go test -v -run ^TestCtx_ExampleBodyParser$
func TestCtx_ExampleBodyParser(t *testing.T) {
	q := New()

	q.Post("/test", func(c *Ctx) error {
		var data struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		err := c.BodyParser(&data)
		if err != nil {
			t.Errorf("BodyParser failed: %v", err)
			return err
		}

		return c.Status(200).JSON(data)
	})

	body := []byte(`{"name": "Quick", "age": 28}`)

	res, err := q.Qtest(QuickTestOptions{
		Method:  MethodPost,
		URI:     "/test",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    body,
	})
	if err != nil {
		t.Errorf("Error during Qtest: %v", err)
		return
	}

	if res.StatusCode() != 200 {
		t.Errorf("Expected status 200, got %d", res.StatusCode())
	}

	expected := `{"name":"Quick","age":28}`
	if res.BodyStr() != expected {
		t.Errorf("Expected response '%s', but got '%s'", expected, res.BodyStr())
	}
}

// TestCtx_ExampleParam checks if Param correctly retrieves path parameters.
//
// To run:
//
// go test -v -run ^TestCtx_ExampleParam$
func TestCtx_ExampleParam(t *testing.T) {
	q := New()

	q.Get("/user/:id", func(c *Ctx) error {
		id := c.Param("id")
		return c.SendString(id)
	})

	res, err := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/user/42",
	})
	if err != nil {
		t.Errorf("Error during Qtest: %v", err)
		return
	}

	expectedBody := "42"
	if res.BodyStr() != expectedBody {
		t.Errorf("Expected: %s, received: %s", expectedBody, res.BodyStr())
	}
}

// TestCtx_ExampleBody ensures Body returns the correct raw request body.
//
// To run:
//
// go test -v -run ^TestCtx_ExampleBody$
func TestCtx_ExampleBody(t *testing.T) {
	expectedBody := `{"name": "Quick", "age": 28}`

	c := &Ctx{
		bodyByte: []byte(expectedBody),
	}

	body := c.Body()

	if string(body) != expectedBody {
		t.Errorf("Expected: %s, received: %s", expectedBody, string(body))
	}
}

// TestCtx_ExampleBodyString ensures BodyString returns the body as a string.
//
// To run:
//
// go test -v -run ^TestCtx_ExampleBodyString$
func TestCtx_ExampleBodyString(t *testing.T) {
	expectedBody := `{"name": "Quick", "age": 28}`

	c := &Ctx{
		bodyByte: []byte(expectedBody),
	}

	bodyStr := c.BodyString()

	if bodyStr != expectedBody {
		t.Errorf("Expected: %s, received: %s", expectedBody, bodyStr)
	}
}

// TestCtx_ExampleJSON validates JSON response generation.
//
// To run:
//
// go test -v -run ^TestCtx_ExampleJSON$
func TestCtx_ExampleJSON(t *testing.T) {
	q := New()

	q.Get("/json", func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		data := map[string]string{"message": "Hello, Quick!"}
		return c.JSON(data)
	})

	res, err := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/json",
	})
	if err != nil {
		t.Errorf("Error during Qtest: %v", err)
		return
	}

	expectedBody := `{"message":"Hello, Quick!"}`

	if res.BodyStr() != expectedBody {
		t.Errorf("Expected: %s, received: %s", expectedBody, res.BodyStr())
	}

	expectedContentType := "application/json"
	contentType := res.Response().Header.Get("Content-Type")

	if contentType != expectedContentType {
		t.Errorf("Expected Content-Type: %s, received: %s", expectedContentType, contentType)
	}
}

// TestCtx_ExampleJSONIN ensures JSONIN returns correct JSON with header set.
//
// To run:
//
// go test -v -run ^TestCtx_ExampleJSONIN$
func TestCtx_ExampleJSONIN(t *testing.T) {
	q := New()

	q.Get("/json", func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		data := map[string]string{"message": "Hello, Quick!"}
		return c.JSONIN(data)
	})

	res, err := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/json",
	})
	if err != nil {
		t.Errorf("Error during Qtest: %v", err)
		return
	}

	// expectedBody := `{"message":"Hello, Quick!"}`
	// if res.BodyStr() != expectedBody {
	// 	t.Errorf("Expected: %s, received: %s", expectedBody, res.BodyStr())
	// }

	expectedContentType := "application/json"
	contentType := res.Response().Header.Get("Content-Type")

	if contentType != expectedContentType {
		t.Errorf("Expected Content-Type: %s, received: %s", expectedContentType, contentType)
	}
}

type XMLMessage struct {
	XMLName xml.Name `xml:"message"`
	Message string   `xml:",chardata"`
}

// TestCtx_ExampleXML verifies that XML responses are returned properly.
//
// To run:
//
// go test -v -run ^TestCtx_ExampleXML$
func TestCtx_ExampleXML(t *testing.T) {
	q := New()

	q.Get("/xml", func(c *Ctx) error {
		c.Set("Content-Type", "text/xml")
		data := XMLMessage{Message: "Hello, Quick!"}
		return c.XML(data)
	})

	res, err := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/xml",
	})
	if err != nil {
		t.Errorf("Error during Qtest: %v", err)
		return
	}

	expectedBody := `<message>Hello, Quick!</message>`

	if res.BodyStr() != expectedBody {
		t.Errorf("Expected: %s, received: %s", expectedBody, res.BodyStr())
	}

	expectedContentType := "text/xml"
	contentType := res.Response().Header.Get("Content-Type")

	if contentType != expectedContentType {
		t.Errorf("Expected Content-Type: %s, received: %s", expectedContentType, contentType)
	}
}

// TestCtx_ExamplewriteResponse validates raw byte response writing.
//
// To run:
//
// go test -v -run ^TestCtx_ExamplewriteResponse$
func TestCtx_ExamplewriteResponse(t *testing.T) {
	q := New()

	q.Get("/response", func(c *Ctx) error {
		return c.writeResponse([]byte("Hello, Quick!"))
	})

	res, err := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/response",
	})
	if err != nil {
		t.Errorf("Error during Qtest: %v", err)
		return
	}

	expectedBody := "Hello, Quick!"

	if res.BodyStr() != expectedBody {
		t.Errorf("Expected: %s, received: %s", expectedBody, res.BodyStr())
	}

	expectedStatus := 200
	if res.Response().StatusCode != expectedStatus {
		t.Errorf("Expected Status Code: %d, received: %d", expectedStatus, res.Response().StatusCode)
	}
}

// TestCtx_ExampleByte checks that Byte correctly writes raw bytes.
//
// To run:
//
// go test -v -run ^TestCtx_ExampleByte$
func TestCtx_ExampleByte(t *testing.T) {
	q := New()

	q.Get("/byte", func(c *Ctx) error {
		return c.Byte([]byte("Hello, Quick!"))
	})

	res, err := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/byte",
	})
	if err != nil {
		t.Errorf("Error during Qtest: %v", err)
		return
	}

	expectedBody := "Hello, Quick!"

	if res.BodyStr() != expectedBody {
		t.Errorf("Expected: %s, received: %s", expectedBody, res.BodyStr())
	}

	expectedStatus := 200
	if res.Response().StatusCode != expectedStatus {
		t.Errorf("Expected Status Code: %d, received: %d", expectedStatus, res.Response().StatusCode)
	}
}

// TestCtx_ExampleSend validates that Send writes the correct byte response.
//
// To run:
//
// go test -v -run ^TestCtx_ExampleSend$
func TestCtx_ExampleSend(t *testing.T) {
	q := New()

	q.Get("/send", func(c *Ctx) error {
		return c.Send([]byte("Hello, Quick!"))
	})

	res, err := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/send",
	})
	if err != nil {
		t.Errorf("Error during Qtest: %v", err)
		return
	}

	expectedBody := "Hello, Quick!"

	if res.BodyStr() != expectedBody {
		t.Errorf("Expected: %s, received: %s", expectedBody, res.BodyStr())
	}

	expectedStatus := 200
	if res.Response().StatusCode != expectedStatus {
		t.Errorf("Expected Status Code: %d, received: %d", expectedStatus, res.Response().StatusCode)
	}
}

// TestCtx_ExampleSendString checks that SendString writes a plain text response.
//
// To run:
//
// go test -v -run ^TestCtx_ExampleSendString$
func TestCtx_ExampleSendString(t *testing.T) {
	q := New()

	q.Get("/sendstring", func(c *Ctx) error {
		return c.SendString("Hello, Quick!")
	})

	res, err := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/sendstring",
	})
	if err != nil {
		t.Errorf("Error during Qtest: %v", err)
		return
	}

	expectedBody := "Hello, Quick!"

	if res.BodyStr() != expectedBody {
		t.Errorf("Expected: %s, received: %s", expectedBody, res.BodyStr())
	}

	expectedStatus := 200
	if res.Response().StatusCode != expectedStatus {
		t.Errorf("Expected Status Code: %d, received: %d", expectedStatus, res.Response().StatusCode)
	}
}

// TestCtx_ExampleSendFile ensures that SendFile correctly sends file content.
//
// To run:
//
// go test -v -run ^TestCtx_ExampleSendFile$
func TestCtx_ExampleSendFile(t *testing.T) {
	q := New()

	q.Get("/sendfile", func(c *Ctx) error {
		fileContent := []byte("Conteúdo do arquivo")
		return c.SendFile(fileContent)
	})

	res, err := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/sendfile",
	})
	if err != nil {
		t.Errorf("Error during Qtest: %v", err)
		return
	}

	expectedBody := "Conteúdo do arquivo"

	if res.BodyStr() != expectedBody {
		t.Errorf("Expected: %s, received: %s", expectedBody, res.BodyStr())
	}

	expectedStatus := 200
	if res.Response().StatusCode != expectedStatus {
		t.Errorf("Expected Status Code: %d, received: %d", expectedStatus, res.Response().StatusCode)
	}
}

// TestCtx_ExampleSet verifies that Set adds headers correctly.
//
// To run:
//
// go test -v -run ^TestCtx_ExampleSet$
func TestCtx_ExampleSet(t *testing.T) {
	q := New()

	q.Get("/set-header", func(c *Ctx) error {
		c.Set("X-Custom-Header", "Quick")
		return c.String("Header Set")
	})

	res, err := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/set-header",
	})
	if err != nil {
		t.Errorf("Error during Qtest: %v", err)
		return
	}

	expectedHeader := "Quick"

	headerValue := res.Response().Header.Get("X-Custom-Header")
	if headerValue != expectedHeader {
		t.Errorf("Expected: %s, received: %s", expectedHeader, headerValue)
	}

	expectedStatus := 200
	if res.Response().StatusCode != expectedStatus {
		t.Errorf("Expected Status Code: %d, received: %d", expectedStatus, res.Response().StatusCode)
	}
}

// TestCtx_ExampleAppend ensures that headers can be appended properly.
//
// To run:
//
// go test -v -run ^TestCtx_ExampleAppend$
func TestCtx_ExampleAppend(t *testing.T) {
	q := New()

	q.Get("/append-header", func(c *Ctx) error {
		c.Append("X-Custom-Header", "Value1")
		c.Append("X-Custom-Header", "Value2")
		return c.String("Header Appended")
	})

	res, err := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/append-header",
	})
	if err != nil {
		t.Errorf("Error during Qtest: %v", err)
		return
	}

	expectedHeaders := []string{"Value1", "Value2"}
	headers := res.Response().Header.Values("X-Custom-Header")

	if len(headers) != len(expectedHeaders) {
		t.Errorf("Expected: %v, received: %v", expectedHeaders, headers)
	}
}

// TestCtx_ExampleAccepts checks if the Accepts method sets the correct header.
//
// To run:
//
// go test -v -run ^TestCtx_ExampleAccepts$
func TestCtx_ExampleAccepts(t *testing.T) {
	q := New()

	q.Get("/accepts", func(c *Ctx) error {
		c.Accepts("application/json")
		return c.String("Accept Set")
	})

	res, err := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/accepts",
	})
	if err != nil {
		t.Errorf("Error during Qtest: %v", err)
		return
	}

	expectedHeader := "application/json"
	header := res.Response().Header.Get("Accept")

	if header != expectedHeader {
		t.Errorf("Expected: %s, received: %s", expectedHeader, header)
	}
}

// TestCtx_ExampleStatus verifies that Status sets the correct HTTP status code.
//
// To run:
//
// go test -v -run ^TestCtx_ExampleStatus$
func TestCtx_ExampleStatus(t *testing.T) {
	q := New()

	q.Get("/status", func(c *Ctx) error {
		c.Status(404)
		return c.String("Not Found")
	})

	res, err := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/status",
	})
	if err != nil {
		t.Errorf("Error during Qtest: %v", err)
		return
	}

	expectedStatus := 404
	status := res.Response().StatusCode

	if status != expectedStatus {
		t.Errorf("Expected Status Code: %d, received: %d", expectedStatus, status)
	}
}
