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
	"testing"
)

// TestQuick_Delete verifies the behavior of the DELETE method handler in the Quick framework.
//
// It tests different scenarios including:
// - DELETE requests to static routes
// - DELETE requests with and without body
// - DELETE requests to parameterized routes
// - DELETE requests to non-existent routes
//
// To run only this test:
//
//	go test -v -count=1 -cover -failfast -run ^TestQuick_Delete$
//
// To generate HTML coverage:
//
//	go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestQuick_Delete$; go tool cover -html=coverage.out
func TestQuick_Delete(t *testing.T) {
	type args struct {
		route       string
		wantCode    int
		isWantedErr bool
		reqBody     []byte
		reqHeaders  map[string]string
	}

	type myType struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	mt := myType{}
	mt.Name = "jeff"
	mt.Age = 35

	testSuccessMockHandler := func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.JSON(mt)
	}

	r := New()
	r.Delete("/", testSuccessMockHandler)
	r.Delete("/test", testSuccessMockHandler)
	r.Delete("/tester/:p1", testSuccessMockHandler)
	r.Delete("/jeff", testSuccessMockHandler)

	if len(r.routes) != 4 {
		t.Errorf("was supose have 4 routes, got: %v", len(r.routes))
		return
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				route:       "/test",
				wantCode:    200,
				isWantedErr: false,
			},
		},
		{
			name: "success_with_body_and_ignore",
			args: args{
				route:       "/test",
				wantCode:    200,
				reqBody:     []byte(`{"name":"jeff", "age":35}`),
				reqHeaders:  map[string]string{"Content-Type": "application/json"},
				isWantedErr: false,
			},
		},
		{
			name: "success_param",
			args: args{
				route:       "/tester/:p1",
				wantCode:    200,
				isWantedErr: false,
			},
		},
		{
			name: "success_without_param",
			args: args{
				route:       "/",
				wantCode:    200,
				isWantedErr: false,
			},
		},
		{
			name: "success_without_param",
			args: args{
				route:       "/jeff",
				wantCode:    200,
				isWantedErr: false,
			},
		},
		{
			name: "error_not_exists_route",
			args: args{
				route:       "/tester/val1/route",
				wantCode:    404,
				isWantedErr: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO method DELETE is not acceptable
			data, err := r.QuickTest("DELETE", tt.args.route, tt.args.reqHeaders, tt.args.reqBody)
			if (!tt.args.isWantedErr) && err != nil {
				t.Errorf("error: %v", err)
				return
			}

			if tt.args.wantCode != data.StatusCode() {
				t.Errorf("%s: was suppose to return %d and %d come", tt.name, tt.args.wantCode, data.StatusCode())
				return
			}

			t.Logf("\nOutputBodyString -> %v", data.BodyStr())
			t.Logf("\nStatusCode -> %d", data.StatusCode())
			t.Logf("\nOutputBody -> %v", string(data.Body())) // I have converted in this example to string but comes []byte as default
			t.Logf("\nResponse -> %v", data.Response())
		})
	}
}

// Test_extractParamsDelete tests the extractParamsDelete function, which binds
// route parameters to the request context in a DELETE route.
//
// This is a placeholder for future detailed test cases.
//
// To run only this test:
//
//	go test -v -count=1 -cover -failfast -run ^Test_extractParamsDelete$
//
// To generate HTML coverage:
//
//	go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^Test_extractParamsDelete$; go tool cover -html=coverage.out
func Test_extractParamsDelete(t *testing.T) {
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
			if got := extractParamsDelete(&tt.args.quick, tt.args.handlerFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractParamsDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}
