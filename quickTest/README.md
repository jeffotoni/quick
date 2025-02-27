# Qtest - HTTP Testing Utility for Quick Framework

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

## 🔥 Conclusion
`Qtest` is a powerful tool for testing in the **Quick** framework, simplifying API validation. It allows:
✅ Testing any **HTTP method**.  
✅ Easily adding **headers, cookies, and query parameters**.  
✅ **Validating status, headers, and response bodies** intuitively.  
✅ **Enabling detailed logs** for debugging.  

With this approach, tests become more structured and maintainable. 🚀🔥

For further optimizations or improvements, feel free to contribute! 😃

