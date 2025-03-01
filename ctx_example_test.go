// Package quick provides a minimalistic and high-performance web framework for Go.
//
// This file contains example-based tests for various functionalities of the Quick framework,
// including request handling, JSON parsing, XML responses, and header manipulations.
//
// These examples demonstrate how to use Quick's context (`Ctx`) methods effectively
package quick

import (
	"errors"
	"fmt"
	"net/http"
)

// This function is named ExampleCtx_GetReqHeadersAll()
// it with the Examples type.
func ExampleCtx_GetReqHeadersAll() {
	q := New()

	q.Get("/headers", func(c *Ctx) error {
		// Retrieve all request headers
		headers := c.GetReqHeadersAll()

		// Print specific headers for demonstration
		fmt.Println(headers["Content-Type"]) // Expected: application/json
		fmt.Println(headers["Accept"])       // Expected: application/xml
		return nil
	})

	// Simulate a GET request with headers
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
		// Retrieve all headers from the request
		headers := c.GetHeadersAll()
		fmt.Println(headers["Content-Type"]) // Expected: application/json
		fmt.Println(headers["Accept"])       // Expected: application/xml
		return nil
	})

	// Simulate a GET request with headers
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
		// Define a struct to map the JSON request body
		var data struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		// Parse JSON body into struct
		err := c.Bind(&data)
		if err != nil {
			fmt.Println("Error in Bind:", err)
			return err
		}

		// Print extracted data
		fmt.Println(data.Name, data.Age)
		return nil
	})

	// Simulate a POST request with JSON data
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
		// Define a struct for JSON parsing
		var data struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		// Parse request body into the struct
		err := c.BodyParser(&data)
		if err != nil {
			fmt.Println("Erro ao analisar o corpo:", err)
			return err
		}

		// Print parsed data
		fmt.Println(data.Name, data.Age)
		return nil
	})

	// Simulate a POST request with JSON data
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
		// Retrieve "id" parameter from the URL path
		id := c.Param("id")
		return c.SendString(id)
	})

	// Simulate a GET request with a path parameter
	res, _ := q.QuickTest("GET", "/user/42", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put: 42
}

// This function is named ExampleCtx_Body()
// it with the Examples type.
func ExampleCtx_Body() {
	// Create a new context with a simulated request body
	c := &Ctx{
		bodyByte: []byte(`{"name": "Quick", "age": 28}`),
	}

	// Retrieve the request body as a byte slice
	body := c.Body()

	// Print the request body as a string
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
		// Define JSON response data
		data := map[string]string{"message": "Hello, Quick!"}
		return c.JSON(data)
	})

	// Simulate a GET request and retrieve JSON response
	res, _ := q.QuickTest("GET", "/json", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put: {"message":"Hello, Quick!"}
}

// This function is named ExampleCtx_XML()
// it with the Examples type.
func ExampleCtx_XML() {
	q := New()

	q.Get("/xml", func(c *Ctx) error {
		// Define XML response structure
		data := struct {
			Message string `xml:"message"`
		}{
			Message: "Hello, Quick!",
		}
		return c.XML(data)
	})

	// Simulate a GET request and retrieve XML response
	res, _ := q.QuickTest("GET", "/xml", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put:<message>Hello, Quick!</message>
}

// This function is named ExampleCtx_writeResponse()
// it with the Examples type.
func ExampleCtx_writeResponse() {
	q := New()

	q.Get("/response", func(c *Ctx) error {
		// Directly write raw byte response
		return c.writeResponse([]byte("Hello, Quick!"))
	})

	// Simulate a GET request
	res, _ := q.QuickTest("GET", "/response", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put: Hello, Quick!
}

// This function is named ExampleCtx_Byte()
// it with the Examples type.
func ExampleCtx_Byte() {
	q := New()

	q.Get("/byte", func(c *Ctx) error {
		// Send raw byte array in response
		return c.Byte([]byte("Hello, Quick!"))
	})

	// Simulate a GET request
	res, _ := q.QuickTest("GET", "/byte", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put: Hello, Quick!
}

// This function is named ExampleCtx_Send()
// it with the Examples type.
func ExampleCtx_Send() {
	q := New()

	q.Get("/send", func(c *Ctx) error {
		// Send raw bytes in response
		return c.Send([]byte("Hello, Quick!"))
	})

	// Simulate a GET request
	res, _ := q.QuickTest("GET", "/send", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put: Hello, Quick!
}

// This function is named ExampleCtx_SendString()
// it with the Examples type.
func ExampleCtx_SendString() {
	q := New()

	q.Get("/sendstring", func(c *Ctx) error {
		// Send string response
		return c.SendString("Hello, Quick!")
	})

	// Simulate a GET request
	res, _ := q.QuickTest("GET", "/sendstring", nil, nil)
	fmt.Println(res.BodyStr())

	// Out put:	Hello, Quick!
}

// This function is named ExampleCtx_String()
// it with the Examples type.
func ExampleCtx_String() {
	q := New()

	q.Get("/string", func(c *Ctx) error {
		// Return a simple string response
		return c.String("Hello, Quick!")
	})

	// Simulate a GET request
	res, _ := q.QuickTest("GET", "/string", nil, nil)
	fmt.Println(res.BodyStr())

	// Out put: Hello, Quick!
}

// This function is named ExampleCtx_SendFile()
// it with the Examples type.
func ExampleCtx_SendFile() {
	q := New()

	q.Get("/sendfile", func(c *Ctx) error {
		// Simulate sending a file as a response
		fileContent := []byte("Conteúdo do arquivo")
		return c.SendFile(fileContent)
	})

	// Simulate a GET request
	res, _ := q.QuickTest("GET", "/sendfile", nil, nil)

	fmt.Println(res.BodyStr())

	// Out put: Conteúdo do arquivo
}

// This function is named ExampleCtx_Set()
// it with the Examples type.
func ExampleCtx_Set() {
	q := New()

	q.Get("/set-header", func(c *Ctx) error {
		// Set a custom response header
		c.Set("X-Custom-Header", "Quick")
		return c.String("Header Set")
	})

	// Simulate a GET request
	res, _ := q.QuickTest("GET", "/set-header", nil, nil)

	fmt.Println(res.Response().Header.Get("X-Custom-Header"))

	// Out put: Quick
}

// This function is named ExampleCtx_Append()
// it with the Examples type.
func ExampleCtx_Append() {
	q := New()

	q.Get("/append-header", func(c *Ctx) error {
		// Append multiple values to a custom header
		c.Append("X-Custom-Header", "Value1")
		c.Append("X-Custom-Header", "Value2")
		return c.String("Header Appended")
	})

	// Simulate a GET request
	res, _ := q.QuickTest("GET", "/append-header", nil, nil)

	fmt.Println(res.Response().Header.Values("X-Custom-Header"))

	// Out put: [Value1 Value2]
}

// This function is named ExampleCtx_Accepts()
// it with the Examples type.
func ExampleCtx_Accepts() {
	q := New()

	q.Get("/accepts", func(c *Ctx) error {
		// Set Accept header
		c.Accepts("application/json")
		return c.String("Accept Set")
	})

	// Simulate a GET request
	res, _ := q.QuickTest("GET", "/accepts", nil, nil)

	fmt.Println(res.Response().Header.Get("Accept"))

	// Out put: application/json
}

// This function is named ExampleCtx_Status()
// it with the Examples type.
func ExampleCtx_Status() {
	q := New()

	q.Get("/status", func(c *Ctx) error {
		// Set status code to 404
		c.Status(404)
		return c.String("Not Found")
	})

	// Simulate a GET request
	res, _ := q.QuickTest("GET", "/status", nil, nil)

	fmt.Println(res.Response().StatusCode)

	// Out put: 404
}

// This function is named ExampleCtx_File()
// it with the Examples type.
func ExampleCtx_File() {
	// Creating a Quick instance
	q := New()

	// Defining a route that serves a specific file
	q.Get("/file", func(c *Ctx) error {
		return c.File("quick.txt") // Serves an existing file
	})

	// Simulating a request to test the route
	res, _ := q.QuickTest("GET", "/file", nil)

	// Printing the expected response
	fmt.Println("Status:", res.StatusCode())

	// Out put: Status: 200
}

// This function is named ExampleCtx_JSONIN()
// it with the Examples type.
func ExampleCtx_JSONIN() {
	// Creating a Quick instance
	q := New()

	// Defining a route that returns a structured JSON
	q.Get("/json", func(c *Ctx) error {
		data := map[string]string{"message": "Hello, Quick!"}
		return c.JSONIN(data)
	})

	// Simulating a request to test the JSON response
	res, _ := q.QuickTest("GET", "/json", nil)

	// Printing the expected response
	fmt.Println("Status:", res.StatusCode())
	fmt.Println("Body:", res.BodyStr())

	// Out put: Status: 200
	// Body:
	// {
	//  "message": "Hello, Quick!"
	// }
}

// This function is named ExampleCtx_FormFileLimit()
// it with the Examples type.
func ExampleCtx_FormFileLimit() {
	// Creating a new context
	ctx := &Ctx{}

	// Setting a file upload limit to 5MB
	err := ctx.FormFileLimit("5MB")

	// Checking if an error occurred while setting the limit
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Upload limit set to:", ctx.uploadFileSize)
	}

	// Out put: Upload limit set to: 5242880
}

// This function is named ExampleCtx_FormFile()
// it with the Examples type.
func ExampleCtx_FormFile() {
	// Simulated uploaded file
	uploadedFile := &UploadedFile{
		Info: FileInfo{
			Filename:    "quick.txt",
			Size:        1024,
			ContentType: "text/plain",
			Bytes:       []byte("File content"),
		},
	}

	// Mocking the FormFiles function externally
	mockFormFiles := func(fieldName string) ([]*UploadedFile, error) {
		if fieldName == "file" {
			return []*UploadedFile{uploadedFile}, nil
		}
		return nil, errors.New("file not found")
	}

	// Calling the mocked function instead of modifying `ctx`
	files, err := mockFormFiles("file")

	// Handling the result
	if err != nil {
		fmt.Println("Error:", err)
	} else if len(files) > 0 {
		fmt.Println("Received file:", files[0].FileName())
	}

	// Out put: Received file: quick.txt
}

// This function is named ExampleCtx_FormFiles()
// it with the Examples type.
func ExampleCtx_FormFiles() {
	// Simulating multiple uploaded files
	uploadedFiles := []*UploadedFile{
		{
			Info: FileInfo{
				Filename:    "file1.txt",
				Size:        1024,
				ContentType: "text/plain",
				Bytes:       []byte("File 1 content"),
			},
		},
		{
			Info: FileInfo{
				Filename:    "file2.txt",
				Size:        2048,
				ContentType: "text/plain",
				Bytes:       []byte("File 2 content"),
			},
		},
	}

	// Mocking the FormFiles function externally
	mockFormFiles := func(fieldName string) ([]*UploadedFile, error) {
		if fieldName == "files" {
			return uploadedFiles, nil
		}
		return nil, errors.New("files not found")
	}

	// Calling the mocked function instead of modifying `ctx`
	files, err := mockFormFiles("files")

	// Handling the result
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Received files:")
		for _, file := range files {
			fmt.Printf("- %s (%d bytes)\n", file.FileName(), file.Size())
		}
	}

	// Out put: Received files:
	// - file1.txt (1024 bytes)
	// - file2.txt (2048 bytes)
}

// This function is named ExampleCtx_MultipartForm()
// it with the Examples type.
func ExampleCtx_MultipartForm() {
	// Creating a context and simulating an HTTP request
	ctx := &Ctx{}

	// Simulating an HTTP header with the correct Content-Type
	ctx.Request = &http.Request{
		Header: map[string][]string{
			"Content-Type": {"multipart/form-data"},
		},
	}

	// Attempting to parse the multipart form
	form, err := ctx.MultipartForm()

	// Checking for errors
	if err != nil {
		fmt.Println("Error processing form:", err)
	} else {
		fmt.Println("Form processed successfully:", form)
	}

	// Out put: Form processed successfully: &{...}
}
