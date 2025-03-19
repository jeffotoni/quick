package maxbody

import (
	"fmt"

	"github.com/jeffotoni/quick"
)

const defaultMaxBytes int64 = 5 * 1024 * 1024 // 5MB

// New creates a middleware that limits the request body size.
//
// Parameters:
//   - maxBytes ...int64: (optional) Maximum allowed request body size in bytes.
//     Defaults to 5MB if not provided.
//
// Returns:
//   - A middleware function that wraps a quick.Handler to enforce body size limits.
func New(maxBytes ...int64) func(next quick.Handler) quick.Handler {
	// Determine the maximum allowed request size
	maxSize := defaultMaxBytes
	if len(maxBytes) > 0 && maxBytes[0] > 0 {
		maxSize = maxBytes[0]
	}

	return func(next quick.Handler) quick.Handler {
		return quick.HandlerFunc(func(c *quick.Ctx) error {
			// Check Content-Length before reading the body

			fmt.Println("maxSize:::", maxSize)
			fmt.Println("ContentLength:::", c.Request.ContentLength)

			if c.Request.ContentLength >= 0 && c.Request.ContentLength > maxSize {
				return c.Status(quick.StatusRequestEntityTooLarge).String("Request body too large")
			}

			// Restrict the body reader to maxSize bytes
			c.Request.Body = quick.MaxBytesReader(c.Response, c.Request.Body, maxSize+1)

			return next.ServeQuick(c)
		})
	}
}

// // New creates a middleware that limits the request body size.
// //
// // Parameters:
// //   - maxBytes ...int64: (optional) Maximum allowed request body size in bytes.
// //     Defaults to 5MB if not provided.
// //
// // Returns:
// //   - A middleware function that wraps an http.Handler to enforce body size limits.
// func New(maxBytes ...int64) func(http.Handler) http.Handler {
// 	// Determine the maximum allowed request size
// 	maxSize := defaultMaxBytes
// 	if len(maxBytes) > 0 && maxBytes[0] > 0 {
// 		maxSize = maxBytes[0]
// 	}

// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			// Check Content-Length before reading the body
// 			if r.ContentLength > maxSize {
// 				http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
// 				return
// 			}

// 			// Restrict the body reader to maxSize bytes
// 			r.Body = http.MaxBytesReader(w, r.Body, maxSize)

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }
