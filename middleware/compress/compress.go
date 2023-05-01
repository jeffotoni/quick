package compress

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

// write in gzip and Header() from http
// this struct is to
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Gzip() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Vary", "Accept-Encoding") // Add Vary header
			gz := gzip.NewWriter(w)

			defer func() {
				err := gz.Close()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "error closing gzip: %+v\n", err)
					return
				}
			}()
			gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
			h.ServeHTTP(gzr, r)
		})
	}
}
