package compress

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// write in gzip and Header() from http
// this struct is to
//type gzipResponseWriter struct {
//	io.Writer
//	http.ResponseWriter
//}
//
//func (w gzipResponseWriter) Write(b []byte) (int, error) {
//	return w.Writer.Write(b)
//}
//
//func Gzip() func(h http.Handler) http.Handler {
//	return func(h http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			w.Header().Set("Content-Encoding", "gzip")
//			w.Header().Set("Vary", "Accept-Encoding") // Add Vary header
//			gz := gzip.NewWriter(w)
//
//			defer func() {
//				err := gz.Close()
//				if err != nil {
//					w.WriteHeader(http.StatusInternalServerError)
//					fmt.Fprintf(w, "error closing gzip: %+v\n", err)
//					return
//				}
//			}()
//			gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
//			h.ServeHTTP(gzr, r)
//		})
//	}
//}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// We ensure that writing is directed to gzip.Writer
func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Gzip creates a middleware to compress the response using gzip.
func Gzip() func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Checks if the client supports gzip
			if !clientSupportsGzip(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Remove Content-Length para evitar conflito com resposta comprimida
			w.Header().Del("Content-Length")

			// Adjust headers to indicate that the response will be gzipped
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Add("Vary", "Accept-Encoding")

			gz := gzip.NewWriter(w)
			defer func() {
				if err := gz.Close(); err != nil {
					// If an error occurs when closing, we send 500
					http.Error(w, fmt.Sprintf("error closing gzip: %v", err), http.StatusInternalServerError)
				}
			}()

			gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
			next.ServeHTTP(gzr, r)
		})
	}
}

// clientSupportsGzip checks if the client sent 'Accept-Encoding: gzip'
func clientSupportsGzip(r *http.Request) bool {
	encoding := r.Header.Get("Accept-Encoding")
	return strings.Contains(strings.ToLower(encoding), "gzip")
}
