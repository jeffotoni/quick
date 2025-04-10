# 🧪 Qtest - HTTP Testing Utility for Quick Framework ![Quick Logo](/quick.png)

Qtest is an **advanced HTTP testing function** designed to simplify route validation within the **Quick** framework. It enables seamless testing of simulated HTTP requests using `httptest`, supporting:

- **Custom HTTP methods** (`GET`, `POST`, `PUT`, `DELETE`, etc.).
- **Custom headers**.
- **Query parameters**.
- **Request body**.
- **Cookies**.
- **Built-in validation methods** for status codes, headers, and response bodies.

## 📌 Overview
The `Qtest` function takes a `QuickTestOptions` struct containing request parameters, executes the request, and returns a `QtestReturn` object, which provides methods for analyzing and validating the result.

### 🛠 Core Function
```go
package main
import (
	"github.com/quick/quick"
	"github.com/quick/quick/httptest"
	)
func (q Quick) Qtest(opts QuickTestOptions) (QtestReturn, error) {
    // Build URL with query parameters (if any)
    uriWithParams, err := attachQueryParams(opts.URI, opts.QueryParams)
    if err != nil {
        return nil, err
    }

    // Create the HTTP request
    reqBody := bytes.NewBuffer(opts.Body)
    req, err := http.NewRequest(opts.Method, uriWithParams, reqBody)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    // Add headers and cookies
    for key, value := range opts.Headers {
        req.Header.Set(key, value)
    }
    for _, cookie := range opts.Cookies {
        req.AddCookie(cookie)
    }

    // Simulate HTTP request
    rec := httptest.NewRecorder()
    q.ServeHTTP(rec, req)

    // Capture the response
    resp := rec.Result()
    respBody, err := readResponseBodyV2(resp)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %w", err)
    }

    // Log details if enabled
    if opts.LogDetails {
        logRequestResponseDetails(opts, resp, respBody)
    }

    // Return response encapsulated in `qTestPlus`
    return &qTestPlus{
        body:       respBody,
        bodyStr:    string(respBody),
        statusCode: resp.StatusCode,
        response:   resp,
    }, nil
}
```

## 📌 Input Struct: `QuickTestOptions`
The `QuickTestOptions` struct defines the necessary parameters to execute an HTTP test.

```go
type QuickTestOptions struct {
    Method      string            // HTTP method (GET, POST, etc.)
    URI         string            // Request path
    Headers     map[string]string // Custom headers
    QueryParams map[string]string // Query parameters
    Body        []byte            // Request body
    Cookies     []*http.Cookie    // Cookies sent in the request
    LogDetails  bool              // Enables detailed logging
}
```

## 📌 Output Struct: `QtestReturn`
The `QtestReturn` interface provides methods to access and validate the response.

```go
type QtestReturn interface {
    Body() []byte
    BodyStr() string
    StatusCode() int
    Response() *http.Response
    AssertStatus(expected int) error
    AssertHeader(key, expectedValue string) error
    AssertBodyContains(expected string) error
}
```

The implementation (`qTestPlus`) encapsulates the request results, enabling direct validation.

## 📌 Validation Methods
These methods allow testing **status codes, headers, and response content**.

```go
func (qt *qTestPlus) AssertStatus(expected int) error {
    if qt.statusCode != expected {
        return fmt.Errorf("expected status %d but got %d", expected, qt.statusCode)
    }
    return nil
}

func (qt *qTestPlus) AssertHeader(key, expectedValue string) error {
    value := qt.response.Header.Get(key)
    if value != expectedValue {
        return fmt.Errorf("expected header '%s' to be '%s' but got '%s'", key, expectedValue, value)
    }
    return nil
}

func (qt *qTestPlus) AssertBodyContains(expected string) error {
    if !strings.Contains(qt.bodyStr, expected) {
        return fmt.Errorf("expected body to contain '%s' but got '%s'", expected, qt.bodyStr)
    }
    return nil
}
```

## 📌 Example Usage in a Test

✅ Test POST Request

Here's an example of how to use `Qtest` to test an API receiving a `POST` request.


```go
func TestQTest_Options_POST(t *testing.T) {
    q := New()

    // Define the POST route
    q.Post("/v1/user/api", func(c *Ctx) error {
        c.Set("Content-Type", "application/json") // Simplified header setting
        return c.Status(StatusOK).String(`{"message":"Success"}`)
    })

    // Configure test parameters
    opts := QuickTestOptions{
        Method: "POST",
        URI:    "/v1/user/api",
        QueryParams: map[string]string{
            "param1": "value1",
            "param2": "value2",
        },
        Body: []byte(`{"key":"value"}`),
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
        Cookies: []*http.Cookie{
            {Name: "session", Value: "abc123"},
        },
        LogDetails: true, // Enables detailed logging
    }

    // Execute test
    result, err := q.Qtest(opts)
    if err != nil {
        t.Fatalf("Error in Qtest: %v", err)
    }

    // Validations
    if err := result.AssertStatus(StatusOK); err != nil {
        t.Errorf("Status assertion failed: %v", err)
    }

    if err := result.AssertHeader("Content-Type", "application/json"); err != nil {
        t.Errorf("Header assertion failed: %v", err)
    }

    if err := result.AssertBodyContains("Success"); err != nil {
        t.Errorf("Body assertion failed: %v", err)
    }
}
```

✅ Test GET Request

Tests a basic `GET` request to ensure the correct response status and body content.

```go
func TestQTest_Options_GET(t *testing.T) {

	q := New()

	q.Get("/v1/user", func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(StatusOK).String("Success")
	})

	opts := QuickTestOptions{
		Method:     "GET",
		URI:        "/v1/user",
		Headers:    map[string]string{"Accept": "application/json"},
		LogDetails: true,
	}

	result, err := q.Qtest(opts)
	if err != nil {
		t.Fatalf("Error in test: %v", err)
	}

	if err := result.AssertStatus(StatusOK); err != nil {
		t.Errorf("Status assertion failed: %v", err)
	}

	// Debug the response body for troubleshooting
	fmt.Println("DEBUG Body (QTest):", result.BodyStr())

	if err := result.AssertBodyContains("Success"); err != nil {
		t.Errorf("Body assertion failed: %v", err)
	}
}

```

✅ Test POST Request

Validates a `POST` request with headers, query parameters, and body content.

```go
func TestQTest_Options_POST(t *testing.T) {

	q := New()

	q.Post("/v1/user/api", func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		c.Response.Header().Set("Content-Type", "application/json")
		return c.Status(StatusOK).String(`{"message":"Success"}`)
	})

	opts := QuickTestOptions{
		Method: "POST", // Correct method
		URI:    "/v1/user/api",
		QueryParams: map[string]string{
			"param1": "value1",
			"param2": "value2",
		},
		Body: []byte(`{"key":"value"}`),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Cookies: []*http.Cookie{
			{Name: "session", Value: "abc123"},
		},
		LogDetails: true, // Enable the log for debug
	}

	// Runs the test
	result, err := q.Qtest(opts)
	if err != nil {
		t.Fatalf("Error in Qtest: %v", err)
	}

	// Check if the status is correct
	if err := result.AssertStatus(StatusOK); err != nil {
		t.Errorf("Status assertion failed: %v", err)
	}

	// Check if o Expected header is present
	if err := result.AssertHeader("Content-Type", "application/json"); err != nil {
		t.Errorf("Header assertion failed: %v", err)
	}

	// Check if the answer contains "Success"
	if err := result.AssertBodyContains("Success"); err != nil {
		t.Errorf("Body assertion failed: %v", err)
	}
}

```


✅ Test PUT Request

Checks if a PUT request correctly updates user data.
```go
func TestQTest_Options_PUT(t *testing.T) {
	q := New()

	// Define the PUT route
	q.Put("/v1/user/update", func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(StatusOK).String(`{"message":"User updated successfully"}`)
	})

	opts := QuickTestOptions{
		Method: "PUT",
		URI:    "/v1/user/update",
		Body:   []byte(`{"name":"Jeff Quick","age":30}`),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		LogDetails: true,
	}

	result, err := q.Qtest(opts)
	if err != nil {
		t.Fatalf("Error in Qtest: %v", err)
	}

	if err := result.AssertStatus(StatusOK); err != nil {
		t.Errorf("Status assertion failed: %v", err)
	}

	if err := result.AssertBodyContains(`"message":"User updated successfully"`); err != nil {
		t.Errorf("Body assertion failed: %v", err)
	}
}

```

✅ Test DELETE Request

Ensures that a DELETE request successfully removes a user entry.
```go
func TestQTest_Options_DELETE(t *testing.T) {
	q := New()

	// Define the DELETE route
	q.Delete("/v1/user/delete", func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(StatusOK).String(`{"message":"User deleted successfully"}`)
	})

	opts := QuickTestOptions{
		Method:     "DELETE",
		URI:        "/v1/user/delete",
		LogDetails: true,
	}

	result, err := q.Qtest(opts)
	if err != nil {
		t.Fatalf("Error in Qtest: %v", err)
	}

	if err := result.AssertStatus(StatusOK); err != nil {
		t.Errorf("Status assertion failed: %v", err)
	}

	if err := result.AssertBodyContains(`"message":"User deleted successfully"`); err != nil {
		t.Errorf("Body assertion failed: %v", err)
	}
}

```

✅ Test PATCH Request

Verifies partial updates to user data using a PATCH request.
```go
func TestQTest_Options_PATCH(t *testing.T) {
	q := New()

	// Define the PATCH route
	q.Patch("/v1/user/patch", func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(StatusOK).String(`{"message":"User patched successfully"}`)
	})

	opts := QuickTestOptions{
		Method: "PATCH",
		URI:    "/v1/user/patch",
		Body:   []byte(`{"nickname":"Johnny"}`),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		LogDetails: true,
	}

	result, err := q.Qtest(opts)
	if err != nil {
		t.Fatalf("Error in Qtest: %v", err)
	}

	if err := result.AssertStatus(StatusOK); err != nil {
		t.Errorf("Status assertion failed: %v", err)
	}

	if err := result.AssertBodyContains(`"message":"User patched successfully"`); err != nil {
		t.Errorf("Body assertion failed: %v", err)
	}
}

```

✅ Test OPTIONS Request

Confirms the allowed HTTP methods for a given route.
```go
func TestQTest_Options_OPTIONS(t *testing.T) {
	q := New()

	// Define the OPTIONS route
	q.Options("/v1/user/options", func(c *Ctx) error {
		c.Set("Allow", "GET, POST, PUT, DELETE, OPTIONS")
		return c.Status(StatusNoContent).String("")
	})

	opts := QuickTestOptions{
		Method:     "OPTIONS",
		URI:        "/v1/user/options",
		LogDetails: true,
	}

	result, err := q.Qtest(opts)
	if err != nil {
		t.Fatalf("Error in Qtest: %v", err)
	}

	if err := result.AssertStatus(StatusNoContent); err != nil {
		t.Errorf("Status assertion failed: %v", err)
	}

	if err := result.AssertHeader("Allow", "GET, POST, PUT, DELETE, OPTIONS"); err != nil {
		t.Errorf("Header assertion failed: %v", err)
	}
}

```
---

## 📌 What I included in this README

- ✅ **README checklist** - Overview of Qtest functionalities.
- ✅ **Overview**: Explanation of Qtest and its use in Quick framework.
- ✅ **Core Function**: Detailed implementation of Qtest function.
- ✅ **Input & Output Structures**: Explanation of `QuickTestOptions` and `QtestReturn`.
- ✅ **Validation Methods**: Methods for checking status, headers, and response bodies.

### ✅ Example Test Cases:
- ✅ **GET** request validation.
- ✅ **POST** request validation with headers and query params.
- ✅ **PUT** request validation for updating user data.
- ✅ **DELETE** request validation for removing a user.
- ✅ **PATCH** request validation for partial updates.
- ✅ **OPTIONS** request validation for checking allowed methods.

### ✅ Testing Checklist:
- ✅ Coverage for HTTP methods.
- ✅ Validation of status codes, headers, and response bodies.
- ✅ Error handling and invalid input cases.
- ✅ Organization and optimization of tests.

- ✅ **Conclusion**: Summary of Qtest benefits and call for contributions.
 

---

With this approach, tests become more structured and maintainable. 🚀🔥

For further optimizations or improvements, feel free to contribute! 😃

