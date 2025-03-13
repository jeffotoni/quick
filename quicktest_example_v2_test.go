// Package quick fornece uma estrutura web rápida e flexível com built-in
// recursos de teste HTTP. Este arquivo contém várias funções de exemplo
// demonstrando o uso dos utilitários de teste do Quick.
package quick

import (
	"fmt"
	"net/http"
)

// This function is named ExampleQuick_Qtest()
// it with the Examples type.
func ExampleQuick_Qtest() {
	// Creating a Quick instance
	q := New()

	// Defining a simple GET route
	q.Get("/hello", func(c *Ctx) error {
		return c.Status(200).String("Hello, Quick!")
	})

	// Defining the request parameters
	opts := QuickTestOptions{
		Method: "GET",
		URI:    "/hello",
	}

	// Performing the HTTP test using the Quick instance
	res, err := q.Qtest(opts)

	// Handling errors
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Printing response details
	fmt.Println("Status Code:", res.StatusCode())
	fmt.Println("Body:", res.BodyStr())

	// Out put: Status Code: 200
	// Body: Hello, Quick!
}

// This function is named Example_qTestPlus_AssertStatus()
// it with the Examples type.
func ExampleQTestPlus_AssertStatus() {
	// Creating a Quick instance
	q := New()

	// Defining a route that returns 404
	q.Get("/notfound", func(c *Ctx) error {
		return c.Status(404).String("Not Found")
	})

	// Performing the HTTP test
	opts := QuickTestOptions{
		Method: "GET",
		URI:    "/notfound",
	}

	res, err := q.Qtest(opts)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Validating the response status
	err = res.AssertStatus(404)
	if err != nil {
		fmt.Println("Assertion failed:", err)
	} else {
		fmt.Println("Status code is correct")
	}

	// Out put: Status code is correct
}

// This function is named Example_qTestPlus_AssertHeader()
// it with the Examples type.
func ExampleQTestPlus_AssertHeader() {
	// Creating a Quick instance
	q := New()

	// Defining a route with a custom header
	q.Get("/header", func(c *Ctx) error {
		c.Set("X-Custom-Header", "QuickFramework")
		return c.Status(200).String("OK")
	})

	// Performing the HTTP test
	opts := QuickTestOptions{
		Method: "GET",
		URI:    "/header",
	}

	res, err := q.Qtest(opts)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Validating the header
	err = res.AssertHeader("X-Custom-Header", "QuickFramework")
	if err != nil {
		fmt.Println("Header assertion failed:", err)
	} else {
		fmt.Println("Header is correct")
	}

	// Out put: Header is correct
}

// This function is named Example_qTestPlus_AssertBodyContains()
// it with the Examples type.
func ExampleQTestPlus_AssertBodyContains() {
	// Creating a Quick instance
	q := New()

	// Defining a route that returns JSON
	q.Get("/json", func(c *Ctx) error {
		data := map[string]string{"message": "Hello, Quick!"}
		return c.JSON(data)
	})

	// Performing the HTTP test
	opts := QuickTestOptions{
		Method: "GET",
		URI:    "/json",
	}

	res, err := q.Qtest(opts)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Validating the response body
	err = res.AssertBodyContains(`"message":"Hello, Quick!"`)
	if err != nil {
		fmt.Println("Body assertion failed:", err)
	} else {
		fmt.Println("Body contains expected content")
	}

	// Out put: Body contains expected content
}

// This function is named ExampleQTestPlus_Body()
// it with the Examples type.
func ExampleQTestPlus_Body() {
	// Simulating a response object with a predefined body
	res := &QTestPlus{
		body: []byte("Hello, Quick!"),
	}

	// Retrieving and printing the body content
	fmt.Println(string(res.Body()))

	// Out put: Hello, Quick!
}

// This function is named ExampleQTestPlus_BodyStr()
// it with the Examples type.
func ExampleQTestPlus_BodyStr() {
	// Simulating a response object with a predefined body
	res := &QTestPlus{
		bodyStr: "Hello, Quick!",
	}

	// Retrieving and printing the body content as a string
	fmt.Println(res.BodyStr())

	// Out put: Hello, Quick!
}

// This function is named ExampleQTestPlus_StatusCode()
// it with the Examples type.
func ExampleQTestPlus_StatusCode() {
	// Simulating a response object with a predefined status code
	res := &QTestPlus{
		statusCode: 200,
	}

	// Retrieving and printing the response status code
	fmt.Println("Status Code:", res.StatusCode())

	// Out put: Status Code: 200
}

// This function is named ExampleQTestPlus_Response()
// it with the Examples type.
func ExampleQTestPlus_Response() {
	// Simulating an HTTP response object
	httpResponse := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
	}

	// Simulating a response object containing the HTTP response
	res := &QTestPlus{
		response: httpResponse,
	}

	// Retrieving and printing the response status
	fmt.Println("Response Status:", res.Response().Status)

	// Out put: Response Status: 200 OK
}
