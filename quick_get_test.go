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
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// TestRouteGET tests whether a GET route returns the expected response.
// The status is expected to be 200 and the response body is "Hello, GET!".
//
// Usage:
//
//	go test -v -run TestRouteGET
func TestRouteGET(t *testing.T) {
	q := New()

	q.Get("/v1/user", func(c *Ctx) error {
		return c.String("Hello, GET!")
	})

	data, err := q.Qtest(QuickTestOptions{
		Method:  MethodGet,
		URI:     "/v1/user",
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		t.Errorf("Error during Qtest: %v", err)
		return
	}

	if data.StatusCode() != 200 {
		t.Errorf("Expected status 200, got %d", data.StatusCode())
	}

	if data.BodyStr() != "Hello, GET!" {
		t.Errorf("Expected body 'Hello, GET!', got '%s'", data.BodyStr())
	}
}

// TestQuick_Get tests multiple GET route scenarios with and without parameters, query strings and error cases.
// Valid routes are expected to return status 200 with expected JSON and invalid routes return error.
//
// Usage:
//
//	go test -v -run TestQuick_Get
func TestQuick_Get(t *testing.T) {
	type args struct {
		route       string
		wantCode    int
		wantOut     string
		isWantedErr bool
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
		fmt.Println("More Requests:", c.MoreRequests)
		return c.JSON(mt)
	}

	r := New()
	r.Get("/test", testSuccessMockHandler)
	r.Get("/tester/:p1", testSuccessMockHandler)
	r.Get("/", testSuccessMockHandler)
	//r.Get("/reg/{[0-9]}", testSuccessMockHandler)
	r.Get("/query", func(c *Ctx) error {
		param := c.Request.URL.Query().Get("name")
		if param == "" {
			return c.Status(400).SendString("Falta o parametro de consulta")
		}
		return c.JSON(myType{Name: param, Age: 30})
	})

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				route:       "/test?some=1",
				wantOut:     `{"name":"jeff","age":35}`,
				wantCode:    200,
				isWantedErr: false,
			},
		},
		{
			name: "success_with_params",
			args: args{
				route:       "/tester/val1",
				wantOut:     `{"name":"jeff","age":35}`,
				wantCode:    200,
				isWantedErr: false,
			},
		},
		{
			name: "success_with_nothing",
			args: args{
				route:       "/",
				wantOut:     `{"name":"jeff","age":35}`,
				wantCode:    200,
				isWantedErr: false,
			},
		},
		// {
		// 	name: "success_with_regex",
		// 	args: args{
		// 		route:       "/reg/1",
		// 		wantOut:     `{"name":"jeff","age":35}`,
		// 		wantCode:    200,
		// 		isWantedErr: false,
		// 	},
		// },
		// {
		// 	name: "success_with_different_regex",
		// 	args: args{
		// 		route:       "/reg/5",
		// 		wantOut:     `{"name":"jeff","age":35}`,
		// 		wantCode:    200,
		// 		isWantedErr: false,
		// 	},
		// },
		{
			name: "error_regex_mismatch",
			args: args{
				route:       "/reg/abc",
				wantOut:     "404 page not found",
				wantCode:    404,
				isWantedErr: true,
			},
		},
		{
			name: "success_query_param",
			args: args{
				route:       "/query?name=Alice",
				wantOut:     `{"name":"Alice","age":30}`,
				wantCode:    200,
				isWantedErr: false,
			},
		},
		{
			name: "error_missing_query_param",
			args: args{
				route:       "/query",
				wantOut:     "Falta o parametro de consulta",
				wantCode:    400,
				isWantedErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data, err := r.Qtest(QuickTestOptions{
				Method:  "GET",
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
		})
	}
}

// Test_extractParamsGet tests the function extractParamsGet responsible for extracting route parameters for GET methods.
// It has not yet implemented cases.
//
// Use:
//
//	go test -v -run Test_extractParamsGet
func Test_extractParamsGet(t *testing.T) {
	type args struct {
		quick       Quick
		pathTmp     string
		paramsPath  string
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
			if got := extractParamsGet(&tt.args.quick, tt.args.pathTmp, tt.args.paramsPath, tt.args.handlerFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractParamsGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestQuick_UseCors tests the CORS middleware application in the Quick framework (commented).
// The test simulates routes that should return JSON with middleware applied.
// It is currently commented, but when activated status 200 is expected and JSON expected.
//
// Use:
//
//  go test -v -run TestQuick_UseCors
// func TestQuick_UseCors(t *testing.T) {
// 	type args struct {
// 		route       string
// 		wantCode    int
// 		wantOut     string
// 		isWantedErr bool
// 		reqHeaders  map[string]string
// 	}

// 	type myType struct {
// 		Name string `json:"name"`
// 		Age  int    `json:"age"`
// 	}

// 	mt := myType{}
// 	mt.Name = "jeff"
// 	mt.Age = 35

// 	testSuccessMockHandler := func(c *Ctx) error {
// 		c.Set("Content-Type", "application/json")
// 		return c.Status(200).JSON(mt)
// 	}

// 	r := New()
// 	r.Use(cors.New(), "cors")
// 	r.Get("/test", testSuccessMockHandler)
// 	r.Get("/tester/:p1", testSuccessMockHandler)
// 	r.Get("/", testSuccessMockHandler)
// 	r.Get("/reg/{[0-9]}", testSuccessMockHandler)

// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		{
// 			name: "success",
// 			args: args{
// 				route:       "/test?some=1",
// 				wantOut:     `{"name":"jeff","age":35}`,
// 				wantCode:    200,
// 				isWantedErr: false,
// 			},
// 		},
// 		// {
// 		// 	name: "success_with_params",
// 		// 	args: args{
// 		// 		route:       "/tester/val1",
// 		// 		wantOut:     `{"name":"jeff","age":35}`,
// 		// 		wantCode:    200,
// 		// 		isWantedErr: false,
// 		// 	},
// 		// },
// 		// {
// 		// 	name: "success_with_nothing",
// 		// 	args: args{
// 		// 		route:       "/",
// 		// 		wantOut:     `{"name":"jeff","age":35}`,
// 		// 		wantCode:    200,
// 		// 		isWantedErr: false,
// 		// 	},
// 		// },
// 		// {
// 		// 	name: "success_with_regex",
// 		// 	args: args{
// 		// 		route:       "/reg/1",
// 		// 		wantOut:     `{"name":"jeff","age":35}`,
// 		// 		wantCode:    200,
// 		// 		isWantedErr: false,
// 		// 	},
// 		// },
// 		// {
// 		// 	name: "error_not_exists_route",
// 		// 	args: args{
// 		// 		route:       "/tester/val1/route",
// 		// 		wantOut:     `404 page not found`,
// 		// 		wantCode:    404,
// 		// 		isWantedErr: true,
// 		// 	},
// 		// },
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			data, err := r.QuickTest("GET", tt.args.route, tt.args.reqHeaders)
// 			if (!tt.args.isWantedErr) && err != nil {
// 				t.Errorf("error: %v", err)
// 				return
// 			}

// 			s := strings.TrimSpace(data.BodyStr())
// 			if s != tt.args.wantOut {
// 				t.Errorf("was suppose to return %s and %s come", tt.args.wantOut, data.BodyStr())
// 				return
// 			}

// 			if tt.args.wantCode != data.StatusCode() {
// 				t.Errorf("was suppose to return %d and %d come", tt.args.wantCode, data.StatusCode())
// 				return
// 			}

// 			t.Logf("outputBody -> %v", data.BodyStr())
// 		})
// 	}
// }
