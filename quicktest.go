package quick

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
)

const logDelimiter = "====================="

// QuickTestReturn defines the interface for handling HTTP test responses.
type QuickTestReturn interface {
	Body() []byte             // Returns the raw response body as a byte slice.
	BodyStr() string          // Returns the response body as a string.
	StatusCode() int          // Returns the HTTP status code.
	Response() *http.Response // Returns the full HTTP response object.
}

type (
	qTest struct {
		body       []byte
		bodyStr    string
		statusCode int
		response   *http.Response
	}

	// QuickMockTestServer defines a mock server configuration for testing.
	QuickMockTestServer struct {
		Client  *http.Client      // HTTP client to interact with the mock server.
		Port    int               // Port on which the mock server runs.
		URI     string            // The request URI for the test.
		Method  string            // The HTTP method (GET, POST, etc.).
		Headers map[string]string // Headers to be included in the request.
		Body    []byte            // Request body content.
	}
)

// QuickTest is a helper function to perform HTTP tests quickly.
//
// It creates an HTTP request using the specified method, URI, and optional body,
// then executes the request within Quick's internal router and returns the response.
//
// Parameters:
//   - method (string): HTTP method (e.g., "GET", "POST").
//   - URI (string): The request path (e.g., "/test/:param").
//   - headers (map[string]string): Optional headers to include in the request.
//   - body (optional []byte): Request body, used for methods like POST.
//
// Returns:
//   - QuickTestReturn: The response object containing the body, status code, and full HTTP response.
//   - error: Any error encountered during the request creation or execution.
//
// Example Usage:
//
//	resp, err := q.QuickTest("GET", "/api/user", nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Status Code:", resp.StatusCode())
//	fmt.Println("Response Body:", resp.BodyStr())
func (q Quick) QuickTest(method, URI string, headers map[string]string, body ...[]byte) (QuickTestReturn, error) {
	requestBody := []byte{}
	if len(body) > 0 {
		requestBody = body[0]
	}

	// Log request details for debugging purposes.
	logRequestDetails(method, URI, len(requestBody))

	// Create the HTTP request.
	req, err := createHTTPRequest(method, URI, headers, requestBody)
	if err != nil {
		return nil, err
	}

	// Create a new HTTP response recorder to capture the output.
	rec := httptest.NewRecorder()
	q.ServeHTTP(rec, req)

	// Retrieve the recorded response.
	resp := rec.Result()
	responseBody, err := readResponseBody(resp.Body)
	if err != nil {
		return nil, err
	}

	return &qTest{
		body:       responseBody,
		bodyStr:    string(responseBody),
		statusCode: resp.StatusCode,
		response:   resp,
	}, nil
}

// createHTTPRequest constructs an HTTP request with the specified method, URI, headers, and body.
//
// Parameters:
//   - method (string): HTTP method (e.g., "GET", "POST").
//   - URI (string): The request path (e.g., "/api/test").
//   - headers (map[string]string): Headers to include in the request.
//   - body ([]byte): The request body content.
//
// Returns:
//   - *http.Request: The constructed HTTP request.
//   - error: Any error encountered while creating the request.
func createHTTPRequest(method, URI string, headers map[string]string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, URI, io.NopCloser(bytes.NewBuffer(body)))
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return req, nil
}

// readResponseBody safely reads and returns the response body as a byte slice.
//
// Parameters:
//   - body (io.ReadCloser): The response body stream.
//
// Returns:
//   - []byte: The content of the response body.
//   - error: Any error encountered while reading the body.
func readResponseBody(body io.ReadCloser) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	defer body.Close()
	return io.ReadAll(body)
}

// logRequestDetails logs the details of an HTTP request.
//
// This function prints the HTTP method, URI, and body length for debugging purposes.
//
// Parameters:
//   - method (string): The HTTP method used in the request.
//   - URI (string): The request path.
//   - bodyLen (int): The length of the request body.
func logRequestDetails(method, URI string, bodyLen int) {
	println(logDelimiter)
	println("Method:", method, "| URI:", URI, "| Body Length:", bodyLen)
	println(logDelimiter)
}

// Body returns the raw response body as a byte slice.
func (qt *qTest) Body() []byte {
	return qt.body
}

// BodyStr returns the response body as a string.
func (qt *qTest) BodyStr() string {
	return qt.bodyStr
}

// StatusCode returns the HTTP status code of the response.
func (qt *qTest) StatusCode() int {
	return qt.statusCode
}

// Response returns the full HTTP response object.
func (qt *qTest) Response() *http.Response {
	return qt.response
}
