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

// QtestReturn represents the response and additional functionality for validation.
type QtestReturn interface {
	Body() []byte
	BodyStr() string
	StatusCode() int
	Response() *http.Response
	AssertStatus(expected int) error
	AssertHeader(key, expectedValue string) error
	AssertBodyContains(expected any) error
}

type QTestPlus struct {
	body       []byte
	bodyStr    string
	statusCode int
	response   *http.Response
}

// QuickTestOptions holds all parameters for the enhanced Qtest function.
type QuickTestOptions struct {
	Method      string
	URI         string
	Headers     map[string]string
	QueryParams map[string]string
	Body        []byte
	Cookies     []*http.Cookie
	LogDetails  bool // Enables request/response logging
}

// Qtest performs HTTP tests with query params, cookies, and validation options.
// The result will Qtest(opts QuickTestOptions) (QtestReturn, error)
func (q Quick) Qtest(opts QuickTestOptions) (QtestReturn, error) {
	// Build query string
	uriWithParams, err := attachQueryParams(opts.URI, opts.QueryParams)
	if err != nil {
		return nil, err
	}

	// Create HTTP request
	reqBody := bytes.NewBuffer(opts.Body)
	req, err := http.NewRequest(opts.Method, uriWithParams, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add headers
	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	// Add cookies
	for _, cookie := range opts.Cookies {
		req.AddCookie(cookie)
	}

	// Simulate HTTP request
	rec := httptest.NewRecorder()
	q.ServeHTTP(rec, req) // Calls the handler

	// Capture response
	resp := rec.Result()

	// Read response body safely
	respBody, err := readResponseBodyV2(resp)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Log details if enabled
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

// attachQueryParams builds URL with query parameters.
// The result will attachQueryParams(uri string, params map[string]string) (string, error)
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

// logRequestResponseDetails logs detailed request and response information.
// The result will logRequestResponseDetails(opts QuickTestOptions, resp *http.Response, body []byte)
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

// readResponseBodyV2 reads response body safely and resets it for reuse.
// The result will readResponseBodyV2(resp *http.Response) ([]byte, error)
func readResponseBodyV2(resp *http.Response) ([]byte, error) {
	if resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Reset response body so it can be read again
	resp.Body = io.NopCloser(bytes.NewReader(body))
	return body, nil
}

// Implement QtestReturn interface methods
///

// The result will Body() []byte
func (qt *QTestPlus) Body() []byte {
	return qt.body
}

// The result will BodyStr() string
func (qt *QTestPlus) BodyStr() string {
	return qt.bodyStr
}

// The result will StatusCode() int
func (qt *QTestPlus) StatusCode() int {
	return qt.statusCode
}

// The result will Response() *http.Response
func (qt *QTestPlus) Response() *http.Response {
	return qt.response
}

// AssertStatus checks if the response status matches the expected.
// The result will AssertStatus(expected int) error
func (qt *QTestPlus) AssertStatus(expected int) error {
	if qt.statusCode != expected {
		return fmt.Errorf("expected status %d but got %d", expected, qt.statusCode)
	}
	return nil
}

// AssertHeader checks if a header has the expected value.
// The result will AssertHeader(key, expectedValue string) error
func (qt *QTestPlus) AssertHeader(key, expectedValue string) error {
	value := qt.response.Header.Get(key)
	if value != expectedValue {
		return fmt.Errorf("expected header '%s' to be '%s' but got '%s'", key, expectedValue, value)
	}
	return nil
}

// AssertBodyContains checks if the response body contains a specific substring.
// The result will AssertBodyContains(expected any) error
func (qt *QTestPlus) AssertBodyContains(expected any) error {
	var expectedStr string

	// Convert expected to string (if it's not already a string)
	switch v := expected.(type) {
	case string:
		expectedStr = v
	default:
		// Convert to JSON string
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to convert expected value to JSON: %w", err)
		}
		expectedStr = string(jsonBytes)
	}

	// Perform assertion
	if !strings.Contains(qt.bodyStr, expectedStr) {
		return fmt.Errorf("expected body to contain '%s' but got '%s'", expectedStr, qt.bodyStr)
	}
	return nil
}
