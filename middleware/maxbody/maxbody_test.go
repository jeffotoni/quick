package maxbody

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func int2prt(x int64) *int64 {
	return &x
}

// go test -v -failfast -count=1 -run ^TestNew$
// go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestNew$; go tool cover -html=coverage.out
func TestNew(t *testing.T) {

	tests := []struct {
		name         string
		testMaxBody  testMaxBody
		maxBodyValue *int64
		wantErr      bool
	}{
		{
			name:         "success_default",
			testMaxBody:  testMaxBodySuccess,
			maxBodyValue: nil,
			wantErr:      false,
		},
		{
			name:         "success_custom",
			testMaxBody:  testMaxBodySuccess,
			maxBodyValue: int2prt(100000000),
			wantErr:      false,
		},
		{
			name:         "error_403",
			testMaxBody:  testMaxBodyFail,
			maxBodyValue: int2prt(DefaultMaxBytes),
			wantErr:      true,
		},
	}

	for _, ti := range tests {
		t.Run(ti.name, func(tt *testing.T) {
			t.Logf("==== TEST %s ====", ti.name)
			h := func(http.Handler) http.Handler { return nil }
			if ti.maxBodyValue != nil {
				h = New(*ti.maxBodyValue)
			} else {
				h = New()
			}

			a := h(ti.testMaxBody.HandlerFunc)
			rec := httptest.NewRecorder()
			a.ServeHTTP(rec, ti.testMaxBody.Request)
			resp := rec.Result()

			if resp.StatusCode != 200 && (!ti.wantErr) {
				tt.Errorf("length is not right")
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
