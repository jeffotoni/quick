package httpclient

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

// cover     ->  go test -v -count=1 -cover -failfast -run ^TestGet$
// coverHTML ->  go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestGet$; go tool cover -html=coverage.out
func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantOut *ClientResponse
		mock    *httpMock
	}{
		{
			name: "success",
			ctx:  context.Background(),
			wantOut: &ClientResponse{
				Body:       []byte(`{"data": "quick is awesome!"}`),
				StatusCode: 200,
				Error:      nil,
			},
			mock: &httpMock{
				response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(`{"data": "quick is awesome!"}`)),
				},
				err: nil,
			},
		},
		{
			name: "error_request",
			ctx:  context.Background(),
			wantOut: &ClientResponse{
				Body:       nil,
				StatusCode: 0,
				Error:      errors.New("request_error"),
			},
			mock: &httpMock{
				response: &http.Response{
					StatusCode: 500,
					Body:       nil,
				},
				err: errors.New("request_error"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(*testing.T) {

			c := Client{
				Ctx:        context.Background(),
				ClientHttp: tc.mock,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}

			resp := c.Get("some/url")

			strBody, wantStrBody := string(resp.Body), string(tc.wantOut.Body)

			if strBody != wantStrBody {
				t.Errorf("want %s and got %s", strBody, wantStrBody)
			}

			if tc.wantOut.Error != nil {
				if resp.Error.Error() != tc.wantOut.Error.Error() {
					t.Errorf("want %v and got %v", tc.wantOut.Error, resp.Error)
				}
			}

			if resp.StatusCode != tc.wantOut.StatusCode {
				t.Errorf("want %d and got %d", tc.wantOut.StatusCode, resp.StatusCode)
			}

		})
	}
}

// cover     ->  go test -v -count=1 -cover -failfast -run ^TestPost$
// coverHTML ->  go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestPost$; go tool cover -html=coverage.out
func TestPost(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		bodyReq string
		wantOut *ClientResponse
		mock    *httpMock
	}{
		{
			name: "success",
			ctx:  context.Background(),
			bodyReq: `{"data": "quick is awesome!"}`,
			wantOut: &ClientResponse{
				Body:       []byte(`{"data": "quick is awesome!"}`),
				StatusCode: 200,
				Error:      nil,
			},
			mock: &httpMock{
				response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(`{"data": "quick is awesome!"}`)),
				},
				err: nil,
			},
		},
		{
			name: "error_request",
			ctx:  context.Background(),
			wantOut: &ClientResponse{
				Body:       nil,
				StatusCode: 0,
				Error:      errors.New("request_error"),
			},
			mock: &httpMock{
				response: &http.Response{
					StatusCode: 500,
					Body:       nil,
				},
				err: errors.New("request_error"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(*testing.T) {

			c := Client{
				Ctx:        context.Background(),
				ClientHttp: tc.mock,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}

			resp := c.Post("some/url", io.NopCloser(strings.NewReader(tc.bodyReq)))

			strBody, wantStrBody := string(resp.Body), string(tc.wantOut.Body)

			if strBody != wantStrBody {
				t.Errorf("want %s and got %s", strBody, wantStrBody)
			}

			if tc.wantOut.Error != nil {
				if resp.Error.Error() != tc.wantOut.Error.Error() {
					t.Errorf("want %v and got %v", tc.wantOut.Error, resp.Error)
				}
			}

			if resp.StatusCode != tc.wantOut.StatusCode {
				t.Errorf("want %d and got %d", tc.wantOut.StatusCode, resp.StatusCode)
			}

		})
	}
}


// cover     ->  go test -v -count=1 -cover -failfast -run ^TestPut$
// coverHTML ->  go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestPut$; go tool cover -html=coverage.out
func TestPut(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		bodyReq string
		wantOut *ClientResponse
		mock    *httpMock
	}{
		{
			name: "success",
			ctx:  context.Background(),
			bodyReq: `{"data": "quick is awesome!"}`,
			wantOut: &ClientResponse{
				Body:       []byte(`{"data": "quick is awesome!"}`),
				StatusCode: 200,
				Error:      nil,
			},
			mock: &httpMock{
				response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(`{"data": "quick is awesome!"}`)),
				},
				err: nil,
			},
		},
		{
			name: "error_request",
			ctx:  context.Background(),
			wantOut: &ClientResponse{
				Body:       nil,
				StatusCode: 0,
				Error:      errors.New("request_error"),
			},
			mock: &httpMock{
				response: &http.Response{
					StatusCode: 500,
					Body:       nil,
				},
				err: errors.New("request_error"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(*testing.T) {

			c := Client{
				Ctx:        context.Background(),
				ClientHttp: tc.mock,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}

			resp := c.Put("some/url", io.NopCloser(strings.NewReader(tc.bodyReq)))

			strBody, wantStrBody := string(resp.Body), string(tc.wantOut.Body)

			if strBody != wantStrBody {
				t.Errorf("want %s and got %s", strBody, wantStrBody)
			}

			if tc.wantOut.Error != nil {
				if resp.Error.Error() != tc.wantOut.Error.Error() {
					t.Errorf("want %v and got %v", tc.wantOut.Error, resp.Error)
				}
			}

			if resp.StatusCode != tc.wantOut.StatusCode {
				t.Errorf("want %d and got %d", tc.wantOut.StatusCode, resp.StatusCode)
			}

		})
	}
}


// cover     ->  go test -v -count=1 -cover -failfast -run ^TestDelete$
// coverHTML ->  go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestDelete$; go tool cover -html=coverage.out
func TestDelete(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantOut *ClientResponse
		mock    *httpMock
	}{
		{
			name: "success",
			ctx:  context.Background(),
			wantOut: &ClientResponse{
				Body:       []byte(`{"data": "quick is awesome!"}`),
				StatusCode: 200,
				Error:      nil,
			},
			mock: &httpMock{
				response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(`{"data": "quick is awesome!"}`)),
				},
				err: nil,
			},
		},
		{
			name: "error_request",
			ctx:  context.Background(),
			wantOut: &ClientResponse{
				Body:       nil,
				StatusCode: 0,
				Error:      errors.New("request_error"),
			},
			mock: &httpMock{
				response: &http.Response{
					StatusCode: 500,
					Body:       nil,
				},
				err: errors.New("request_error"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(*testing.T) {

			c := Client{
				Ctx:        context.Background(),
				ClientHttp: tc.mock,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}

			resp := c.Delete("some/url")

			strBody, wantStrBody := string(resp.Body), string(tc.wantOut.Body)

			if strBody != wantStrBody {
				t.Errorf("want %s and got %s", strBody, wantStrBody)
			}

			if tc.wantOut.Error != nil {
				if resp.Error.Error() != tc.wantOut.Error.Error() {
					t.Errorf("want %v and got %v", tc.wantOut.Error, resp.Error)
				}
			}

			if resp.StatusCode != tc.wantOut.StatusCode {
				t.Errorf("want %d and got %d", tc.wantOut.StatusCode, resp.StatusCode)
			}

		})
	}
}

func BenchmarkGet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Get("http://localhost:8000")
	}
}


func BenchmarkPost(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Post("http://localhost:8000", io.NopCloser(strings.NewReader(`{"data": "quick is awesome!"}`)))
	}
}


func BenchmarkPut(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Put("http://localhost:8000", io.NopCloser(strings.NewReader(`{"data": "quick is awesome!"}`)))
	}
}


func BenchmarkDelete(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Delete("http://localhost:8000")
	}
}
