# 📂 Middleware - Quick Framework ![Quick Logo](/quick.png)

The **`middleware`** directory contains useful middleware implementations for the Quick Framework, making it easy to integrate common features such as authentication, compression, logging, request size limits, and UUID tracking.

📌 **`What are Middlewares?`**

Middlewares are functions that intercept HTTP requests before they reach the final handler. They allow:

- Validation (e.g., authentication and security policies)
- Request/response modification (e.g., GZIP compression)
- Logging and monitoring (e.g., request logging and UUID tracking)

---

### 📜 Middlewares Available

🔐 BasicAuth
Provides HTTP Basic Authentication, requiring a username and password to access protected routes.

- Can be applied globally or to specific routes.
- Supports authentication via environment variables.
- Allows custom implementation if needed.

---

### 📦 Compress
The Compress Middleware in Quick provides automatic GZIP compression for HTTP responses, reducing response sizes and improving performance.

#### 📌 Features:

- ✅ Automatic compression detection (Accept-Encoding: gzip)
- ✅ Transparent response compression without modifying business logic
- ✅ Bandwidth efficiency improvement
#### 🚀 Middleware Implementations in Quick
The Quick framework supports multiple styles of middleware implementation for GZIP compression:

#### 1️⃣ Native Quick Implementation (quick.Handler)
This is the standard and recommended approach in Quick. It follows the framework’s native middleware pattern using quick.Handler.
```go
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
```
📌 ✅ This is the default implementation used in Quick.

---
#### 2️⃣ Alternative Syntax: quick.HandlerFunc
Quick also supports quick.HandlerFunc, which allows a function-based syntax instead of using the Handler interface.

```go
func Gzip() func(next quick.HandlerFunc) quick.HandlerFunc {
	return func(next quick.HandlerFunc) quick.HandlerFunc {
		return func(c *quick.Ctx) error {
			if !clientSupportsGzip(c) {
				// If the client does not support gzip, proceed without compression
				return next(c)
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

			return next(c)
		}
	}
}
```
#### 📌 ✅ This is functionally equivalent to the default version but follows a different syntax.
---

#### 3️⃣ Pure net/http Middleware Support
Quick also supports native `net/http` middleware, making it compatible with standard Go HTTP handlers.

```go
func clientSupportsGzipHttp(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.Header.Get("Accept-Encoding")), "gzip")
}

// Gzip creates a middleware to compress the response using gzip.
func Gzip() func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Checks if the client supports gzip
			if !clientSupportsGzipHttp(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Remove Content-Length to avoid conflict with compressed response
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

			gzr := gzipResponseWriter{writer: gz, ResponseWriter: w}
			next.ServeHTTP(gzr, r)
		})
	}
}
```
📌 ✅ Useful when integrating Quick with net/http-based applications.

### Middleware Implementation

#### 🔹 Using Quick Default Middleware (quick.Handler)
```go
package main

import (
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/compress"
)

func main() {
	q := quick.New()

	// Enable GZIP compression
	q.Use(compress.Gzip())

	// Define a compressed response route
	q.Get("/v1/compress", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		c.Set("Accept-Encoding", "gzip")

		type response struct {
			Msg     string              `json:"msg"`
			Headers map[string][]string `json:"headers"`
		}

		return c.Status(200).JSON(&response{
			Msg:     "Quick ❤️",
			Headers: c.Headers,
		})
	})

	log.Fatal(q.Listen("0.0.0.0:8080"))
}

```
#### 🔹 Using Quick HandlerFunc Middleware (quick.HandlerFunc)
```go
package main

import (
	"log"
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/compress"
)

func main() {
	q := quick.New()

	// Enable GZIP middleware using HandlerFunc version
	q.Use(compress.Gzip())

	// Define a compressed response route
	q.Get("/v1/compress", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		type response struct {
			Msg     string              `json:"msg"`
			Headers map[string][]string `json:"headers"`
		}

		return c.Status(200).JSON(&response{
			Msg:     "Quick ❤️",
			Headers: c.Headers,
		})
	})

	log.Fatal(q.Listen(":8080"))
}
```
#### 🔹 Using Pure net/http Middleware
```go
package main

import (
	"log"
	"net/http"
	"github.com/jeffotoni/quick/middleware/compress"
)

func main() {
	mux := http.NewServeMux()

	// Route with compression enabled using the middleware
	mux.Handle("/v1/compress", compress.Gzip()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Hello, net/http with Gzip!"}`))
	})))

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```
---

#### 🌐 CORS (Cross-Origin Resource Sharing)
Controls how your API can be accessed from different domains.

- Restricts which domains, methods, and headers are allowed.
- Helps prevent CORS errors in browsers.
- Configurable via allowed origins, headers, and credentials.

---

#### 📜 Logger (Request Logging)
Logs incoming HTTP requests, helping in monitoring and debugging.

- Logs request method, path, response time, and status code.
- Can be integrated with structured logging tools.
- Helps with API usage tracking and debugging.

---

#### 📏 Maxbody (Request Size Limiter)
Restricts the maximum request body size to prevent clients from sending excessively large payloads.

- Avoids excessive memory usage.
- Can prevent attacks such as DoS (Denial-of-Service).
- Returns a 413 Payload Too Large error when exceeded.

---

#### 🔄 MsgUUID
Assigns a UUID (Universally Unique Identifier) to each request.

- Allows easy tracking of requests in logs.
- Useful for distributed systems where tracing requests across services is required.
- Adds a unique identifier to every request automatically.

---

### 🚧 **Coming soon!**
- Etag
- Pprof
- Proxy
- RequestID
- Skip
- Timeout

