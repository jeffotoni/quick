// Package compress provides middleware for compressing HTTP responses using gzip.
package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/jeffotoni/quick"
)

// gzipWriterPool maintains a pool of gzip writers to optimize performance.
// Writers are reused to reduce memory allocations and improve efficiency.
var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(io.Discard)
	},
}

// gzipResponseWriter wraps an http.ResponseWriter and compresses the output.
// It overrides the Write method to pass the response through a gzip writer.
type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

// Write compresses the response data using gzip before writing it to the client.
//
// Parameters:
//   - b []byte: The data to be written.
//
// Returns:
//   - int: The number of bytes written.
//   - error: Any error encountered while writing.
func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

// clientSupportsGzip checks if the client accepts gzip encoding.
//
// Parameters:
//   - c *quick.Ctx: The request context containing headers.
//
// Returns:
//   - bool: True if the client supports gzip encoding, false otherwise.
func clientSupportsGzip(c *quick.Ctx) bool {
	return strings.Contains(strings.ToLower(c.GetHeader("Accept-Encoding")), "gzip")
}

// Gzip creates a middleware that compresses HTTP responses using gzip.
//
// This middleware checks if the client supports gzip encoding before applying compression.
// It modifies the response headers and wraps the response writer in a gzip writer.
//
// Usage:
//
//	q.Use(compress.Gzip())
//
// The middleware follows Quick's middleware format: func(next quick.Handler) quick.Handler.
//
// Returns:
//   - func(quick.Handler) quick.Handler: A middleware function for response compression.
func Gzip() func(next quick.Handler) quick.Handler {
	return func(next quick.Handler) quick.Handler {
		return quick.HandlerFunc(func(c *quick.Ctx) error {
			if !clientSupportsGzip(c) {
				// If the client does not support gzip, proceed without compression
				return next.ServeQuick(c)
			}

			// Remove Content-Length to avoid conflicts with compression
			c.Del("Content-Length")
			c.Set("Content-Encoding", "gzip")
			c.Add("Vary", "Accept-Encoding")

			// Retrieve a gzip writer from the pool
			gz := gzipWriterPool.Get().(*gzip.Writer)
			defer gzipWriterPool.Put(gz)

			gz.Reset(c.Response)
			defer gz.Close()

			// Wrap the response writer with gzipResponseWriter
			gzr := gzipResponseWriter{ResponseWriter: c.Response, writer: gz}
			c.Response = gzr

			//return next(c)
			return next.ServeQuick(c)
		})
	}
}

/// version use HandlerFunc
// func Gzip() func(next quick.HandlerFunc) quick.HandlerFunc {
// 	return func(next quick.HandlerFunc) quick.HandlerFunc {
// 		return func(c *quick.Ctx) error {
// 			if !clientSupportsGzip(c) {
// 				// If the client does not support gzip, proceed without compression
// 				return next(c)
// 			}

// 			// Remove Content-Length to avoid conflicts with compression
// 			c.Del("Content-Length")
// 			c.Set("Content-Encoding", "gzip")
// 			c.Add("Vary", "Accept-Encoding")

// 			// Retrieve a gzip writer from the pool
// 			gz := gzipWriterPool.Get().(*gzip.Writer)
// 			defer gzipWriterPool.Put(gz)

// 			gz.Reset(c.Response)
// 			defer gz.Close()

// 			// Wrap the response writer with gzipResponseWriter
// 			gzr := gzipResponseWriter{ResponseWriter: c.Response, writer: gz}
// 			c.Response = gzr

// 			return next(c)
// 		}
// 	}
// }

//// version net/http pure
// func clientSupportsGzipHttp(r *http.Request) bool {
// 	return strings.Contains(strings.ToLower(r.Header.Get("Accept-Encoding")), "gzip")
// }

// // // Gzip creates a middleware to compress the response using gzip.
// func Gzip() func(h http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			// Checks if the client supports gzip
// 			if !clientSupportsGzipHttp(r) {
// 				next.ServeHTTP(w, r)
// 				return
// 			}

// 			// Remove Content-Length para evitar conflito com resposta comprimida
// 			w.Header().Del("Content-Length")

// 			// Adjust headers to indicate that the response will be gzipped
// 			w.Header().Set("Content-Encoding", "gzip")
// 			w.Header().Add("Vary", "Accept-Encoding")

// 			gz := gzip.NewWriter(w)
// 			defer func() {
// 				if err := gz.Close(); err != nil {
// 					// If an error occurs when closing, we send 500
// 					http.Error(w, fmt.Sprintf("error closing gzip: %v", err), http.StatusInternalServerError)
// 				}
// 			}()

// 			gzr := gzipResponseWriter{writer: gz, ResponseWriter: w}
// 			next.ServeHTTP(gzr, r)
// 		})
// 	}
// }
