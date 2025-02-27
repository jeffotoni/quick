package quick

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
)

const logDelimiter = "====================="

type QuickTestReturn interface {
	Body() []byte
	BodyStr() string
	StatusCode() int
	Response() *http.Response
}

type (
	qTest struct {
		body       []byte
		bodyStr    string
		statusCode int
		response   *http.Response
	}

	QuickMockTestServer struct {
		Client  *http.Client
		Port    int
		URI     string
		Method  string
		Headers map[string]string
		Body    []byte
	}
)

// QuickTest: Helper function to make HTTP tests quickly.
// Required Params: method (e.g., GET, POST), URI (path only, e.g., /test/:param)
// Optional Param: body (optional; use only when necessary)
func (q Quick) QuickTest(method, URI string, headers map[string]string, body ...[]byte) (QuickTestReturn, error) {
	requestBody := []byte{}
	if len(body) > 0 {
		requestBody = body[0]
	}

	logRequestDetails(method, URI, len(requestBody))

	req, err := createHTTPRequest(method, URI, headers, requestBody)
	if err != nil {
		return nil, err
	}

	rec := httptest.NewRecorder()
	q.ServeHTTP(rec, req)

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

// createHTTPRequest: Encapsulates the creation of an HTTP request.
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
func readResponseBody(body io.ReadCloser) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	defer body.Close()
	return io.ReadAll(body)
}

// logRequestDetails logs the details of the HTTP request.
func logRequestDetails(method, URI string, bodyLen int) {
	println(logDelimiter)
	println("Method:", method, "| URI:", URI, "| Body Length:", bodyLen)
	println(logDelimiter)
}

func (qt *qTest) Body() []byte {
	return qt.body
}

func (qt *qTest) BodyStr() string {
	return qt.bodyStr
}

func (qt *qTest) StatusCode() int {
	return qt.statusCode
}

func (qt *qTest) Response() *http.Response {
	return qt.response
}
