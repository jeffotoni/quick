package quick

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/gojeffotoni/quick/internal/concat"
)

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

// QuickTest: This Method is a helper function to make tests with quick more quickly
// Required Params: Method (GET, POST, PUT, DELETE...), URI (only the path. Example: /test/:myParam)
// Optional Param: Body (If you don't want to define one, just ignore)
func (q Quick) QuickTest(method, URI string, headers map[string]string, body ...[]byte) (QuickTestReturn, error) {
	var buffBody []byte

	if len(body) > 0 {
		buffBody = body[0]
	}

	req, err := http.NewRequest(method, URI, io.NopCloser(bytes.NewBuffer(buffBody)))

	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	rec := httptest.NewRecorder()
	q.ServeHTTP(rec, req)

	resp := rec.Result()

	if resp.Body != nil {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return &qTest{
			body:       b,
			bodyStr:    string(b),
			statusCode: resp.StatusCode,
			response:   resp,
		}, nil
	}
	return nil, errors.New("return body is empty")
}

func (q Quick) QuickTestListen(qs QuickMockTestServer) (QuickTestReturn, error) {
	port := strconv.Itoa(qs.Port)

	port = concat.String(":", port)
	URI := concat.String("http://0.0.0.0", port, qs.URI)

	req, err := http.NewRequest(qs.Method, URI, io.NopCloser(bytes.NewBuffer(qs.Body)))
	if err != nil {
		return nil, err
	}

	for k, v := range qs.Headers {
		req.Header.Set(k, v)
	}

	go q.Listen(port)

	// This is a wait time to start the server in go routine
	time.Sleep(time.Millisecond * 100)

	if qs.Client == nil {
		qs.Client = http.DefaultClient
	}

	resp, err := qs.Client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &qTest{
		body:       b,
		bodyStr:    string(b),
		statusCode: resp.StatusCode,
		response:   resp,
	}, nil
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
