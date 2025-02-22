package goquick

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCtx_Bind(t *testing.T) {
	type fields struct {
		Response  http.ResponseWriter
		Request   *http.Request
		resStatus int
		bodyByte  []byte
		JsonStr   string
		Headers   map[string][]string
		Params    map[string]string
		Query     map[string]string
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
				Response:  tt.fields.Response,
				Request:   tt.fields.Request,
				resStatus: tt.fields.resStatus,
				bodyByte:  tt.fields.bodyByte,
				JsonStr:   tt.fields.JsonStr,
				Headers:   tt.fields.Headers,
				Params:    tt.fields.Params,
				Query:     tt.fields.Query,
			}
			if err := c.Bind(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Ctx.Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCtx_BodyParser(t *testing.T) {
	type fields struct {
		Response  http.ResponseWriter
		Request   *http.Request
		resStatus int
		bodyByte  []byte
		JsonStr   string
		Headers   map[string][]string
		Params    map[string]string
		Query     map[string]string
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
				Response:  tt.fields.Response,
				Request:   tt.fields.Request,
				resStatus: tt.fields.resStatus,
				bodyByte:  tt.fields.bodyByte,
				JsonStr:   tt.fields.JsonStr,
				Headers:   tt.fields.Headers,
				Params:    tt.fields.Params,
				Query:     tt.fields.Query,
			}
			if err := c.BodyParser(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Ctx.BodyParser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCtx_Param(t *testing.T) {
	type fields struct {
		Response  http.ResponseWriter
		Request   *http.Request
		resStatus int
		bodyByte  []byte
		JsonStr   string
		Headers   map[string][]string
		Params    map[string]string
		Query     map[string]string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Ctx{
				Response:  tt.fields.Response,
				Request:   tt.fields.Request,
				resStatus: tt.fields.resStatus,
				bodyByte:  tt.fields.bodyByte,
				JsonStr:   tt.fields.JsonStr,
				Headers:   tt.fields.Headers,
				Params:    tt.fields.Params,
				Query:     tt.fields.Query,
			}
			if got := c.Param(tt.args.key); got != tt.want {
				t.Errorf("Ctx.Param() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCtx_Body(t *testing.T) {
	type fields struct {
		Response  http.ResponseWriter
		Request   *http.Request
		resStatus int
		bodyByte  []byte
		JsonStr   string
		Headers   map[string][]string
		Params    map[string]string
		Query     map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Ctx{
				Response:  tt.fields.Response,
				Request:   tt.fields.Request,
				resStatus: tt.fields.resStatus,
				bodyByte:  tt.fields.bodyByte,
				JsonStr:   tt.fields.JsonStr,
				Headers:   tt.fields.Headers,
				Params:    tt.fields.Params,
				Query:     tt.fields.Query,
			}
			if got := c.Body(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ctx.Body() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCtx_BodyString(t *testing.T) {
	type fields struct {
		Response  http.ResponseWriter
		Request   *http.Request
		resStatus int
		bodyByte  []byte
		JsonStr   string
		Headers   map[string][]string
		Params    map[string]string
		Query     map[string]string
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
				Response:  tt.fields.Response,
				Request:   tt.fields.Request,
				resStatus: tt.fields.resStatus,
				bodyByte:  tt.fields.bodyByte,
				JsonStr:   tt.fields.JsonStr,
				Headers:   tt.fields.Headers,
				Params:    tt.fields.Params,
				Query:     tt.fields.Query,
			}
			if got := c.BodyString(); got != tt.want {
				t.Errorf("Ctx.BodyString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCtx_JSON(t *testing.T) {
	type fields struct {
		Response  http.ResponseWriter
		Request   *http.Request
		resStatus int
		bodyByte  []byte
		JsonStr   string
		Headers   map[string][]string
		Params    map[string]string
		Query     map[string]string
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
				Response:  tt.fields.Response,
				Request:   tt.fields.Request,
				resStatus: tt.fields.resStatus,
				bodyByte:  tt.fields.bodyByte,
				JsonStr:   tt.fields.JsonStr,
				Headers:   tt.fields.Headers,
				Params:    tt.fields.Params,
				Query:     tt.fields.Query,
			}
			if err := c.JSON(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Ctx.JSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCtx_XML(t *testing.T) {
	type fields struct {
		Response  http.ResponseWriter
		Request   *http.Request
		resStatus int
		bodyByte  []byte
		JsonStr   string
		Headers   map[string][]string
		Params    map[string]string
		Query     map[string]string
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
				Response:  tt.fields.Response,
				Request:   tt.fields.Request,
				resStatus: tt.fields.resStatus,
				bodyByte:  tt.fields.bodyByte,
				JsonStr:   tt.fields.JsonStr,
				Headers:   tt.fields.Headers,
				Params:    tt.fields.Params,
				Query:     tt.fields.Query,
			}
			if err := c.XML(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Ctx.XML() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCtx_writeResponse(t *testing.T) {
	type fields struct {
		Response  http.ResponseWriter
		Request   *http.Request
		resStatus int
		bodyByte  []byte
		JsonStr   string
		Headers   map[string][]string
		Params    map[string]string
		Query     map[string]string
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
				Response:  tt.fields.Response,
				Request:   tt.fields.Request,
				resStatus: tt.fields.resStatus,
				bodyByte:  tt.fields.bodyByte,
				JsonStr:   tt.fields.JsonStr,
				Headers:   tt.fields.Headers,
				Params:    tt.fields.Params,
				Query:     tt.fields.Query,
			}
			if err := c.writeResponse(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Ctx.writeResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// go test -v -failfast -count=1 -run ^TestCtx_Byte$
func TestCtx_Byte(t *testing.T) {
	type args struct {
		response string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success_string",
			args: args{
				response: `"data": "gopher"`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := httptest.NewRecorder()

			c := &Ctx{
				Response: x,
			}

			if err := c.Byte([]byte(tt.args.response)); (err != nil) != tt.wantErr {
				t.Errorf("Ctx.String() error = %v, wantErr %v", err, tt.wantErr)
			}

			result := x.Result()
			if result.Body != nil {
				defer result.Body.Close()
				b, err := io.ReadAll(result.Body)
				if err != nil {
					t.Errorf("error: %v", err)
				}

				if string(b) != tt.args.response {
					t.Errorf("was suppose to have header value: %s and got %s", tt.args.response, string(b))
				}
			}
		})
	}
}

// go test -v -failfast -count=1 -run ^TestCtx_Send$
func TestCtx_Send(t *testing.T) {
	type args struct {
		response string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success_string",
			args: args{
				response: `"data": "jeffotoni send dados all"`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := httptest.NewRecorder()

			c := &Ctx{
				Response: x,
			}

			if err := c.Send([]byte(tt.args.response)); (err != nil) != tt.wantErr {
				t.Errorf("Ctx.String() error = %v, wantErr %v", err, tt.wantErr)
			}

			result := x.Result()
			if result.Body != nil {
				defer result.Body.Close()
				b, err := io.ReadAll(result.Body)
				if err != nil {
					t.Errorf("error: %v", err)
				}

				if string(b) != tt.args.response {
					t.Errorf("was suppose to have header value: %s and got %s", tt.args.response, string(b))
				}
			}
		})
	}
}

func TestCtx_SendString(t *testing.T) {
	type fields struct {
		Response  http.ResponseWriter
		Request   *http.Request
		resStatus int
		bodyByte  []byte
		JsonStr   string
		Headers   map[string][]string
		Params    map[string]string
		Query     map[string]string
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
				Response:  tt.fields.Response,
				Request:   tt.fields.Request,
				resStatus: tt.fields.resStatus,
				bodyByte:  tt.fields.bodyByte,
				JsonStr:   tt.fields.JsonStr,
				Headers:   tt.fields.Headers,
				Params:    tt.fields.Params,
				Query:     tt.fields.Query,
			}
			if err := c.SendString(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Ctx.SendString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// go test -v -failfast -count=1 -run ^TestCtx_String$
func TestCtx_String(t *testing.T) {
	type args struct {
		response string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success_string",
			args: args{
				response: `"data": "gopher"`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := httptest.NewRecorder()

			c := &Ctx{
				Response: x,
			}

			if err := c.String(tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("Ctx.String() error = %v, wantErr %v", err, tt.wantErr)
			}

			result := x.Result()
			if result.Body != nil {
				defer result.Body.Close()
				b, err := io.ReadAll(result.Body)
				if err != nil {
					t.Errorf("error: %v", err)
				}

				if string(b) != tt.args.response {
					t.Errorf("was suppose to have header value: %s and got %s", tt.args.response, string(b))
				}
			}
		})
	}
}

func TestCtx_SendFile(t *testing.T) {
	type fields struct {
		Response  http.ResponseWriter
		Request   *http.Request
		resStatus int
		bodyByte  []byte
		JsonStr   string
		Headers   map[string][]string
		Params    map[string]string
		Query     map[string]string
	}
	type args struct {
		file []byte
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
				Response:  tt.fields.Response,
				Request:   tt.fields.Request,
				resStatus: tt.fields.resStatus,
				bodyByte:  tt.fields.bodyByte,
				JsonStr:   tt.fields.JsonStr,
				Headers:   tt.fields.Headers,
				Params:    tt.fields.Params,
				Query:     tt.fields.Query,
			}
			if err := c.SendFile(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("Ctx.SendFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// go test -v -failfast -count=1 -run ^TestCtx_Set$
func TestCtx_Set(t *testing.T) {
	type fields struct {
		Response http.ResponseWriter
	}

	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantHeaderValue string
		wantError       bool
	}{
		{
			name: "success_Set_Headers",
			fields: fields{
				Response: func() http.ResponseWriter { x := httptest.NewRecorder(); return x }(),
			},
			args: args{
				key:   "my-key",
				value: "my-header-value",
			},
			wantHeaderValue: "my-header-value",
			wantError:       false,
		},
		{
			name: "wrong_header_check",
			fields: fields{
				Response: func() http.ResponseWriter { x := httptest.NewRecorder(); return x }(),
			},
			args: args{
				key:   "my-key",
				value: "my-header-valuee",
			},
			wantHeaderValue: "my-header-value",
			wantError:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Ctx{
				Response: tt.fields.Response,
			}

			c.Set(tt.args.key, tt.args.value)

			got := c.Response.Header().Get(tt.args.key)

			if (!tt.wantError) && got != tt.wantHeaderValue {
				t.Errorf("was suppose to have header value: %s and got %s", tt.wantHeaderValue, got)
			}
		})
	}
}

func TestCtx_Append(t *testing.T) {
	type fields struct {
		Response  http.ResponseWriter
		Request   *http.Request
		resStatus int
		bodyByte  []byte
		JsonStr   string
		Headers   map[string][]string
		Params    map[string]string
		Query     map[string]string
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantLen int
	}{
		{
			name: "should be able to create a new header",
			fields: fields{
				Response: httptest.NewRecorder(),
			},
			args: args{
				key:   "Append",
				value: "one",
			},
			wantLen: 1,
		},
		{
			name: "should be able to append to existing header",
			fields: fields{
				Response: func() http.ResponseWriter { x := httptest.NewRecorder(); x.Header().Set("Append", "one"); return x }(),
			},
			args: args{
				key:   "Append",
				value: "two",
			},
			wantLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Ctx{
				Response:  tt.fields.Response,
				Request:   tt.fields.Request,
				resStatus: tt.fields.resStatus,
				bodyByte:  tt.fields.bodyByte,
				JsonStr:   tt.fields.JsonStr,
				Headers:   tt.fields.Headers,
				Params:    tt.fields.Params,
				Query:     tt.fields.Query,
			}
			c.Append(tt.args.key, tt.args.value)

			if len(c.Response.Header().Values(tt.args.key)) != tt.wantLen {
				t.Errorf("c.Append(): want %v, got %v", tt.wantLen, len(c.Response.Header().Values(tt.args.key)))
			}
		})
	}
}

func TestCtx_Accepts(t *testing.T) {
	type fields struct {
		Response  http.ResponseWriter
		Request   *http.Request
		resStatus int
		bodyByte  []byte
		JsonStr   string
		Headers   map[string][]string
		Params    map[string]string
		Query     map[string]string
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
				Response:  tt.fields.Response,
				Request:   tt.fields.Request,
				resStatus: tt.fields.resStatus,
				bodyByte:  tt.fields.bodyByte,
				JsonStr:   tt.fields.JsonStr,
				Headers:   tt.fields.Headers,
				Params:    tt.fields.Params,
				Query:     tt.fields.Query,
			}
			if got := c.Accepts(tt.args.acceptType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ctx.Accepts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCtx_Status(t *testing.T) {
	type fields struct {
		Response  http.ResponseWriter
		Request   *http.Request
		resStatus int
		bodyByte  []byte
		JsonStr   string
		Headers   map[string][]string
		Params    map[string]string
		Query     map[string]string
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
				Response:  tt.fields.Response,
				Request:   tt.fields.Request,
				resStatus: tt.fields.resStatus,
				bodyByte:  tt.fields.bodyByte,
				JsonStr:   tt.fields.JsonStr,
				Headers:   tt.fields.Headers,
				Params:    tt.fields.Params,
				Query:     tt.fields.Query,
			}
			if got := c.Status(tt.args.status); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ctx.Status() = %v, want %v", got, tt.want)
			}
		})
	}
}
