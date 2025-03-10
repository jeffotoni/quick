package quick

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// TestCtx_Bind tests if the Bind function correctly binds request data to a struct
// The will test func TestCtx_Bind(t *testing.T)
// go test -v -run ^TestCtx_Bind$
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

// TestCtx_BodyParser tests whether BodyParser correctly parses the request body
// The will test func TestCtx_BodyParser(t *testing.T)
// go test -v -run ^TestCtx_BodyParser$
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

// TestCtx_Param checks if Param retrieves the correct value for a given parameter key
// The will test func TestCtx_Param(t *testing.T)
// go test -v -run ^TestCtx_Param$
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

// TestCtx_Body verifies that Body() returns the expected request body content
// The will test func TestCtx_Body(t *testing.T)
// go test -v -run ^TestCtx_Body$
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

// TestCtx_BodyString verifies that BodyString() correctly returns the body as a string
// The will test func TestCtx_BodyString(t *testing.T)
// go test -v -run ^TestCtx_BodyString$
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

// TestCtx_Methods_JSON verifies that JSON responses are properly returned
// The will test func TestCtx_Methods_JSON(t *testing.T)
// go test -v -run ^TestCtx_Methods_JSON$
func TestCtx_Methods_JSON(t *testing.T) {

	q := New()

	q.Post("/json", func(c *Ctx) error {
		data := map[string]string{"message": "Hello, JSON!"}
		return c.JSON(data)
	})

	data, err := q.QuickTest("POST", "/json", nil)
	if err != nil {
		t.Errorf("Error during QuickTest: %v", err)
		return
	}

	if data.StatusCode() != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, data.StatusCode())
	}

	expected := `{"message":"Hello, JSON!"}`
	if !bytes.Equal(bytes.TrimSpace(data.Body()), []byte(expected)) {
		t.Errorf("Expected JSON body '%s', got '%s'", expected, data.BodyStr())
	}
}

// TestCtx_JSON verifies that the JSON function correctly encodes the response body
// The will test func TestCtx_JSON(t *testing.T)
// go test -v -run ^TestCtx_JSON$
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

// TestCtx_XML ensures that XML responses are properly returned
// The will test func TestCtx_XML(t *testing.T)
// go test -v -run ^TestCtx_XML$
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

// TestCtx_writeResponse tests the function that writes raw response bytes
// The will test func TestCtx_writeResponse(t *testing.T)
// go test -v -run ^TestCtx_writeResponse$
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

// TestCtx_Byte checks if byte responses are correctly sent
// The will test func TestCtx_Byte(t *testing.T)
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

// TestCtx_Send verifies the function that sends raw byte responses
// The will test func TestCtx_Send(t *testing.T)
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

// TestCtx_SendString verifies that the SendString function correctly sends string responses
// The will test func TestCtx_SendString(t *testing.T)
// go test -v -run ^TestCtx_SendString$
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

// TestCtx_String ensures that the String function correctly formats string responses
// The will test func TestCtx_String(t *testing.T)
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

// TestCtx_SendFile verifies that files are properly sent using SendFile
// The will test func TestCtx_SendFile(t *testing.T)
// go test -v -run ^TestCtx_SendFile$
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

// TestCtx_Set ensures that headers can be correctly set using Set
// The will test func TestCtx_Set(t *testing.T)
// go test -v -run ^TestCtx_Set$
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

// TestCtx_Append verifies that new header values can be appended correctly
// The will test func TestCtx_Append(t *testing.T)
// go test -v -run ^TestCtx_Append$
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

// TestCtx_Accepts checks whether content types are correctly validated
// The will test func TestCtx_Accepts(t *testing.T)
// go test -v -run ^TestCtx_Accepts$
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

// TestCtx_Status ensures that status codes are properly set
// The will test func TestCtx_Status(t *testing.T)
// go test -v -run ^TestCtx_Status$
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
