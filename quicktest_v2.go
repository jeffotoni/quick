// Qtest is an advanced HTTP testing function designed to
// facilitate route validation in the Quick framework.
//
// It allows you to test simulated HTTP requests using httptest, supporting:
//   - Métodos HTTP personalizados (GET, POST, PUT, DELETE, etc.)
//   - Cabeçalhos (Headers) personalizados.
//   - Parâmetros de Query (QueryParams).
//   - Corpo da Requisição (Body).
//   - Cookies.
//   - Validações embutidas para status, headers e corpo da resposta.
//
// The Qtest function receives a QuickTestOptions structure containing the request
// parameters, executes the call and returns a QtestReturn object, which provides methods
// for analyzing and validating the result.
package quick

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

// QtestReturn defines an interface for validating HTTP test responses.
//
// This interface provides methods to retrieve response details such as body, status code,
// headers, and to perform assertions for testing HTTP responses.
//
// Example Usage:
//
//	resp, err := q.Qtest(quick.QuickTestOptions{
//	    Method: quick.MethodGet,
//	    URI:    "/test",
//	})
//
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	if err := resp.AssertStatus(200); err != nil {
//	    log.Fatal(err)
//	}
type QtestReturn interface {
	Body() []byte
	BodyStr() string
	StatusCode() int
	Response() *http.Response
	AssertStatus(expected int) error
	AssertHeader(key, expectedValue string) error
	AssertBodyContains(expected any) error
}

// QTestPlus implements QtestReturn, encapsulating HTTP response details for testing.
type QTestPlus struct {
	body       []byte
	bodyStr    string
	statusCode int
	response   *http.Response
}

// QuickTestOptions defines configuration options for executing an HTTP test request.
//
// Example Usage:
//
//	opts := quick.QuickTestOptions{
//	    Method:      quick.MethodPost,
//	    URI:         "/submit",
//	    Headers:     map[string]string{"Content-Type": "application/json"},
//	    QueryParams: map[string]string{"id": "123"},
//	    Body:        []byte(`{"name":"John Doe"}`),
//	}
//
//	resp, err := q.Qtest(opts)
type QuickTestOptions struct {
	Method      string            // HTTP method (e.g., "GET", "POST")
	URI         string            // Target URI for the request
	Headers     map[string]string // Request headers
	QueryParams map[string]string // Query parameters to append to the URI
	Body        []byte            // Request body payload
	Cookies     []*http.Cookie    // Cookies to include in the request
	LogDetails  bool              // Enable logging of request and response details
}

// Qtest performs an HTTP request using QuickTestOptions and returns a response handler.
//
// This method executes a test HTTP request within the Quick framework, allowing validation
// of the response status, headers, and body.
//
// Example Usage:
//
//	resp, err := q.Qtest(quick.QuickTestOptions{
//	    Method: quick.MethodGet,
//	    URI:    "/test",
//	})
//
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	if err := resp.AssertStatus(200); err != nil {
//	    log.Fatal(err)
//	}
//
// Returns:
//   - QtestReturn: An interface for response validation
//   - error: Error encountered during request execution, if any.
func (q Quick) Qtest(opts QuickTestOptions) (QtestReturn, error) {
	uriWithParams, err := attachQueryParams(opts.URI, opts.QueryParams)
	if err != nil {
		return nil, err
	}

	reqBody := bytes.NewBuffer(opts.Body)
	req, err := http.NewRequest(opts.Method, uriWithParams, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	// Add cookies
	for _, cookie := range opts.Cookies {
		req.AddCookie(cookie)
	}

	// Simulate HTTP request execution
	rec := httptest.NewRecorder()
	q.ServeHTTP(rec, req)

	// Capture response
	resp := rec.Result()
	respBody, err := readResponseBodyV2(resp)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if opts.LogDetails {
		logRequestResponseDetails(opts, resp, respBody)
	}

	return &QTestPlus{
		body:       respBody,
		bodyStr:    string(respBody),
		statusCode: resp.StatusCode,
		response:   resp,
	}, nil
}

// attachQueryParams appends query parameters to a URI.
//
// Example Usage:
//
//	newURI, err := attachQueryParams("/search", map[string]string{"q": "golang"})
//
// Returns:
//   - string: The URI with query parameters appended
//   - error: Returns an error if the URI parsing fails.
func attachQueryParams(uri string, params map[string]string) (string, error) {
	if len(params) == 0 {
		return uri, nil
	}
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	query := u.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	u.RawQuery = query.Encode()
	return u.String(), nil
}

// logRequestResponseDetails logs request and response information for debugging.
//
// If LogDetails is enabled in QuickTestOptions, this function prints the details
// of the HTTP request and response to the console.
//
// Example Usage:
//
//	logRequestResponseDetails(opts, response, responseBody)
func logRequestResponseDetails(opts QuickTestOptions, resp *http.Response, body []byte) {
	fmt.Println("========================================")
	fmt.Printf("Request: %s %s\n", opts.Method, opts.URI)
	fmt.Printf("Request Body: %s\n", string(opts.Body))
	fmt.Println("--- Response ---")
	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Headers: %+v\n", resp.Header)
	fmt.Printf("Body: %s\n", string(body))
	fmt.Println("========================================")
}

// readResponseBodyV2 safely reads and resets the response body for reuse.
//
// Example Usage:
//
//	body, err := readResponseBodyV2(response)
//
// Returns:
//   - []byte: The response body as a byte slice
//   - error: Returns an error if body reading fails.
func readResponseBodyV2(resp *http.Response) ([]byte, error) {
	if resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body = io.NopCloser(bytes.NewReader(body)) // Reset response body for reuse
	return body, nil
}

// Body returns the response body as a byte slice.
//
// Returns:
//   - []byte: The response body.
func (qt *QTestPlus) Body() []byte {
	return qt.body
}

// BodyStr returns the response body as a string.
//
// Returns:
//   - string: The response body as a string.
func (qt *QTestPlus) BodyStr() string {
	return qt.bodyStr
}

// StatusCode retrieves the HTTP status code of the response.
//
// Returns:
//   - int: The HTTP status code.
func (qt *QTestPlus) StatusCode() int {
	return qt.statusCode
}

// Response returns the raw *http.Response for advanced validation.
//
// Returns:
//   - *http.Response: The full HTTP response object.
func (qt *QTestPlus) Response() *http.Response {
	return qt.response
}

// AssertStatus verifies if the response status matches the expected status.
//
// Example Usage:
//
//	err := resp.AssertStatus(200)
//	if err != nil {
//	    t.Errorf("Unexpected status: %v", err)
//	}
//
// Returns:
//   - error: Returns an error if the expected status does not match the actual status.
func (qt *QTestPlus) AssertStatus(expected int) error {
	if qt.statusCode != expected {
		return fmt.Errorf("expected status %d but got %d", expected, qt.statusCode)
	}
	return nil
}

// AssertHeader verifies if the specified header has the expected value.
//
// Example Usage:
//
//	err := resp.AssertHeader("Content-Type", "application/json")
//
// Returns:
//   - error: Returns an error if the header does not match the expected value.
func (qt *QTestPlus) AssertHeader(key, expectedValue string) error {
	value := qt.response.Header.Get(key)
	if value != expectedValue {
		return fmt.Errorf("expected header '%s' to be '%s' but got '%s'", key, expectedValue, value)
	}
	return nil
}

// AssertBodyContains checks if the response body contains the expected content.
//
// Example Usage:
//
//	err := resp.AssertBodyContains("Success")
//
// Returns:
//   - error: Returns an error if the expected content is not found in the body.
func (qt *QTestPlus) AssertBodyContains(expected any) error {
	var expectedStr string

	switch v := expected.(type) {
	case string:
		expectedStr = v
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to convert expected value to JSON: %w", err)
		}
		expectedStr = string(jsonBytes)
	}

	if !strings.Contains(qt.bodyStr, expectedStr) {
		return fmt.Errorf("expected body to contain '%s' but got '%s'", expectedStr, qt.bodyStr)
	}
	return nil
}
