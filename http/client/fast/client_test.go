package fast

import (
	"context"
	"errors"

	"testing"
)

// cover     ->  go test -v -count=1 -cover -failfast -run ^TestGet$
// coverHTML ->  go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestGet$; go tool cover -html=coverage.out
func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		URL     string
		ctx     context.Context
		wantOut *ClientResponse
		mock    *httpMock
		wantErr bool
	}{
		{
			name: "success",
			ctx:  context.Background(),
			URL:  "https://letsgoquick.com",
			wantOut: &ClientResponse{
				Body:       []byte(`{"data": "quick is awesome!"}`),
				StatusCode: 200,
			},
			wantErr: false,
			mock: &httpMock{
				err: nil,
			},
		},
		{
			name: "error_request",
			ctx:  context.Background(),
			URL:  "https://letsgoquick.com",
			wantOut: &ClientResponse{
				Body:       nil,
				StatusCode: 0,
			},
			wantErr: true,
			mock: &httpMock{
				err: errors.New("request_error"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(*testing.T) {

			c := Client{
				Ctx:            context.Background(),
				ClientFastHttp: tc.mock,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Timeout: 10000,
			}

			resp, err := c.Get(tc.URL)
			if (!tc.wantErr) && err != nil {
				t.Errorf("want nil and got %v", err)
			}

			if resp != nil {
				strBody, wantStrBody := string(resp.Body), string(tc.wantOut.Body)

				if strBody != wantStrBody {
					t.Errorf("want %s and got %s", strBody, wantStrBody)
				}

				if resp.StatusCode != tc.wantOut.StatusCode {
					t.Errorf("want %d and got %d", tc.wantOut.StatusCode, resp.StatusCode)
				}
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
		URL     string
		bodyReq string
		wantOut *ClientResponse
		mock    *httpMock
		wantErr bool
	}{
		{
			name:    "success",
			ctx:     context.Background(),
			URL:     "https://letsgoquick.com",
			bodyReq: `{"data": "quick is awesome!"}`,
			wantOut: &ClientResponse{
				Body:       []byte(`{"data": "quick is awesome!"}`),
				StatusCode: 200,
			},
			wantErr: false,
			mock: &httpMock{
				err: nil,
			},
		},
		{
			name: "error_request",
			ctx:  context.Background(),
			URL:  "https://letsgoquick.com",
			wantOut: &ClientResponse{
				Body:       nil,
				StatusCode: 0,
			},
			wantErr: true,
			mock: &httpMock{
				err: errors.New("request_error"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(*testing.T) {

			c := Client{
				Ctx:            context.Background(),
				ClientFastHttp: tc.mock,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}

			resp, err := c.Post(tc.URL, []byte(tc.bodyReq))
			if (!tc.wantErr) && err != nil {
				t.Errorf("want nil and got %v", err)
			}

			if resp != nil {
				strBody, wantStrBody := string(resp.Body), string(tc.wantOut.Body)

				if strBody != wantStrBody {
					t.Errorf("want %s and got %s", strBody, wantStrBody)
				}

				if resp.StatusCode != tc.wantOut.StatusCode {
					t.Errorf("want %d and got %d", tc.wantOut.StatusCode, resp.StatusCode)
				}
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
		URL     string
		wantOut *ClientResponse
		mock    *httpMock
		wantErr bool
	}{
		{
			name:    "success",
			ctx:     context.Background(),
			bodyReq: `{"data": "quick is awesome!"}`,
			URL:     "https://letsgoquick.com",
			wantOut: &ClientResponse{
				Body:       []byte(`{"data": "quick is awesome!"}`),
				StatusCode: 200,
			},
			wantErr: false,
			mock: &httpMock{
				err: nil,
			},
		},
		{
			name: "error_request",
			ctx:  context.Background(),
			wantOut: &ClientResponse{
				Body:       nil,
				StatusCode: 0,
			},
			wantErr: true,
			URL:     "https://letsgoquick.com",
			mock: &httpMock{
				err: errors.New("request_error"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(*testing.T) {

			c := Client{
				Ctx:            context.Background(),
				ClientFastHttp: tc.mock,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}

			resp, err := c.Put(tc.URL, []byte(tc.bodyReq))
			if (!tc.wantErr) && err != nil {
				t.Errorf("want nil and got %v", err)
			}

			if resp != nil {
				strBody, wantStrBody := string(resp.Body), string(tc.wantOut.Body)

				if strBody != wantStrBody {
					t.Errorf("want %s and got %s", strBody, wantStrBody)
				}

				if resp.StatusCode != tc.wantOut.StatusCode {
					t.Errorf("want %d and got %d", tc.wantOut.StatusCode, resp.StatusCode)
				}
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
		URL     string
		wantOut *ClientResponse
		mock    *httpMock
		wantErr bool
	}{
		{
			name: "success",
			ctx:  context.Background(),
			URL:  "https://letsgoquick.com",
			wantOut: &ClientResponse{
				Body:       []byte(`{"data": "quick is awesome!"}`),
				StatusCode: 200,
			},
			wantErr: false,
			mock: &httpMock{
				err: nil,
			},
		},
		{
			name: "error_request",
			ctx:  context.Background(),
			URL:  "https://letsgoquick.com",
			wantOut: &ClientResponse{
				Body:       nil,
				StatusCode: 0,
			},
			wantErr: true,
			mock: &httpMock{
				err: errors.New("request_error"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(*testing.T) {

			c := Client{
				Ctx:            context.Background(),
				ClientFastHttp: tc.mock,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}

			resp, err := c.Delete(tc.URL)

			if (!tc.wantErr) && err != nil {
				t.Errorf("want nil and got %v", err)
			}

			if resp != nil {
				strBody, wantStrBody := string(resp.Body), string(tc.wantOut.Body)

				if strBody != wantStrBody {
					t.Errorf("want %s and got %s", strBody, wantStrBody)
				}

				if resp.StatusCode != tc.wantOut.StatusCode {
					t.Errorf("want %d and got %d", tc.wantOut.StatusCode, resp.StatusCode)
				}
			}

		})
	}
}

// cover     ->  go test -v -count=1 -cover -failfast -run ^TestExternalGet$
// coverHTML ->  go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestExternalGet$; go tool cover -html=coverage.out
func TestExternalGet(t *testing.T) {
	tests := []struct {
		name    string
		URL     string
		ctx     context.Context
		wantOut *ClientResponse
		wantErr bool
	}{
		{
			name: "success",
			ctx:  context.Background(),
			URL:  "https://letsgoquick.com",
			wantOut: &ClientResponse{
				Body:       []byte(letsgoquickOutMock),
				StatusCode: 200,
			},
			wantErr: false,
		},
		{
			name: "error_request_unsupported_protocol",
			ctx:  context.Background(),
			URL:  "letsgoquick.com",
			wantOut: &ClientResponse{
				Body:       nil,
				StatusCode: 0,
			},
			wantErr: true,
		},
		{
			name: "error_request_no_such_host",
			ctx:  context.Background(),
			URL:  "https://letsgoquick.co",
			wantOut: &ClientResponse{
				Body:       nil,
				StatusCode: 0,
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(*testing.T) {

			resp, err := Get(tc.URL)
			if (!tc.wantErr) && err != nil {
				t.Errorf("want nil and got %v", err)
			}

			if resp != nil {

				removeSpaces(&resp.Body)
				strBody, wantStrBody := string(resp.Body), string(tc.wantOut.Body)

				if strBody != wantStrBody {
					t.Errorf("want %s and got %s", strBody, wantStrBody)
				}

				if resp.StatusCode != tc.wantOut.StatusCode {
					t.Errorf("want %d and got %d", tc.wantOut.StatusCode, resp.StatusCode)
				}
			}

		})
	}
}

// cover     ->  go test -v -count=1 -cover -failfast -run ^TestExternalPost$
// coverHTML ->  go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestExternalPost$; go tool cover -html=coverage.out
func TestExternalPost(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		URL     string
		bodyReq string
		wantOut *ClientResponse
		mock    *httpMock
		wantErr bool
	}{
		{
			name:    "success",
			ctx:     context.Background(),
			URL:     "https://go.dev/",
			bodyReq: `{"data": "quick is awesome!"}`,
			wantOut: &ClientResponse{
				StatusCode: 200,
			},
			wantErr: false,
			mock: &httpMock{
				err: nil,
			},
		},
		{
			name:    "error_unsuported_procol",
			ctx:     context.Background(),
			URL:     "letsgoquick.com",
			bodyReq: `{"data": "quick is awesome!"}`,
			wantOut: &ClientResponse{
				StatusCode: 200,
			},
			wantErr: true,
		},
		{
			name: "error_request_no_such_host",
			ctx:  context.Background(),
			URL:  "https://letsgoquick.co",
			wantOut: &ClientResponse{
				Body:       nil,
				StatusCode: 0,
			},
			wantErr: true,
		},
		{
			name: "error_request_403",
			ctx:  context.Background(),
			URL:  "https://letsgoquick.com",
			wantOut: &ClientResponse{
				StatusCode: 403,
			},
			wantErr: true,
			mock: &httpMock{
				err: errors.New("request_error"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(*testing.T) {

			resp, err := Post(tc.URL, []byte(tc.bodyReq))
			if (!tc.wantErr) && err != nil {
				t.Errorf("want nil and got %v", err)
			}

			if resp != nil {
				removeSpaces(&resp.Body)
				if resp.StatusCode != tc.wantOut.StatusCode {
					t.Errorf("want %d and got %d", tc.wantOut.StatusCode, resp.StatusCode)
				}
			}

		})
	}
}

// cover     ->  go test -v -count=1 -cover -failfast -run ^TestExternalPut$
// coverHTML ->  go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestExternalPut$; go tool cover -html=coverage.out
func TestExternalPut(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		URL     string
		bodyReq string
		wantOut *ClientResponse
		mock    *httpMock
		wantErr bool
	}{
		{
			name:    "success",
			ctx:     context.Background(),
			URL:     "https://go.dev/",
			bodyReq: `{"data": "quick is awesome!"}`,
			wantOut: &ClientResponse{
				StatusCode: 200,
			},
			wantErr: false,
		},
		{
			name:    "error_unsuported_procol",
			ctx:     context.Background(),
			URL:     "letsgoquick.com",
			bodyReq: `{"data": "quick is awesome!"}`,
			wantOut: &ClientResponse{
				StatusCode: 200,
			},
			wantErr: true,
		},
		{
			name: "error_request_no_such_host",
			ctx:  context.Background(),
			URL:  "https://letsgoquick.co",
			wantOut: &ClientResponse{
				Body:       nil,
				StatusCode: 0,
			},
			wantErr: true,
		},
		{
			name: "error_request_403",
			ctx:  context.Background(),
			URL:  "https://letsgoquick.com",
			wantOut: &ClientResponse{
				StatusCode: 403,
			},
			wantErr: true,
			mock: &httpMock{
				err: errors.New("request_error"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(*testing.T) {

			resp, err := Put(tc.URL, []byte(tc.bodyReq))
			if (!tc.wantErr) && err != nil {
				t.Errorf("want nil and got %v", err)
			}

			if resp != nil {
				removeSpaces(&resp.Body)
				if resp.StatusCode != tc.wantOut.StatusCode {
					t.Errorf("want %d and got %d", tc.wantOut.StatusCode, resp.StatusCode)
				}
			}

		})
	}
}

// cover     ->  go test -v -count=1 -cover -failfast -run ^TestExternalDelete$
// coverHTML ->  go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestExternalDelete$; go tool cover -html=coverage.out
func TestExternalDelete(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		URL     string
		bodyReq string
		wantOut *ClientResponse
		mock    *httpMock
		wantErr bool
	}{
		{
			name:    "success",
			ctx:     context.Background(),
			URL:     "https://go.dev/",
			bodyReq: `{"data": "quick is awesome!"}`,
			wantOut: &ClientResponse{
				StatusCode: 200,
			},
			wantErr: false,
		},
		{
			name:    "error_unsuported_procol",
			ctx:     context.Background(),
			URL:     "letsgoquick.com",
			bodyReq: `{"data": "quick is awesome!"}`,
			wantOut: &ClientResponse{
				StatusCode: 200,
			},
			wantErr: true,
		},
		{
			name: "error_request_no_such_host",
			ctx:  context.Background(),
			URL:  "https://letsgoquick.co",
			wantOut: &ClientResponse{
				Body:       nil,
				StatusCode: 0,
			},
			wantErr: true,
		},
		{
			name: "error_request_403",
			ctx:  context.Background(),
			URL:  "https://letsgoquick.com",
			wantOut: &ClientResponse{
				StatusCode: 403,
			},
			wantErr: true,
			mock: &httpMock{
				err: errors.New("request_error"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(*testing.T) {

			resp, err := Delete(tc.URL)
			if (!tc.wantErr) && err != nil {
				t.Errorf("want nil and got %v", err)
			}

			if resp != nil {
				removeSpaces(&resp.Body)
				if resp.StatusCode != tc.wantOut.StatusCode {
					t.Errorf("want %d and got %d", tc.wantOut.StatusCode, resp.StatusCode)
				}
			}

		})
	}
}

func BenchmarkGet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Get("https://letsgoquick:8000")
	}
}

func BenchmarkPost(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Post("https://letsgoquick:8000", []byte(`{"data": "quick is awesome!"}`))
	}
}

func BenchmarkPut(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Put("https://letsgoquick:8000", []byte(`{"data": "quick is awesome!"}`))
	}
}

func BenchmarkDelete(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Delete("https://letsgoquick:8000")
	}
}
