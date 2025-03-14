package quick

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

/*
+==============================================================================================================================+
#     To test the entire package and check the coverage you can run those commands below:                                      #
#                                                                                                                              #
#     coverage     -> go test -v -count=1 -failfast -cover -coverprofile=coverage.out ./...                                    #
#     coverageHTML -> go test -v -count=1 -failfast -cover -coverprofile=coverage.out ./...; go tool cover -html=coverage.out  #
+==============================================================================================================================+
*/
// cover     ->  go test -v -count=1 -cover -failfast -run ^TestRouteGET$
// coverHTML ->  go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestRouteGET$; go tool cover -html=coverage.out

func TestRouteGET(t *testing.T) {
	q := New()

	q.Get("/v1/user", func(c *Ctx) error {
		return c.String("Hello, GET!")
	})

	data, err := q.QuickTest("GET", "/v1/user", nil)
	if err != nil {
		t.Errorf("Error during QuickTest: %v", err)
		return
	}

	if data.StatusCode() != 200 {
		t.Errorf("Expected status 200, got %d", data.StatusCode())
	}

	if data.BodyStr() != "Hello, GET!" {
		t.Errorf("Expected body 'Hello, GET!', got '%s'", data.BodyStr())
	}
}

// cover     ->  go test -v -count=1 -cover -failfast -run ^TestQuick_Get$
// coverHTML ->  go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestQuick_Get$; go tool cover -html=coverage.out
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

			data, err := r.QuickTest("GET", tt.args.route, tt.args.reqHeaders)
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

// cover     ->  go test -v -count=1 -cover -failfast -run ^Test_extractParamsGet$
// coverHTML ->  go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^Test_extractParamsGet$; go tool cover -html=coverage.out
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

// cover     -> go test -v -count=1 -cover -failfast -run ^TestQuick_UseCors$
// coverHTML -> go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestQuick_UseCors$; go tool cover -html=coverage.out
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
