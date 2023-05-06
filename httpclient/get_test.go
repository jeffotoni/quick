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
			confHttpGo := NewGoClientConfig(
				ClientGoConfig{
					Transport: &http.Transport{
						DisableKeepAlives: false,
						MaxIdleConns:      100,
					},
					Timeout: 0,
				},
			)

			confHttpGo = tc.mock

			cl := New(HttpClient{
				Ctx:        tc.ctx,
				ClientHttp: confHttpGo,
			})

			resp := cl.Get("")

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

func BenchmarkNew(b *testing.B) {
	for n := 0; n < b.N; n++ {
		New()
	}
}

func BenchmarkGet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		c := New()
		c.Get("")
	}
}
