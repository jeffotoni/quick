package quick

import (
	"encoding/xml"
	"fmt"
	"testing"
)

// This function is named ExampleCtx_GetReqHeadersAll()
// it with the Examples type.
func ExampleCtx_GetReqHeadersAll() {
	q := New()

	q.Get("/headers", func(c *Ctx) error {
		headers := c.GetReqHeadersAll()
		fmt.Println(headers["Content-Type"])
		fmt.Println(headers["Accept"])
		return nil
	})

	res, _ := q.QuickTest("GET", "/headers", map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/xml",
	}, nil)

	fmt.Println(res.BodyStr())

	// Out put:
	// [application/json]
	// [application/xml]
}

// This function is named ExampleCtx_GetHeadersAll()
// it with the Examples type.
func ExampleCtx_GetHeadersAll() {
	q := New()

	q.Get("/headers", func(c *Ctx) error {
		headers := c.GetHeadersAll()
		fmt.Println(headers["Content-Type"])
		fmt.Println(headers["Accept"])
		return nil
	})

	res, _ := q.QuickTest("GET", "/headers", map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/xml",
	}, nil)

	fmt.Println(res.BodyStr())

	// Out put:
	// [application/json]
	// [application/xml]
}

// This function is named ExampleCtx_Bind()
// it with the Examples type.
func ExampleCtx_Bind() {
	q := New()

	q.Post("/bind", func(c *Ctx) error {
		var data struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		err := c.Bind(&data)
		if err != nil {
			fmt.Println("Erro ao fazer bind:", err)
			return err
		}

		fmt.Println(data.Name, data.Age)
		return nil
	})

	body := []byte(`{"name": "Quick", "age": 30}`)

	res, _ := q.QuickTest("POST", "/bind", map[string]string{"Content-Type": "application/json"}, body)

	fmt.Println(res.BodyStr())

	// Out put: Quick 30
}

// This function is named ExampleCtx_BodyParser()
// it with the Examples type.
func ExampleCtx_BodyParser() {
	q := New()

	q.Post("/parse", func(c *Ctx) error {
		var data struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		err := c.BodyParser(&data)
		if err != nil {
			fmt.Println("Erro ao analisar o corpo:", err)
			return err
		}

		fmt.Println(data.Name, data.Age)
		return nil
	})

	body := []byte(`{"name": "Quick", "age": 28}`)

	res, _ := q.QuickTest("POST", "/parse", map[string]string{"Content-Type": "application/json"}, body)

	fmt.Println(res.BodyStr())

	// Out put: Quick 28
}

// This function is named ExampleCtx_Param()
// it with the Examples type.
func ExampleCtx_Param() {
	q := New()

	q.Get("/user/:id", func(c *Ctx) error {
		id := c.Param("id")
		return c.SendString(id)
	})

	res, _ := q.QuickTest("GET", "/user/42", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put: 42
}

// This function is named ExampleCtx_Body()
// it with the Examples type.
func ExampleCtx_Body() {
	c := &Ctx{
		bodyByte: []byte(`{"name": "Quick", "age": 28}`),
	}

	body := c.Body()

	fmt.Println(string(body))

	// Out put: {"name": "Quick", "age": 28}
}

// This function is named ExampleCtx_Body()
// it with the Examples type.
func ExampleCtx_BodyString() {
	c := &Ctx{
		bodyByte: []byte(`{"name": "Quick", "age": 28}`),
	}

	bodyStr := c.BodyString()

	fmt.Println(bodyStr)

	// Out put: {"name": "Quick", "age": 28}
}

// This function is named ExampleCtx_JSON()
// it with the Examples type.
func ExampleCtx_JSON() {
	q := New()

	q.Get("/json", func(c *Ctx) error {
		data := map[string]string{"message": "Hello, Quick!"}
		return c.JSON(data)
	})

	res, _ := q.QuickTest("GET", "/json", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put: {"message":"Hello, Quick!"}
}

// This function is named ExampleCtx_XML()
// it with the Examples type.
func ExampleCtx_XML() {
	q := New()

	q.Get("/xml", func(c *Ctx) error {
		data := struct {
			Message string `xml:"message"`
		}{
			Message: "Hello, Quick!",
		}
		return c.XML(data)
	})

	res, _ := q.QuickTest("GET", "/xml", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put:<message>Hello, Quick!</message>
}

// This function is named ExampleCtx_writeResponse()
// it with the Examples type.
func ExampleCtx_writeResponse() {
	q := New()

	q.Get("/response", func(c *Ctx) error {
		return c.writeResponse([]byte("Hello, Quick!"))
	})

	res, _ := q.QuickTest("GET", "/response", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put: Hello, Quick!
}

// This function is named ExampleCtx_Byte()
// it with the Examples type.
func ExampleCtx_Byte() {
	q := New()

	q.Get("/byte", func(c *Ctx) error {
		return c.Byte([]byte("Hello, Quick!"))
	})

	res, _ := q.QuickTest("GET", "/byte", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put: Hello, Quick!
}

// This function is named ExampleCtx_Send()
// it with the Examples type.
func ExampleCtx_Send() {
	q := New()

	q.Get("/send", func(c *Ctx) error {
		return c.Send([]byte("Hello, Quick!"))
	})

	res, _ := q.QuickTest("GET", "/send", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put: Hello, Quick!
}

// This function is named ExampleCtx_SendString()
// it with the Examples type.
func ExampleCtx_SendString() {
	q := New()

	q.Get("/sendstring", func(c *Ctx) error {
		return c.SendString("Hello, Quick!")
	})

	res, _ := q.QuickTest("GET", "/sendstring", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put:	Hello, Quick!
}

// This function is named ExampleCtx_String()
// it with the Examples type.
func ExampleCtx_String() {
	q := New()

	q.Get("/string", func(c *Ctx) error {
		return c.String("Hello, Quick!")
	})

	res, _ := q.QuickTest("GET", "/string", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put: Hello, Quick!
}

// This function is named ExampleCtx_SendFile()
// it with the Examples type.
func ExampleCtx_SendFile() {
	q := New()

	q.Get("/sendfile", func(c *Ctx) error {
		fileContent := []byte("Conteúdo do arquivo")
		return c.SendFile(fileContent)
	})

	res, _ := q.QuickTest("GET", "/sendfile", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put: Conteúdo do arquivo
}

// This function is named ExampleCtx_Set()
// it with the Examples type.
func ExampleCtx_Set() {
	q := New()

	q.Get("/set-header", func(c *Ctx) error {
		c.Set("X-Custom-Header", "Quick")
		return c.String("Header Set")
	})

	res, _ := q.QuickTest("GET", "/set-header", nil, nil)

	fmt.Println(res.Response().Header.Get("X-Custom-Header"))

	// Out put: Quick
}

// This function is named ExampleCtx_Append()
// it with the Examples type.
func ExampleCtx_Append() {
	q := New()

	q.Get("/append-header", func(c *Ctx) error {
		c.Append("X-Custom-Header", "Value1")
		c.Append("X-Custom-Header", "Value2")
		return c.String("Header Appended")
	})

	res, _ := q.QuickTest("GET", "/append-header", nil, nil)

	fmt.Println(res.Response().Header.Values("X-Custom-Header"))

	// Out put: [Value1 Value2]
}

// This function is named ExampleCtx_Accepts()
// it with the Examples type.
func ExampleCtx_Accepts() {
	q := New()

	q.Get("/accepts", func(c *Ctx) error {
		c.Accepts("application/json")
		return c.String("Accept Set")
	})

	res, _ := q.QuickTest("GET", "/accepts", nil, nil)

	fmt.Println(res.Response().Header.Get("Accept"))

	// Out put: application/json
}

// This function is named ExampleCtx_Status()
// it with the Examples type.
func ExampleCtx_Status() {
	q := New()

	q.Get("/status", func(c *Ctx) error {
		c.Status(404)
		return c.String("Not Found")
	})

	res, _ := q.QuickTest("GET", "/status", nil, nil)

	fmt.Println(res.Response().StatusCode)

	// Out put: 404
}

// go test -v -run ^TestCtx_GetReqHeadersAll
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

// go test -v -run ^TestCtx_GetHeadersAll
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

// go test -v -run ^TestCtx_ExampleBind
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

	body := []byte(`{"name": "Quick", "age": 30}`)

	res, err := q.QuickTest("POST", "/bind", map[string]string{"Content-Type": "application/json"}, body)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	if res.StatusCode() != 200 {
		t.Errorf("Expected status 200, got %d", res.StatusCode())
	}

	expected := `{"name":"Quick","age":30}`
	if res.BodyStr() != expected {
		t.Errorf("Expected response '%s', but got '%s'", expected, res.BodyStr())
	}
}

// go test -v -run ^TestCtx_ExampleBodyParser
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

	res, err := q.QuickTest("POST", "/test", map[string]string{"Content-Type": "application/json"}, body)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	if res.StatusCode() != 200 {
		t.Errorf("Expected status 200, got %d", res.StatusCode())
	}

	expected := `{"name":"Quick","age":28}`
	if res.BodyStr() != expected {
		t.Errorf("Expected response '%s', but got '%s'", expected, res.BodyStr())
	}
}

// go test -v -run ^TestCtx_ExampleParam
func TestCtx_ExampleParam(t *testing.T) {
	q := New()

	q.Get("/user/:id", func(c *Ctx) error {
		id := c.Param("id")
		return c.SendString(id)
	})

	res, err := q.QuickTest("GET", "/user/42", nil, nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	expectedBody := "42"
	if res.BodyStr() != expectedBody {
		t.Errorf("Expected: %s, received: %s", expectedBody, res.BodyStr())
	}
}

// go test -v -run ^TestCtx_ExampleBody
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

// go test -v -run ^TestCtx_ExampleBodyString
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

// go test -v -run ^TestCtx_ExampleJSON
func TestCtx_ExampleJSON(t *testing.T) {
	q := New()

	q.Get("/json", func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		data := map[string]string{"message": "Hello, Quick!"}
		return c.JSON(data)
	})

	res, err := q.QuickTest("GET", "/json", nil, nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
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

// go test -v -run ^TestCtx_ExampleJSONIN
func TestCtx_ExampleJSONIN(t *testing.T) {
	q := New()

	q.Get("/json", func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		data := map[string]string{"message": "Hello, Quick!"}
		return c.JSONIN(data)
	})

	res, err := q.QuickTest("GET", "/json", nil, nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
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

// go test -v -run ^TestCtx_ExampleXML
func TestCtx_ExampleXML(t *testing.T) {
	q := New()

	q.Get("/xml", func(c *Ctx) error {
		data := XMLMessage{Message: "Hello, Quick!"}
		return c.XML(data)
	})

	res, err := q.QuickTest("GET", "/xml", nil, nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
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

// go test -v -run ^TestCtx_ExampleXML
func TestCtx_ExamplewriteResponse(t *testing.T) {
	q := New()

	q.Get("/response", func(c *Ctx) error {
		return c.writeResponse([]byte("Hello, Quick!"))
	})

	res, err := q.QuickTest("GET", "/response", nil, nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
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

// go test -v -run ^TestCtx_ExampleByte
func TestCtx_ExampleByte(t *testing.T) {
	q := New()

	q.Get("/byte", func(c *Ctx) error {
		return c.Byte([]byte("Hello, Quick!"))
	})

	res, err := q.QuickTest("GET", "/byte", nil, nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
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

// go test -v -run ^TestCtx_ExampleSend
func TestCtx_ExampleSend(t *testing.T) {
	q := New()

	q.Get("/send", func(c *Ctx) error {
		return c.Send([]byte("Hello, Quick!"))
	})

	res, err := q.QuickTest("GET", "/send", nil, nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
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

// go test -v -run ^TestCtx_ExampleSendString
func TestCtx_ExampleSendString(t *testing.T) {
	q := New()

	q.Get("/sendstring", func(c *Ctx) error {
		return c.SendString("Hello, Quick!")
	})

	res, err := q.QuickTest("GET", "/sendstring", nil, nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
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

// go test -v -run ^TestCtx_ExampleSendFile
func TestCtx_ExampleSendFile(t *testing.T) {
	q := New()

	q.Get("/sendfile", func(c *Ctx) error {
		fileContent := []byte("Conteúdo do arquivo")
		return c.SendFile(fileContent)
	})

	res, err := q.QuickTest("GET", "/sendfile", nil, nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
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

// go test -v -run ^TestCtx_ExampleSet
func TestCtx_ExampleSet(t *testing.T) {
	q := New()

	q.Get("/set-header", func(c *Ctx) error {
		c.Set("X-Custom-Header", "Quick")
		return c.String("Header Set")
	})

	res, err := q.QuickTest("GET", "/set-header", nil, nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
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

// go test -v -run ^TestCtx_ExampleAppend
func TestCtx_ExampleAppend(t *testing.T) {
	q := New()

	q.Get("/append-header", func(c *Ctx) error {
		c.Append("X-Custom-Header", "Value1")
		c.Append("X-Custom-Header", "Value2")
		return c.String("Header Appended")
	})

	res, err := q.QuickTest("GET", "/append-header", nil, nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	expectedHeaders := []string{"Value1", "Value2"}
	headers := res.Response().Header.Values("X-Custom-Header")

	if len(headers) != len(expectedHeaders) {
		t.Errorf("Expected: %v, received: %v", expectedHeaders, headers)
	}
}

// go test -v -run ^TestCtx_ExampleAccepts
func TestCtx_ExampleAccepts(t *testing.T) {
	q := New()

	q.Get("/accepts", func(c *Ctx) error {
		c.Accepts("application/json")
		return c.String("Accept Set")
	})

	res, err := q.QuickTest("GET", "/accepts", nil, nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	expectedHeader := "application/json"
	header := res.Response().Header.Get("Accept")

	if header != expectedHeader {
		t.Errorf("Expected: %s, received: %s", expectedHeader, header)
	}
}

// go test -v -run ^TestCtx_ExampleStatus
func TestCtx_ExampleStatus(t *testing.T) {
	q := New()

	q.Get("/status", func(c *Ctx) error {
		c.Status(404)
		return c.String("Not Found")
	})

	res, err := q.QuickTest("GET", "/status", nil, nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	expectedStatus := 404
	status := res.Response().StatusCode

	if status != expectedStatus {
		t.Errorf("Expected Status Code: %d, received: %d", expectedStatus, status)
	}
}
