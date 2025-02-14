package quick

import (
	"strings"
	"testing"
)

// cover     ->  go test -v -count=1 -cover -failfast -run ^TestQuick_GroupGet$
// coverHTML ->  go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestQuick_GroupGet$; go tool cover -html=coverage.out
func TestQuick_GroupGet(t *testing.T) {
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
		return c.JSON(mt)
	}

	r := New()
	g1 := r.Group("/v1/my")
	g1.Get("/test", testSuccessMockHandler)
	r.Get("/tester/:p1", testSuccessMockHandler)
	r.Get("/", testSuccessMockHandler)
	//g1.Get(`/reg/{(\S+)}/route`, testSuccessMockHandler)

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				route:       "/v1/my/test?some=val1",
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
		{
			name: "success_with_regex",
			args: args{
				route:       "/my/reg/88/route",
				wantOut:     `{"name":"jeff","age":35}`,
				wantCode:    200,
				isWantedErr: false,
			},
		},
		{
			name: "error_not_exists_route",
			args: args{
				route:       "/tester/val1/route",
				wantOut:     `404 page not found`,
				wantCode:    404,
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
