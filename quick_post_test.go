package goquick

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/jeffotoni/goquick/internal/concat"
)

// cover     -> go test -v -count=1 -cover -failfast -run ^TestQuick_Post$
// coverHTML -> go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestQuick_Post$; go tool cover -html=coverage.out
func TestQuick_Post(t *testing.T) {
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

	type XmlData struct {
		XMLName xml.Name `xml:"data"`
		Name    string   `xml:"name"`
		Age     int      `xml:"age"`
	}

	type myXmlType struct {
		XMLName xml.Name `xml:"MyXMLType"`
		Data    XmlData  `xml:"data"`
	}

	testSuccessMockHandler := func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		b := c.Body()
		resp := concat.String(`"data":`, string(b))
		c.Status(200)
		return c.SendString(resp)
	}

	testSuccessMockHandlerBodyStr := func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		resp := concat.String(`"data":`, c.BodyString())
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

	testSuccessMockXml := func(c *Ctx) error {
		c.Set("Content-Type", ContentTypeTextXML)
		mtx := new(myXmlType)
		if err := c.Bind(&mtx); err != nil {
			t.Errorf("error: %v", err)
		}
		return c.Status(200).XML(mtx)
	}

	r := New()
	r.Post("/test", testSuccessMockHandler)
	r.Post("/testStr", testSuccessMockHandlerBodyStr)
	r.Post("/tester/:p1", testSuccessMockHandler)
	r.Post("/", testSuccessMockHandlerString)
	r.Post("/bind", testSuccessMockHandlerBind)
	r.Post("/test/xml", testSuccessMockXml)

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
			name: "success_body_str",
			args: args{
				route:       "/testStr",
				wantCode:    200,
				wantOut:     `"data":{"name":"jeff", "age":35}`,
				isWantedErr: false,
				reqBody:     []byte(`{"name":"jeff", "age":35}`),
			},
		},
		{
			name: "success_param",
			args: args{
				route:       "/tester/some",
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
		{
			name: "success_xml",
			args: args{
				route:       "/test/xml",
				wantCode:    200,
				wantOut:     `<MyXMLType><data><name>Jeff</name><age>35</age></data></MyXMLType>`,
				isWantedErr: false,
				reqBody:     []byte(`<MyXMLType><data><name>Jeff</name><age>35</age></data></MyXMLType>`),
				reqHeaders:  map[string]string{"Content-Type": ContentTypeTextXML},
			},
		},
		{
			name: "success_different_valid_json",
			args: args{
				route:       "/test",
				wantCode:    200,
				wantOut:     `"data":{"title":"Test","status":true}`,
				isWantedErr: false,
				reqBody:     []byte(`{"title":"Test","status":true}`),
				reqHeaders:  map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			name: "success_empty_body",
			args: args{
				route:       "/test",
				wantCode:    200,
				wantOut:     `"data":{}`,
				isWantedErr: false,
				reqBody:     []byte(`{}`),
				reqHeaders:  map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			name: "success_json_with_numbers",
			args: args{
				route:       "/test",
				wantCode:    200,
				wantOut:     `"data":{"value":12345,"percentage":99.9}`,
				isWantedErr: false,
				reqBody:     []byte(`{"value":12345,"percentage":99.9}`),
				reqHeaders:  map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			name: "success_xml_with_different_data",
			args: args{
				route:       "/test/xml",
				wantCode:    200,
				wantOut:     `<MyXMLType><data><name>Maria</name><age>28</age></data></MyXMLType>`,
				isWantedErr: false,
				reqBody:     []byte(`<MyXMLType><data><name>Maria</name><age>28</age></data></MyXMLType>`),
				reqHeaders:  map[string]string{"Content-Type": "application/xml"},
			},
		},
		{
			name: "success_longer_json",
			args: args{
				route:       "/test",
				wantCode:    200,
				wantOut:     `"data":{"name":"jeff","age":35,"city":"São Paulo","country":"Brazil"}`,
				isWantedErr: false,
				reqBody:     []byte(`{"name":"jeff","age":35,"city":"São Paulo","country":"Brazil"}`),
				reqHeaders:  map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			name: "success_json_with_array",
			args: args{
				route:       "/test",
				wantCode:    200,
				wantOut:     `"data":{"items":["item1","item2","item3"]}`,
				isWantedErr: false,
				reqBody:     []byte(`{"items":["item1","item2","item3"]}`),
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

func Test_extractParamsPost(t *testing.T) {
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
			if got := extractParamsPost(&tt.args.quick, tt.args.handlerFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractParamsPost() = %v, want %v", got, tt.want)
			}
		})
	}
}
