package quick

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/jeffotoni/goquick/internal/concat"
)

// cover     ->  go test -v -count=1 -cover -failfast -run ^TestQuick_GroupPost$
// coverHTML ->  go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestQuick_GroupPost$; go tool cover -html=coverage.out
func TestQuick_GroupPost(t *testing.T) {
	type args struct {
		route       string
		wantCode    int
		wantOut     string
		isWantedErr bool
		reqBody     []byte
		reqHeaders  map[string]string
	}

	type myType struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	testSuccessMockHandler := func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		b := c.Body()
		resp := concat.String(`"data":`, string(b))
		c.Status(200)
		return c.SendString(resp)
	}

	testSuccessMockHandlerString := func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		mt := new(myType)
		if err := c.BodyParser(mt); err != nil {
			t.Errorf("error: %v", err)
		}
		b, _ := json.Marshal(mt)
		resp := concat.String(`"data":`, string(b))
		c.Status(200)
		return c.String(resp)
	}

	testSuccessMockHandlerBind := func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		mt := new(myType)
		if err := c.Bind(&mt); err != nil {
			t.Errorf("error: %v", err)
		}
		b, _ := json.Marshal(mt)
		resp := concat.String(`"data":`, string(b))
		c.Status(200)
		return c.String(resp)
	}

	r := New()
	g1 := r.Group("/my/group")
	g1.Post("/test", testSuccessMockHandler)
	g1.Post("/tester/:p1", testSuccessMockHandler)
	g2 := r.Group("/my")
	g2.Post("/", testSuccessMockHandlerString)
	r.Post("/bind", testSuccessMockHandlerBind)
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				route:       "/my/group/test",
				wantCode:    200,
				wantOut:     `"data":{"name":"jeff", "age":35}`,
				isWantedErr: false,
				reqBody:     []byte(`{"name":"jeff", "age":35}`),
			},
		},
		{
			name: "success_param",
			args: args{
				route:       "/my/group/tester/:p1",
				wantCode:    200,
				wantOut:     `"data":{"name":"jeff", "age":35}`,
				isWantedErr: false,
				reqBody:     []byte(`{"name":"jeff", "age":35}`),
			},
		},
		{
			name: "success_without_param",
			args: args{
				route:       "/my/",
				wantCode:    200,
				wantOut:     `"data":{"name":"jeff","age":35}`,
				isWantedErr: false,
				reqBody:     []byte(`{"name":"jeff","age":35}`),
				reqHeaders:  map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			name: "success_bind",
			args: args{
				route:       "/bind",
				wantCode:    200,
				wantOut:     `"data":{"name":"jeff","age":35}`,
				isWantedErr: false,
				reqBody:     []byte(`{"name":"jeff","age":35}`),
				reqHeaders:  map[string]string{"Content-Type": "application/json"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data, err := r.QuickTest("POST", tt.args.route, tt.args.reqHeaders, tt.args.reqBody)
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
