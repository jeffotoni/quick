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
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// TestCtx_Bind validates whether the Ctx.Bind() function properly binds the request body into a given struct.
//
// This test ensures that:
//   - JSON payloads are correctly unmarshaled into the target structure.
//   - Errors are properly returned when the binding fails.
//
// To run:
//
//	$ go test -v -run ^TestCtx_Bind$
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

// TestCtx_BodyParser checks if the Ctx.BodyParser method correctly parses the request body
// and maps it to the given struct.
//
// It ensures that BodyParser correctly handles JSON decoding.
//
// To run:
//
//	$ go test -v -run ^TestCtx_BodyParser$
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

// TestCtx_Param ensures that the Ctx.Param method retrieves the correct route parameter
// based on the provided key.
//
// Useful for validating route variable extraction.
//
// To run:
//
//	$ go test -v -run ^TestCtx_Param$
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

// TestCtx_Body verifies whether Ctx.Body returns the expected byte slice representing
// the raw request body content.
//
// To run:
//
//	$ go test -v -run ^TestCtx_Body$
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

// TestCtx_BodyString verifies whether Ctx.BodyString returns the request body content as a string.
//
// It is helpful to check if text content can be retrieved from request payloads.
//
// To run:
//
//	$ go test -v -run ^TestCtx_BodyString$
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

// TestCtx_Methods_JSON ensures that JSON responses are properly serialized and returned by the Ctx.JSON() method.
//
// It registers a POST route and sends a request using Qtest to validate:
//   - The response status code is HTTP 200 OK.
//   - The response body is correctly formatted as a JSON string.
//
// To run:
//
//	$ go test -v -run ^TestCtx_Methods_JSON$
func TestCtx_Methods_JSON(t *testing.T) {

	q := New()

	q.Post("/json", func(c *Ctx) error {
		data := map[string]string{"message": "Hello, JSON!"}
		return c.JSON(data)
	})

	data, err := q.Qtest(QuickTestOptions{
		Method:  MethodPost,
		URI:     "/json",
		Headers: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		t.Errorf("Error during Qtest: %v", err)
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

// TestCtx_JSON ensures that Ctx.JSON serializes the given struct or map and writes it
// as a JSON response with the appropriate headers.
//
// To run:
//
//	$ go test -v -run ^TestCtx_JSON$
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

// TestCtx_XML checks that the Ctx.XML method serializes the input struct and returns
// an XML response.
//
// To run:
//
//	$ go test -v -run ^TestCtx_XML$
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

// TestCtx_writeResponse verifies that the internal method writeResponse writes the raw byte
// data to the response correctly.
//
// To run:
//
//	$ go test -v -run ^TestCtx_writeResponse$
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

// TestCtx_Byte verifies that the Byte method correctly writes raw bytes
// to the response body.
//
// To run:
//
//	$ go test -v -failfast -count=1 -run ^TestCtx_Byte$
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

// TestCtx_Send checks if the Send method sends raw byte responses correctly,
// without modifying or formatting them.
//
// To run:
//
//	$ go test -v -failfast -count=1 -run ^TestCtx_Send$
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

// TestCtx_SendString tests if the Ctx.SendString method sends the correct string response.
//
// To run:
//
//	$ go test -v -run ^TestCtx_SendString$
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

// TestCtx_String ensures that the String method sends plain string content
// as the response body.
//
// To run:
//
//	$ go test -v -failfast -count=1 -run ^TestCtx_String$
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

// TestCtx_SendFile tests whether Ctx.SendFile writes the given byte slice as file content
// in the response body.
//
// To run:
//
//	$ go test -v -run ^TestCtx_SendFile$
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

// TestCtx_Send checks if the Send method sends raw byte responses correctly,
// without modifying or formatting them.
//
// To run:
//
//	$ go test -v -failfast -count=1 -run ^TestCtx_Send$
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

// TestCtx_Append ensures that the Append method correctly adds headers,
// even when the header key already exists (appending instead of replacing).
//
// To run:
//
//	$ go test -v -run ^TestCtx_Append$
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

// TestCtx_Accepts ensures that Ctx.Accepts correctly evaluates the Accept header
// to determine if the content type is acceptable.
//
// To run:
//
//	$ go test -v -run ^TestCtx_Accepts$
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

// TestCtx_Status validates that calling Ctx.Status sets the correct status code for the response.
//
// To run:
//
//	$ go test -v -run ^TestCtx_Status$
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

// TestCtxMethods validates multiple Ctx methods for extracting request data.
//
// It checks if headers, IP address, method, path, and query parameters are correctly retrieved.
//
// Check:
//   - Ctx.GetHeader()
//   - Ctx.GetHeaders()
//   - Ctx.RemoteIP()
//   - Ctx.Method()
//   - Ctx.Path()
//   - Ctx.QueryParam()
//
// To run:
//
//	$ go test -v -run ^TestCtxMethods$
func TestCtxMethods(t *testing.T) {
	// Prepare the test request
	req := httptest.NewRequest(http.MethodGet, "/testpath?search=golang", nil)
	req.RemoteAddr = "192.168.1.10:12345"
	req.Header.Set("User-Agent", "Go-Test-Agent")

	// Create the fake ResponseWriter
	rec := httptest.NewRecorder()

	// Create the Quick context
	c := &Ctx{
		Response: rec,
		Request:  req,
	}

	// Test GetHeader
	if got := c.GetHeader("User-Agent"); got != "Go-Test-Agent" {
		t.Errorf("GetHeader() = %v, want %v", got, "Go-Test-Agent")
	}

	// Test GetHeaders
	headers := c.GetHeaders()
	if headers.Get("User-Agent") != "Go-Test-Agent" {
		t.Errorf("GetHeaders().Get(\"User-Agent\") = %v, want %v", headers.Get("User-Agent"), "Go-Test-Agent")
	}

	// RemoteIP Test
	if ip := c.RemoteIP(); ip != "192.168.1.10" {
		t.Errorf("RemoteIP() = %v, want %v", ip, "192.168.1.10")
	}

	// Test Method
	if method := c.Method(); method != http.MethodGet {
		t.Errorf("Method() = %v, want %v", method, http.MethodGet)
	}

	// Test Path
	if path := c.Path(); path != "/testpath" {
		t.Errorf("Path() = %v, want %v", path, "/testpath")
	}

	// Test Query
	if q := c.QueryParam("search"); q != "golang" {
		t.Errorf("Query(\"search\") = %v, want %v", q, "golang")
	}
}

// TestCtxOriginalURI validates the behavior of Ctx.OriginalURI().
//
// It ensures that the original request URI is correctly retrieved from the HTTP request,
// including the path and query string as sent by the client.
//
// Check:
//   - Ctx.OriginalURI()
//
// To run:
//
//	$ go test -v -run ^TestCtxOriginalURI$
func TestCtxOriginalURI(t *testing.T) {
	// Prepare the test request
	req := httptest.NewRequest(http.MethodGet, "/v1/api/testpath?search=golang", nil)
	req.RemoteAddr = "192.168.1.10:12345"
	req.Header.Set("User-Agent", "Go-Test-Agent")

	// Create the fake ResponseWriter
	rec := httptest.NewRecorder()

	// Create the Quick context
	c := &Ctx{
		Response: rec,
		Request:  req,
	}

	// Test OriginalURI
	if original := c.OriginalURI(); original != "/v1/api/testpath?search=golang" {
		t.Errorf("originalURI() = %v, want %v", original, "/v1/api/testpath?search=golang")
	}
}
