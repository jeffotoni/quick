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
