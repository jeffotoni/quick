package compress

import (
	"net/http/httptest"
	"testing"
)

// go test -v -failfast -count=1 -run ^TestGzip$
// go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestGzip$; go tool cover -html=coverage.out
func TestGzip(t *testing.T) {

	tests := []struct {
		name            string
		testcompress    testcompress
		contentEncoding string
		acceptEncoding  string
	}{
		{
			name:            "success",
			testcompress:    testcompressSuccess,
			contentEncoding: "gzip",
		},
		{
			name:            "success_accept_encoding",
			testcompress:    testcompressSuccess,
			contentEncoding: "gzip",
			acceptEncoding:  "gzip",
		},
	}

	for _, ti := range tests {
		t.Run(ti.name, func(tt *testing.T) {
			t.Logf("==== TEST %s ====", ti.name)
			h := Gzip()
			a := h(ti.testcompress.HandlerFunc)
			rec := httptest.NewRecorder()
			// ti.testcompress.Request.Header.Set("Content-Encoding", ti.contentEncoding)
			// ti.testcompress.Request.Header.Set("Accept-Encoding", ti.acceptEncoding)
			a.ServeHTTP(rec, ti.testcompress.Request)
			resp := rec.Result()
			if resp.Header.Get("Content-Encoding") == "" && len(ti.contentEncoding) == 0 {
				tt.Errorf("was expected a Content-Encoding and nothing came")
			}
		})
	}
}

// go test -bench=. -benchtime=1s -benchmem
func BenchmarkNew(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Gzip()
	}
}
