// Package quick provides a fast and flexible web framework with built-in
// HTTP testing utilities. This file contains example functions demonstrating
// how to use Quick's testing utilities.
package quick

import (
	"fmt"
	"net/http"
)

// ExampleQuick_Qtest demonstrates how to use Qtest to test a simple GET route.
//
// The "/hello" route returns a response with status 200 and body "Hello, Quick!".
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

// ExampleQTestPlus_AssertStatus demonstrates how to verify the response status code.
//
// The "/notfound" route returns a 404 status, which is validated in the assertion.
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

// ExampleQTestPlus_AssertHeader demonstrates how to validate an HTTP response header.
//
// The "/header" route sets an "X-Custom-Header", which is then validated.
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

// ExampleQTestPlus_AssertBodyContains demonstrates how to validate the response body.
//
// The "/json" route returns a JSON response, and the test verifies the presence of the "message" key.
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

// ExampleQTestPlus_Body demonstrates how to retrieve the response body as a byte slice.
//
// The simulated response object contains "Hello, Quick!" as its body content.
func ExampleQTestPlus_Body() {
	// Simulating a response object with a predefined body
	res := &QTestPlus{
		body: []byte("Hello, Quick!"),
	}

	// Retrieving and printing the body content
	fmt.Println(string(res.Body()))

	// Out put: Hello, Quick!
}

// ExampleQTestPlus_BodyStr demonstrates how to retrieve the response body as a string.
//
// The simulated response object contains "Hello, Quick!" as its body content.
func ExampleQTestPlus_BodyStr() {
	// Simulating a response object with a predefined body
	res := &QTestPlus{
		bodyStr: "Hello, Quick!",
	}

	// Retrieving and printing the body content as a string
	fmt.Println(res.BodyStr())

	// Out put: Hello, Quick!
}

// ExampleQTestPlus_StatusCode demonstrates how to retrieve the response status code.
//
// The simulated response object contains a status code of 200.
func ExampleQTestPlus_StatusCode() {
	res := &QTestPlus{
		statusCode: 200,
	}

	// Retrieving and printing the response status code
	fmt.Println("Status Code:", res.StatusCode())

	// Out put: Status Code: 200
}

// ExampleQTestPlus_Response demonstrates how to retrieve the complete HTTP response object.
//
// The simulated HTTP response object contains a status of "200 OK".
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

// ExampleQTestPlus_AssertNoHeader demonstrates how to check that a header is not present in the response.
//
// The simulated HTTP response does not include the "X-Powered-By" header.
func ExampleQTestPlus_AssertNoHeader() {
	q := New()

	q.Get("/no-header", func(c *Ctx) error {
		return c.String("No custom header here.")
	})

	res, err := q.Qtest(QuickTestOptions{
		Method: "GET",
		URI:    "/no-header",
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = res.AssertNoHeader("X-Powered-By")
	if err != nil {
		fmt.Println("Assertion failed:", err)
	} else {
		fmt.Println("Header is not present as expected")
	}

	// Out put: Header is not present as expected
}

// ExampleQTestPlus_AssertString demonstrates how to compare the response body with an expected string.
//
// The simulated HTTP response body contains "pong".
func ExampleQTestPlus_AssertString() {

	q := New()

	q.Get("/ping", func(c *Ctx) error {
		return c.String("pong")
	})

	res, _ := q.Qtest(QuickTestOptions{
		Method: "GET",
		URI:    "/ping",
	})

	err := res.AssertString("pong")
	if err != nil {
		fmt.Println("Assertion failed:", err)
	} else {
		fmt.Println("Body matches expected string")
	}

	// Out put: Body matches expected string
}

// ExampleQTestPlus_AssertHeaderHasPrefix demonstrates how to verify if a header starts with a prefix.
//
// The simulated HTTP response includes "X-Version" header with value "v1.2.3".
func ExampleQTestPlus_AssertHeaderHasPrefix() {
	q := New()

	q.Get("/prefix", func(c *Ctx) error {
		c.Set("X-Version", "v1.2.3")
		return c.String("OK")
	})

	res, _ := q.Qtest(QuickTestOptions{
		Method: "GET",
		URI:    "/prefix",
	})

	err := res.AssertHeaderHasPrefix("X-Version", "v1")
	if err != nil {
		fmt.Println("Assertion failed:", err)
	} else {
		fmt.Println("Header has expected prefix")
	}

	// Out put: Header has expected prefix
}

// ExampleQTestPlus_AssertHeaderHasValueInSet demonstrates how to verify if a header value matches one of the allowed values.
//
// The simulated HTTP response includes "X-Env" header with value "staging".
func ExampleQTestPlus_AssertHeaderHasValueInSet() {

	q := New()

	q.Get("/variant", func(c *Ctx) error {
		c.Set("X-Env", "staging")
		return c.String("OK")
	})

	res, _ := q.Qtest(QuickTestOptions{
		Method: "GET",
		URI:    "/variant",
	})

	err := res.AssertHeaderHasValueInSet("X-Env", []string{"dev", "staging", "prod"})
	if err != nil {
		fmt.Println("Assertion failed:", err)
	} else {
		fmt.Println("Header value is in allowed set")
	}

	// Out put: Header value is in allowed set
}

// ExampleQTestPlus_AssertHeaderContains demonstrates how to verify if a header contains a substring.
//
// The simulated HTTP response includes "X-Custom" header with value "PoweredByQuick".
func ExampleQTestPlus_AssertHeaderContains() {
	q := New()

	q.Get("/header", func(c *Ctx) error {
		c.Set("X-Custom", "PoweredByQuick")
		return c.String("OK")
	})

	res, _ := q.Qtest(QuickTestOptions{
		Method: "GET",
		URI:    "/header",
	})

	err := res.AssertHeaderContains("X-Custom", "Quick")
	if err != nil {
		fmt.Println("Assertion failed:", err)
	} else {
		fmt.Println("Header contains expected substring")
	}

	// Out put: Header contains expected substring
}
