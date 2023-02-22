package quick

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Quick
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuick_Post(t *testing.T) {
	type fields struct {
		routes      []Route
		middlewares []http.HandlerFunc
		mux         *http.ServeMux
	}
	type args struct {
		pattern     string
		handlerFunc func(*Ctx)
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
				routes:      tt.fields.routes,
				middlewares: tt.fields.middlewares,
				mux:         tt.fields.mux,
			}
			r.Post(tt.args.pattern, tt.args.handlerFunc)
		})
	}
}

func Test_extractParamsPost(t *testing.T) {
	type args struct {
		pathTmp     string
		handlerFunc func(*Ctx)
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
			if got := extractParamsPost(tt.args.pathTmp, tt.args.handlerFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractParamsPost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCtx_Body(t *testing.T) {
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
				JSON:     tt.fields.JSON,
				BodyByte: tt.fields.BodyByte,
				JsonStr:  tt.fields.JsonStr,
			}
			if err := c.Body(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Ctx.Body() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCtx_BodyString(t *testing.T) {
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
	tests := []struct {
		name   string
		fields fields
		want   string
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
				JSON:     tt.fields.JSON,
				BodyByte: tt.fields.BodyByte,
				JsonStr:  tt.fields.JsonStr,
			}
			if got := c.BodyString(); got != tt.want {
				t.Errorf("Ctx.BodyString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuick_Get(t *testing.T) {
	type fields struct {
		routes      []Route
		middlewares []http.HandlerFunc
		mux         *http.ServeMux
	}
	type args struct {
		pattern     string
		handlerFunc func(*Ctx)
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
				routes:      tt.fields.routes,
				middlewares: tt.fields.middlewares,
				mux:         tt.fields.mux,
			}
			r.Get(tt.args.pattern, tt.args.handlerFunc)
		})
	}
}

func Test_extractParamsGet(t *testing.T) {
	type args struct {
		pathTmp     string
		paramsPath  string
		handlerFunc func(*Ctx)
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
			if got := extractParamsGet(tt.args.pathTmp, tt.args.paramsPath, tt.args.handlerFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractParamsGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuick_Use(t *testing.T) {
	type fields struct {
		routes      []Route
		middlewares []http.HandlerFunc
		mux         *http.ServeMux
	}
	type args struct {
		middleware http.HandlerFunc
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
				routes:      tt.fields.routes,
				middlewares: tt.fields.middlewares,
				mux:         tt.fields.mux,
			}
			r.Use(tt.args.middleware)
		})
	}
}

func TestQuick_ServeHTTP(t *testing.T) {
	type fields struct {
		routes      []Route
		middlewares []http.HandlerFunc
		mux         *http.ServeMux
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
				routes:      tt.fields.routes,
				middlewares: tt.fields.middlewares,
				mux:         tt.fields.mux,
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
				JSON:     tt.fields.JSON,
				BodyByte: tt.fields.BodyByte,
				JsonStr:  tt.fields.JsonStr,
			}
			if err := c.Json(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Ctx.Json() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCtx_Byte(t *testing.T) {
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
		b []byte
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
				JSON:     tt.fields.JSON,
				BodyByte: tt.fields.BodyByte,
				JsonStr:  tt.fields.JsonStr,
			}
			if err := c.Byte(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Ctx.Byte() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCtx_SendString(t *testing.T) {
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
		s string
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
				JSON:     tt.fields.JSON,
				BodyByte: tt.fields.BodyByte,
				JsonStr:  tt.fields.JsonStr,
			}
			if err := c.SendString(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Ctx.SendString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCtx_Set(t *testing.T) {
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
		key   string
		value string
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
			c := &Ctx{
				Response: tt.fields.Response,
				Request:  tt.fields.Request,
				Headers:  tt.fields.Headers,
				Params:   tt.fields.Params,
				Query:    tt.fields.Query,
				JSON:     tt.fields.JSON,
				BodyByte: tt.fields.BodyByte,
				JsonStr:  tt.fields.JsonStr,
			}
			c.Set(tt.args.key, tt.args.value)
		})
	}
}

func TestCtx_Accepts(t *testing.T) {
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
		acceptType string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Ctx
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
				JSON:     tt.fields.JSON,
				BodyByte: tt.fields.BodyByte,
				JsonStr:  tt.fields.JsonStr,
			}
			if got := c.Accepts(tt.args.acceptType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ctx.Accepts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCtx_Status(t *testing.T) {
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
		status int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Ctx
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
				JSON:     tt.fields.JSON,
				BodyByte: tt.fields.BodyByte,
				JsonStr:  tt.fields.JsonStr,
			}
			if got := c.Status(tt.args.status); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ctx.Status() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuick_ListRoutes(t *testing.T) {
	type fields struct {
		routes      []Route
		middlewares []http.HandlerFunc
		mux         *http.ServeMux
	}
	tests := []struct {
		name   string
		fields fields
		want   []Route
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Quick{
				routes:      tt.fields.routes,
				middlewares: tt.fields.middlewares,
				mux:         tt.fields.mux,
			}
			if got := r.ListRoutes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Quick.ListRoutes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuick_Listen(t *testing.T) {
	type fields struct {
		routes      []Route
		middlewares []http.HandlerFunc
		mux         *http.ServeMux
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
			r := &Quick{
				routes:      tt.fields.routes,
				middlewares: tt.fields.middlewares,
				mux:         tt.fields.mux,
			}
			if err := r.Listen(tt.args.addr); (err != nil) != tt.wantErr {
				t.Errorf("Quick.Listen() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
