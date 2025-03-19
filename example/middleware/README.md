# 📂 Middleware - Quick Framework ![Quick Logo](/quick.png)

The **`middleware`** directory contains useful middleware implementations for the Quick Framework, making it easy to integrate common features such as authentication, compression, logging, request size limits, and UUID tracking.

📌 **`What are Middlewares?`**

Middlewares are functions that intercept HTTP requests before they reach the final handler. They allow:

- Validation (e.g., authentication and security policies)
- Request/response modification (e.g., GZIP compression)
- Logging and monitoring (e.g., request logging and UUID tracking)

---

## 📜 Middlewares Available

🔐 BasicAuth
Provides HTTP Basic Authentication, requiring a username and password to access protected routes.

- Can be applied globally or to specific routes.
- Supports authentication via environment variables.
- Allows custom implementation if needed.

---

## 📦 Compress
Enables automatic GZIP compression of HTTP responses to reduce response size and improve performance.

- Detects if the client supports compression (Accept-Encoding: gzip).
- Compresses responses transparently without modifying business logic.
- Improves bandwidth efficiency.
---

## 🌐 CORS (Cross-Origin Resource Sharing)
Controls how your API can be accessed from different domains.

- Restricts which domains, methods, and headers are allowed.
- Helps prevent CORS errors in browsers.
- Configurable via allowed origins, headers, and credentials.

---

## 📜 Logger (Request Logging)
The `logger` middleware captures HTTP request details, helping with monitoring, debugging, and analytics.

#### 🚀 Key Features:
- ✅ Logs request method, path, response time, and status code.
- ✅ Supports multiple formats: text, json, and slog (structured logging).
- ✅ Helps track API usage and debugging.
- ✅ Customizable log patterns and additional fields.

#### 📝 Default Logging (Text Format)
This example applies logging in text format with custom log fields.

```go
package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/logger"
)

func main() {

	q := quick.New()

	// Apply the logger middleware with custom configuration
	q.Use(logger.New(logger.Config{
		Format:  "text", // Available formats: "text", "json", "slog"
		Pattern: "[${level}] ${ip} ${method} - ${latency} user_id=${user_id} trace=${trace}\n",
		Level:   "DEBUG", // Logging level: "DEBUG", "INFO", "WARN", "ERROR"
		CustomFields: map[string]string{ // Custom fields included in logs
			"user_id": "12345",
			"trace":   "xyz",
		},
	}))

	// Define a route that logs request details
	q.Get("/v1/logger", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// Return a JSON response
		return c.Status(200).JSON(quick.M{
			"msg": "Quick ❤️",
		})
	})

	// Start the server
	q.Listen("0.0.0.0:8080")
}
```
---
### 🛠️ Structured Logging (Slog Format)

This example uses structured logging (slog) for better log parsing.

```go
package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/logger"
)

func main() {

	q := quick.New()

	// Apply logger middleware with structured logging (slog)
	q.Use(logger.New(logger.Config{
		Format: "slog",
		Level:  "DEBUG",
		Pattern: "[${level}] ${ip} ${method} ${path} - ${latency} " +
			"user=${user_id} trace=${trace}\n",
		CustomFields: map[string]string{
			"user_id": "99999",
			"trace":   "abcdef",
		},
	}))

	// Define a route with structured logging
	q.Get("/v1/logger/slog", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		return c.Status(200).JSON(quick.M{
			"msg": "Structured logging with slog",
		})
	})

	// Start the server
	q.Listen("0.0.0.0:8080")
}
```
---
### 📦 JSON Logging (Machine-Readable)

Ideal for log aggregation systems, this example logs in JSON format.

```go
package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/logger"
)

func main() {

	q := quick.New()

	// Apply logger with JSON format for structured logging
	q.Use(logger.New(logger.Config{
		Format: "json",
		Level:  "INFO",
	}))

	// Define a logging route
	q.Get("/v1/logger/json", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		return c.Status(200).JSON(quick.M{
			"msg": "JSON logging example",
		})
	})

	// Start the server
	q.Listen("0.0.0.0:8080")
}

```

---

## 📏 Maxbody (Request Size Limiter)
Restricts the maximum request body size to prevent clients from sending excessively large payloads.

- Avoids excessive memory usage.
- Can prevent attacks such as DoS (Denial-of-Service).
- Returns a 413 Payload Too Large error when exceeded.

---

## 🔄 MsgUUID
Assigns a UUID (Universally Unique Identifier) to each request.

- Allows easy tracking of requests in logs.
- Useful for distributed systems where tracing requests across services is required.
- Adds a unique identifier to every request automatically.

---

## 🚧 **Coming soon!**
- Etag
- Pprof
- Proxy
- RequestID
- Skip
- Timeout

