package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	sucessReq = httptest.NewRequest("GET", "/some", nil)
)

// go test -v -cover -count=1 -failfast -run ^Test_middleware_Auth$
func Test_middleware_Auth(t *testing.T) {

	type fields struct {
		Route http.HandlerFunc
	}
	type args struct {
		apiTokenKey string
		value       string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		Request  *http.Request
		Recorder *httptest.ResponseRecorder
		wantOut  []byte
		headers  map[string]string
	}{
		{
			name: "success",
			fields: fields{
				Route: mockMiddlewareFunc,
			},
			args: args{
				apiTokenKey: "api-token",
				value:       "myValue",
			},
			Request:  sucessReq,
			Recorder: httptest.NewRecorder(),
			wantOut:  []byte(`{"data":"mock"}`),
			headers:  map[string]string{"api-token": "myValue"},
		},
		{
			name: "fail_auth",
			fields: fields{
				Route: mockMiddlewareFunc,
			},
			args: args{
				apiTokenKey: "api-token",
				value:       "myValue",
			},
			Request:  sucessReq,
			Recorder: httptest.NewRecorder(),
			headers:  map[string]string{"api-token": "myVal"},
			wantOut:  []byte(`{"error": "api-key is missing or it isn't correct", "code": 401}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.headers {
				tt.Request.Header.Set(k, v)
			}
			m := NewMiddleware(tt.fields.Route)
			m.Auth(tt.args.apiTokenKey, tt.args.value)
			m.Route().ServeHTTP(tt.Recorder, tt.Request)
			resp := tt.Recorder.Result()
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("error: %v", err)
			}
			if string(data) != string(tt.wantOut) {
				t.Errorf("\ngot: %s\nwant: %s\n", string(data), string(tt.wantOut))
			}
		})
	}
}

func mockMiddlewareFunc(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte(`{"data":"mock"}`))
}
