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

| Feature                                             | Status | Notes / Observations                                 |
|-----------------------------------------------------|:------:|------------------------------------------------------|
| `X-XSS-Protection` header                           |   ✅   | Legacy protection, still included                    |
| `X-Content-Type-Options: nosniff` header            |   ✅   | Prevents MIME sniffing attacks                       |
| `X-Frame-Options` header                            |   ✅   | Helps prevent clickjacking                           |
| `Content-Security-Policy` header                    |   ✅   | Defaults to `default-src 'self'`                     |
| `CSPReportOnly` support                             |   ✅   | Optional report-only mode for CSP                    |
| `Referrer-Policy` header                            |   ✅   | Defaults to `no-referrer`                            |
| `Permissions-Policy` header                         |   ✅   | Controls browser features like camera, mic, etc.     |
| `Strict-Transport-Security (HSTS)` support          |   ✅   | Adds HSTS for HTTPS requests                         |
| HSTS options: `maxAge`, `includeSubDomains`, `preload` | ✅   | Fully customizable                                  |
| `Cache-Control` header                              |   ✅   | Defaults to no-cache, improves response integrity    |
| `Cross-Origin-Embedder-Policy` header               |   ✅   | Required for certain advanced browser APIs           |
| `Cross-Origin-Opener-Policy` header                 |   ✅   | Isolates browsing contexts                           |
| `Cross-Origin-Resource-Policy` header               |   ✅   | Restricts resource access                            |
| `Origin-Agent-Cluster` header                       |   ✅   | Enables memory isolation in browsers                 |
| `X-DNS-Prefetch-Control` header                     |   ✅   | Controls browser DNS prefetching                     |
| `X-Download-Options` header                         |   ✅   | Prevents automatic downloads (IE-specific)           |
| `X-Permitted-Cross-Domain-Policies` header          |   ✅   | Blocks Flash and Silverlight legacy access           |
| `Next func(c)` to skip middleware dynamically       |   ✅   | Allows conditional header injection per route        |
| Secure defaults applied when no options are provided|   ✅   | Based on OWASP and best practices                    |
| Option naming compatible with Fiber                 |   ✅   | Enables easier migration from Fiber to Quick         |
| Built-in TLS simulation support in `Qtest`          |   ✅   | Enables full testing of HTTPS-only behavior          |
| Full HTTP method coverage in `Qtest`                |   ✅   | GET, POST, PUT, PATCH, DELETE, OPTIONS supported     |
| Extended Qtest assertions (headers, body, etc.)     |   ✅   | Includes `AssertString`, `AssertNoHeader`, and more  |

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
---

### 📌 cURL
```bash
$ curl -X GET 'http://localhost:8080/v1/user'
```

### 📌 Response

```bash
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