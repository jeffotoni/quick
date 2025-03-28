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
	"strings"
	"testing"
)

// TestQuick_DeleteGroup ensures that DELETE requests to grouped and standalone routes
// return the correct status codes and response bodies.
//
//	Coverage:
//
// To measure test coverage and open it in HTML:
//
//	go test -v -count=1 -cover -failfast -run ^TestQuick_DeleteGroup$
//	go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestQuick_DeleteGroup$; go tool cover -html=coverage.out
func TestQuick_DeleteGroup(t *testing.T) {
	type args struct {
		route       string
		wantOut     string
		wantCode    int
		isWantedErr bool
		reqHeaders  map[string]string
	}

	testSuccessMockHandler := func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(204).SendString("")
	}

	r := New()
	r.Delete("/", testSuccessMockHandler)
	g1 := r.Group("/del/group")
	g1.Delete("/user", testSuccessMockHandler)
	g1.Delete("/tester/:p1", testSuccessMockHandler)
	r.Delete("/jeff", testSuccessMockHandler)

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				route:       "/",
				wantOut:     ``,
				wantCode:    204,
				isWantedErr: false,
			},
		},
		{
			name: "success",
			args: args{
				route:       "/del/group/user",
				wantOut:     ``,
				wantCode:    204,
				isWantedErr: false,
			},
		},
		{
			name: "success_param",
			args: args{
				route:       "/del/group/tester/:p1",
				wantOut:     ``,
				wantCode:    204,
				isWantedErr: false,
			},
		},
		{
			name: "success_without_param",
			args: args{
				route:       "/jeff",
				wantOut:     ``,
				wantCode:    204,
				isWantedErr: false,
				reqHeaders:  map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			name: "success_without_param",
			args: args{
				route:       "/nao/existe/esta/rota",
				wantOut:     "404 page not found",
				wantCode:    404,
				isWantedErr: false,
				reqHeaders:  map[string]string{"Content-Type": "application/json"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data, err := r.Qtest(QuickTestOptions{
				Method:  "DELETE",
				URI:     tt.args.route,
				Headers: tt.args.reqHeaders,
			})

			if (!tt.args.isWantedErr) && err != nil {
				t.Errorf("error: %v", err)
				return
			}

			s := strings.TrimSpace(data.BodyStr())
			if s != tt.args.wantOut {
				t.Errorf("was suppose to return %s and %s come", tt.args.wantOut, data.BodyStr())
				return
			}

			if tt.args.wantCode != data.StatusCode() {
				t.Errorf("was suppose to return %d and %d come", tt.args.wantCode, data.StatusCode())
				return
			}

			t.Logf("outputBody -> %v", data.BodyStr())

			// t.Logf("\nOutputBodyString -> %v", data.BodyStr())
			// t.Logf("\nStatusCode -> %d", data.StatusCode())
			// t.Logf("\nOutputBody -> %v", string(data.Body())) // I have converted in this example to string but comes []byte as default
			// t.Logf("\nResponse -> %v", data.Response())
		})
	}
}
