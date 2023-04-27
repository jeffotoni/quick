package logger

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// go test -v -failfast -count=1 -run ^TestNew$
// go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestNew$; go tool cover -html=coverage.out
func TestNew(t *testing.T) {

	// type args struct {
	// 	Config Config
	// }

	tests := []struct {
		name string
		// args       args
		testLogger testLogger
		wantErr    bool
	}{
		{
			name:       "success",
			testLogger: testLoggerSuccess,
			wantErr:    false,
		},
		{
			name:       "success_with_body",
			testLogger: testLoggerSuccessBody,
			wantErr:    false,
		},
	}

	for _, ti := range tests {
		t.Run(ti.name, func(tt *testing.T) {
			t.Logf("==== TEST %s ====", ti.name)
			h := New()
			a := h(ti.testLogger.HandlerFunc)
			rec := httptest.NewRecorder()
			a.ServeHTTP(rec, ti.testLogger.Request)
			resp := rec.Result()

			if resp.StatusCode != 200 && (!ti.wantErr) {
				tt.Errorf("there was suppose to have %v and got %v", 200, resp.StatusCode)
			}
		})
	}
}

func TestLoggerMw(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Quick in action!"))
	})

	mw := New()
	hmw := mw(handler)
	var req = &http.Request{
		Header:     http.Header{},
		Host:       "localhost:3000",
		Method:     "GET",
		RemoteAddr: "127.0.0.1:3000",
		URL: &url.URL{
			Scheme: "http",
			Host:   "quick.com",
		},
	}
	rr := httptest.NewRecorder()
	hmw.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned an incorrect status: expected %v, obtained %v", http.StatusOK, status)
	}

	expected := "Quick in action!"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned an incorrect status: expected %v, obtained %v", expected, rr.Body.String())
	}
}

func TestLoggerMw500(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Quick in action 20223!"))
	})

	middleware := New()
	handlerWithMiddleware := middleware(handler)

	// req, err := http.NewRequest("POST", "http://localhost:3000", bytes.NewBuffer([]byte("request body")))
	// if err != nil {
	// 	t.Fatalf("Erro ao criar a requisição: %v", err)
	// }

	var req = &http.Request{
		Header:     http.Header{},
		Host:       "localhost:3000",
		RemoteAddr: "127.0.0.1",
		Method:     "POST",
		URL: &url.URL{
			Scheme: "http",
			Host:   "letsgoquick.com",
		},
	}

	// Defina RemoteAddr com um valor inválido para causar um erro interno no servidor
	req.RemoteAddr = "invalid"

	rr := httptest.NewRecorder()
	handlerWithMiddleware.ServeHTTP(rr, req)

	// Verifique se o status da resposta está correto
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned an incorrect status: expected %v, obtained %v", http.StatusInternalServerError, status)
	}
}

// go test -bench=. -benchtime=1s -benchmem
func BenchmarkNew(b *testing.B) {
	for n := 0; n < b.N; n++ {
		New()
	}
}

// func BenchmarkNew2(b *testing.B) {
// 	for n := 0; n < b.N; n++ {
// 		New2()
// 	}
// }
