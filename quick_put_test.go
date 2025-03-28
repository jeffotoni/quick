// Package quick provides a high-performance, minimalistic web framework for Go.
//
// This file contains **unit tests** for various functionalities of the Quick framework.
// These tests ensure that the core features of Quick work as expected.
//
// 📌 To run all unit tests, use:
//
//	$ go test -v ./...
//	$ go test -v
package quick

import (
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/jeffotoni/quick/internal/concat"
)

// TestQuick_Put validates the behavior of PUT requests using various route configurations.
//
// This test ensures that:
// - PUT handlers return the correct status code and body content.
// - Routes with and without parameters behave as expected.
// - JSON payloads are correctly echoed back in the response body.
//
// The test uses a mock handler that sets the Content-Type to "application/json"
// Usage:
// cover     -> go test -v -count=1 -cover -failfast -run ^TestQuick_Put$
// coverHTML -> go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestQuick_Put$; go tool cover -html=coverage.out
func TestQuick_Put(t *testing.T) {
	type args struct {
		route       string
		wantCode    int
		wantOut     string
		isWantedErr bool
		reqBody     []byte
		reqHeaders  map[string]string
	}

	testSuccessMockHandler := func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		b := c.Body()
		resp := concat.String(`"data":`, string(b))
		c.Byte([]byte(resp))
		return nil
	}

	r := New()
	r.Put("/", testSuccessMockHandler)
	r.Put("/test", testSuccessMockHandler)
	r.Put("/tester/:p1", testSuccessMockHandler)
	r.Put("/jeff", testSuccessMockHandler)

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				route:       "/test",
				wantCode:    200,
				wantOut:     `"data":{"name":"jeff", "age":35}`,
				isWantedErr: false,
				reqBody:     []byte(`{"name":"jeff", "age":35}`),
			},
		},
		{
			name: "success_param",
			args: args{
				route:       "/tester/:p1",
				wantCode:    200,
				wantOut:     `"data":{"name":"jeff", "age":35}`,
				isWantedErr: false,
				reqBody:     []byte(`{"name":"jeff", "age":35}`),
			},
		},
		{
			name: "success_without_param",
			args: args{
				route:       "/",
				wantCode:    200,
				wantOut:     `"data":{"name":"jeff", "age":35}`,
				isWantedErr: false,
				reqBody:     []byte(`{"name":"jeff", "age":35}`),
				reqHeaders:  map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			name: "success_without_param",
			args: args{
				route:       "/jeff",
				wantCode:    200,
				wantOut:     `"data":{"name":"jeff", "age":35}`,
				isWantedErr: false,
				reqBody:     []byte(`{"name":"jeff", "age":35}`),
				reqHeaders:  map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			name: "success_json_empty",
			args: args{
				route:       "/jeff",
				wantCode:    200,
				wantOut:     `"data":{}`,
				isWantedErr: false,
				reqBody:     []byte(`{}`),
				reqHeaders:  map[string]string{"Content-Type": "application/json"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data, err := r.Qtest(QuickTestOptions{
				Method:  "PUT",
				URI:     tt.args.route,
				Headers: tt.args.reqHeaders,
				Body:    tt.args.reqBody,
			})

			if (!tt.args.isWantedErr) && err != nil {
				t.Errorf("error: %v", err)
				return
			}

			s := strings.TrimSpace(data.BodyStr())
			if s != tt.args.wantOut {
				t.Errorf("route %s -> was suppose to return %s and %s come", tt.args.route, tt.args.wantOut, data.BodyStr())
				return
			}

			if tt.args.wantCode != data.StatusCode() {
				t.Errorf("was suppose to return %d and %d come", tt.args.wantCode, data.StatusCode())
				return
			}

			t.Logf("\nOutputBodyString -> %v", data.BodyStr())
			t.Logf("\nStatusCode -> %d", data.StatusCode())
			t.Logf("\nOutputBody -> %v", string(data.Body())) // I have converted in this example to string but comes []byte as default
			t.Logf("\nResponse -> %v", data.Response())
		})
	}
}

// Test_extractParamsPut validates the behavior of the `extractParamsPut` function.
//
// This test checks whether the function properly generates an HTTP handler for PUT requests,
// wiring the provided Quick instance and handler function.
//
// Note: Test cases should be added to verify specific route handling, request parsing, and
// response logic. This scaffolding ensures that future implementations can be verified
// using table-driven testing.
//
// Usage:
// - cover:     go test -v -count=1 -cover -failfast -run ^Test_extractParamsPut$
// - coverHTML: go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^Test_extractParamsPut$; go tool cover -html=coverage.out
func Test_extractParamsPut(t *testing.T) {
	type args struct {
		quick       Quick
		handlerFunc func(*Ctx) error
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractParamsPut(&tt.args.quick, tt.args.handlerFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractParamsPut() = %v, want %v", got, tt.want)
			}
		})
	}
}
