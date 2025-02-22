package quick

import (
	"net/http"
	"reflect"
	"testing"
)

// cover     -> go test -v -count=1 -cover -failfast -run ^TestQuick_Delete$
// coverHTML -> go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestQuick_Delete$; go tool cover -html=coverage.out
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
