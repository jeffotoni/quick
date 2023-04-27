package logger

import (
	"net/http/httptest"
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

// go test -bench=. -benchtime=1s -benchmem
func BenchmarkNew(b *testing.B) {
	for n := 0; n < b.N; n++ {
		New()
	}
}
