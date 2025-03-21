## 🛡️ Helmet Middleware in Quick ![Quick Logo](/quick.png)

**Helmet** is a middleware this package provides sensible security defaults while allowing full customization.

---

### ✨ Features

- Sets common security-related HTTP headers
- Provides secure defaults
- Easily customizable via `Options` struct
- Supports skipping middleware per request

---

### 🛡️ Default Headers

By default, the middleware sets the following headers:

- X-XSS-Protection
- X-Content-Type-Options
- X-Frame-Options
- Content-Security-Policy
- Referrer-Policy
- Permissions-Policy
- Cross-Origin-Embedder-Policy
- Cross-Origin-Opener-Policy
- Cross-Origin-Resource-Policy
- Origin-Agent-Cluster
- X-DNS-Prefetch-Control
- X-Download-Options
- X-Permitted-Cross-Domain-Policies
- Strict-Transport-Security (only for HTTPS requests)
- Cache-Control

---

### 🧩 Example Usage
```go
package main

import (
	"github.com/jeffotoni/quick"
	"github.com/seuusuario/helmet"
)

func main() {
	q := quick.New()

	// Use Helmet middleware with default security headers
	q.Use(helmet.Helmet())

	// Simple route to test headers
	q.Get("/v1/user", func(c *quick.Ctx) error {

		// list all headers
		headers := make(map[string]string)
		for k, v := range c.Response.Header() {
			if len(v) > 0 {
				headers[k] = v[0]
			}
		}
		return c.Status(200).JSONIN(headers)
	})

	q.Listen("0.0.0.0:8080")
}
```
### 📌 cURL
```bash
$ curl -X GET 'http://localhost:8080/v1/user'
```

### 📥 Example Output

Here's an example of the response headers returned:

```go
{
  "Cache-Control": "no-cache, no-store, must-revalidate",
  "Content-Security-Policy": "default-src 'self'",
  "Cross-Origin-Embedder-Policy": "require-corp",
  "Cross-Origin-Opener-Policy": "same-origin",
  "Cross-Origin-Resource-Policy": "same-origin",
  "Origin-Agent-Cluster": "?1",
  "Referrer-Policy": "no-referrer",
  "X-Content-Type-Options": "nosniff",
  "X-DNS-Prefetch-Control": "off",
  "X-Download-Options": "noopen",
  "X-Frame-Options": "SAMEORIGIN",
  "X-Permitted-Cross-Domain-Policies": "none",
  "X-XSS-Protection": "0"
}
```
---
### ⚙️ Custom Configuration

You can override any of the default headers by providing an Options struct:

```go
app.Use(helmet.Helmet(helmet.Options{
	XSSProtection:         "1; mode=block",
	ContentSecurityPolicy: "default-src 'self'; script-src 'none'",
	XFrameOptions:         "DENY",
	HSTSMaxAge:            63072000,
	HSTSPreloadEnabled:    true,
}))
```

You can also use the `Next` function to conditionally skip the middleware:

```go
app.Use(helmet.Helmet(helmet.Options{
	Next: func(c *quick.Ctx) bool {
		// Skip for health checks
		return c.Path() == "/health"
	},
}))
```
