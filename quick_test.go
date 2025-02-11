package quick

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jeffotoni/quick/internal/concat"
	"github.com/jeffotoni/quick/middleware/cors"
)

/*
+==============================================================================================================================+
#     To test the entire package and check the coverage you can run those commands below:                                      #
#                                                                                                                              #
#     coverage     -> go test -v -count=1 -failfast -cover -coverprofile=coverage.out ./...                                    #
#     coverageHTML -> go test -v -count=1 -failfast -cover -coverprofile=coverage.out ./...; go tool cover -html=coverage.out  #
+==============================================================================================================================+
*/

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
			data, err := r.QuickTest("PUT", tt.args.route, tt.args.reqHeaders, tt.args.reqBody)
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

func TestQuick_ServeStaticFile(t *testing.T) {
	type fields struct {
		routes  []*Route
		mux     *http.ServeMux
		handler http.Handler
	}
	type args struct {
		pattern     string
		handlerFunc func(*Ctx) error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Quick{
				routes: tt.fields.routes,

				mux:     tt.fields.mux,
				handler: tt.fields.handler,
			}
			r.Get(tt.args.pattern, tt.args.handlerFunc)
		})
	}
}

func TestQuick_ServeHTTP(t *testing.T) {
	type fields struct {
		routes  []*Route
		mux     *http.ServeMux
		handler http.Handler
	}
	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quick{
				routes: tt.fields.routes,

				mux:     tt.fields.mux,
				handler: tt.fields.handler,
			}
			q.ServeHTTP(tt.args.w, tt.args.req)
		})
	}
}

func TestCtx_Json(t *testing.T) {
	type fields struct {
		Response http.ResponseWriter
		Request  *http.Request
		Headers  map[string][]string
		Params   map[string]string
		Query    map[string]string
		JSON     map[string]interface{}
		BodyByte []byte
		JsonStr  string
	}
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Ctx{
				Response: tt.fields.Response,
				Request:  tt.fields.Request,
				Headers:  tt.fields.Headers,
				Params:   tt.fields.Params,
				Query:    tt.fields.Query,
				//JSON:     tt.fields.JSON,
				bodyByte: tt.fields.BodyByte,
				JsonStr:  tt.fields.JsonStr,
			}
			if err := c.JSON(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Ctx.Json() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuick_GetRoute(t *testing.T) {
	type fields struct {
		routes  []*Route
		mux     *http.ServeMux
		handler http.Handler
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Route
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Quick{
				routes: tt.fields.routes,

				mux:     tt.fields.mux,
				handler: tt.fields.handler,
			}
			if got := r.GetRoute(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Quick.GetRoute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuick_Listen(t *testing.T) {
	type fields struct {
		routes  []*Route
		mux     *http.ServeMux
		handler http.Handler
	}
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quick{
				routes: tt.fields.routes,

				mux:     tt.fields.mux,
				handler: tt.fields.handler,
			}
			if err := q.Listen(tt.args.addr); (err != nil) != tt.wantErr {
				t.Errorf("Quick.Listen() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func benchmarkWriteToStdout(b *testing.B, size int) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		os.Stdout.Write(make([]byte, size))
	}
}

func benchmarkPrintln(b *testing.B, size int) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fmt.Println(make([]byte, size))
	}
}

func BenchmarkWriteToStdout_10Bytes(b *testing.B) {
	benchmarkWriteToStdout(b, 10)
}

func BenchmarkPrintln_10Bytes(b *testing.B) {
	benchmarkPrintln(b, 10)
}

func BenchmarkWriteToStdout_100Bytes(b *testing.B) {
	benchmarkWriteToStdout(b, 100)
}

func BenchmarkPrintln_100Bytes(b *testing.B) {
	benchmarkPrintln(b, 100)
}

func BenchmarkWriteToStdout_1000Bytes(b *testing.B) {
	benchmarkWriteToStdout(b, 1000)
}

func BenchmarkPrintln_1000Bytes(b *testing.B) {
	benchmarkPrintln(b, 1000)
}

// go test -v -count=1 -failfast -run ^Test_extractParamsPattern$
func Test_extractParamsPattern(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name              string
		args              args
		wantPath          string
		wantParams        string
		wantPartternExist string
	}{
		{
			name: "should ble able to extract 1 param",
			args: args{
				pattern: "/v1/customer/:param1",
			},
			wantPath:          "/v1/customer",
			wantParams:        "/:param1",
			wantPartternExist: "/v1/customer/:param1",
		},
		{
			name: "should ble able to extract 2 params",
			args: args{
				pattern: "/v1/customer/params/:param1/:param2",
			},
			wantPath:          "/v1/customer/params",
			wantParams:        "/:param1/:param2",
			wantPartternExist: "/v1/customer/params/:param1/:param2",
		},
		{
			name: "should ble able to extract 3 params",
			args: args{
				pattern: "/v1/customer/params/:param1/:param2/some/:param3",
			},
			wantPath:          "/v1/customer/params",
			wantParams:        "/:param1/:param2/some/:param3",
			wantPartternExist: "/v1/customer/params/:param1/:param2/some/:param3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, gotParams, gotPartternExist := extractParamsPattern(tt.args.pattern)
			if gotPath != tt.wantPath {
				t.Errorf("extractParamsPattern() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
			if gotParams != tt.wantParams {
				t.Errorf("extractParamsPattern() gotParams = %v, want %v", gotParams, tt.wantParams)
			}
			if gotPartternExist != tt.wantPartternExist {
				t.Errorf("extractParamsPattern() gotPartternExist = %v, want %v", gotPartternExist, tt.wantPartternExist)
			}
		})
	}
}

// cover     -> go test -v -count=1 -cover -failfast -run ^TestQuick_UseCors$
// coverHTML -> go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestQuick_UseCors$; go tool cover -html=coverage.out
func TestQuick_UseCors(t *testing.T) {
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
		return c.Status(200).JSON(mt)
	}

	r := New()
	r.Use(cors.New(), "cors")
	r.Get("/test", testSuccessMockHandler)
	r.Get("/tester/:p1", testSuccessMockHandler)
	r.Get("/", testSuccessMockHandler)
	r.Get("/reg/{[0-9]}", testSuccessMockHandler)

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
		// {
		// 	name: "success_with_params",
		// 	args: args{
		// 		route:       "/tester/val1",
		// 		wantOut:     `{"name":"jeff","age":35}`,
		// 		wantCode:    200,
		// 		isWantedErr: false,
		// 	},
		// },
		// {
		// 	name: "success_with_nothing",
		// 	args: args{
		// 		route:       "/",
		// 		wantOut:     `{"name":"jeff","age":35}`,
		// 		wantCode:    200,
		// 		isWantedErr: false,
		// 	},
		// },
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
		// 	name: "error_not_exists_route",
		// 	args: args{
		// 		route:       "/tester/val1/route",
		// 		wantOut:     `404 page not found`,
		// 		wantCode:    404,
		// 		isWantedErr: true,
		// 	},
		// },
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

// cover -> go test -v -cover -run TestDefaultConfig
// cover -> go test -v -cover -run TestQuickInitializationWithCustomConfig
// cover -> go test -v -cover -run TestQuickInitializationDefaults
// cover -> go test -v -cover -run TestQuickInitializationWithZeroValues
func TestDefaultConfig(t *testing.T) {
	expectedConfig := Config{
		BodyLimit:         2 * 1024 * 1024,
		MaxBodySize:       2 * 1024 * 1024,
		MaxHeaderBytes:    1 * 1024 * 1024,
		RouteCapacity:     1000,
		MoreRequests:      290,
		ReadTimeout:       0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		ReadHeaderTimeout: 0,
	}

	if defaultConfig != expectedConfig {
		t.Errorf("esperado %+v, mas obteve %+v", expectedConfig, defaultConfig)
	}
}

func TestQuickInitializationWithCustomConfig(t *testing.T) {
	customConfig := Config{
		BodyLimit:         4 * 1024 * 1024,
		MaxBodySize:       4 * 1024 * 1024,
		MaxHeaderBytes:    2 * 1024 * 1024,
		RouteCapacity:     500,
		MoreRequests:      500,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       2 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
	}

	q := New(customConfig)

	if q.config != customConfig {
		t.Errorf("esperado %+v, mas obteve %+v", customConfig, q.config)
	}
}

func TestQuickInitializationDefaults(t *testing.T) {
	q := New()

	if q.config.BodyLimit != defaultConfig.BodyLimit {
		t.Errorf("BodyLimit incorreto: esperado %d, obteve %d", defaultConfig.BodyLimit, q.config.BodyLimit)
	}
	if q.config.MaxBodySize != defaultConfig.MaxBodySize {
		t.Errorf("MaxBodySize incorreto: esperado %d, obteve %d", defaultConfig.MaxBodySize, q.config.MaxBodySize)
	}
	if q.config.MoreRequests != defaultConfig.MoreRequests {
		t.Errorf("MoreRequests incorreto: esperado %d, obteve %d", defaultConfig.MoreRequests, q.config.MoreRequests)
	}
}

func TestQuickInitializationWithZeroValues(t *testing.T) {
	zeroConfig := Config{}
	q := New(zeroConfig)

	if q.config.RouteCapacity != 1000 {
		t.Errorf("RouteCapacity incorreto: esperado 1000, obteve %d", q.config.RouteCapacity)
	}
}
