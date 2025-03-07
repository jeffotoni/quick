package compress

import (
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// failOnSecondWriteWriter fails on the second attempt to write, simulating an error in the gzip flush.
type failOnSecondWriteWriter struct {
	writeCount int
	underlying io.Writer
}

func (f *failOnSecondWriteWriter) Write(p []byte) (int, error) {
	f.writeCount++
	// The first write is normal
	if f.writeCount == 2 {
		// Second write fails
		return 0, errors.New("forced error on close")
	}
	return f.underlying.Write(p)
}

// customResponseRecorder injects a failOnSecondWriteWriter in place of the ResponseRecorder's Writer.
func TestGzipMiddleware(t *testing.T) {
	t.Run("No Accept-Encoding => no compression", func(t *testing.T) {
		handler := Gzip()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		// We wait for the response without compression
		if enc := w.Header().Get("Content-Encoding"); enc != "" {
			t.Errorf("Expected no Content-Encoding, got '%s'", enc)
		}
		if !strings.Contains(w.Body.String(), "Hello, World!") {
			t.Errorf("Response body mismatch, got '%s'", w.Body.String())
		}
	})

	t.Run("Accept-Encoding => compressed with gzip", func(t *testing.T) {
		handler := Gzip()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, Gopher!"))
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept-Encoding", "gzip")

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		// Check if it is marked as gzip
		if enc := w.Header().Get("Content-Encoding"); enc != "gzip" {
			t.Errorf("Expected Content-Encoding=gzip, got '%s'", enc)
		}
		if !strings.Contains(w.Header().Get("Vary"), "Accept-Encoding") {
			t.Errorf("Expected Vary=Accept-Encoding, got '%s'", w.Header().Get("Vary"))
		}

		// Unzip to confirm content
		gzReader, err := gzip.NewReader(w.Body)
		if err != nil {
			t.Fatalf("Failed to create gzip reader: %v", err)
		}
		defer gzReader.Close()

		unzipped, err := io.ReadAll(gzReader)
		if err != nil {
			t.Fatalf("Failed to read unzipped content: %v", err)
		}
		if string(unzipped) != "Hello, Gopher!" {
			t.Errorf("Unzipped content mismatch, got '%s'", string(unzipped))
		}
	})
}

type customResponseRecorder struct {
	*httptest.ResponseRecorder
	failWriter *failOnSecondWriteWriter
}

func (c *customResponseRecorder) Write(p []byte) (int, error) {
	return c.failWriter.Write(p)
}

func (c *customResponseRecorder) WriteHeader(code int) {
	c.ResponseRecorder.WriteHeader(code)
}

func (c *customResponseRecorder) Header() http.Header {
	return c.ResponseRecorder.Header()
}
