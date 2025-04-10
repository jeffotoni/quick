
## 📏 Maxbody (Request Size Limiter)
Restricts the maximum request body size to prevent clients from sending excessively large payloads.

- ✅ Avoids excessive memory usage.
- ✅ Can prevent attacks such as DoS (Denial-of-Service).
- ✅ Returns a 413 Payload Too Large error when exceeded.

### 🔹 Simple Example (Using maxbody.New)
This example limits the request body size using maxbody.New(), which applies the restriction globally.

```go
package main

import (
    "log"

    "github.com/jeffotoni/quick"
    "github.com/jeffotoni/quick/middleware/maxbody"
)

func main() {
    q := quick.New()

    // Middleware to enforce a 50KB request body limit
    q.Use(maxbody.New(50000)) // 50KB

    // Define a route that accepts a request body
    q.Post("/v1/user/maxbody/any", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")

        log.Printf("Body received: %s", c.BodyString())
        return c.Status(200).Send(c.Body())
    })

    log.Fatal(q.Listen("0.0.0.0:8080"))
}
```
---
### 🔹 Advanced Example (Using MaxBytesReader)

This example applies MaxBytesReader for additional security by enforcing the body size limit at the request handling level.

```go
package main

import (
    "io"
    "log"
    "net/http"

    "github.com/jeffotoni/quick"
)

const maxBodySize = 1024 // 1KB

func main() {
    q := quick.New()

    // Define a route that applies MaxBytesReader for additional protection
    q.Post("/v1/user/maxbody/max", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")

        // Limit request body size to 1KB
        c.Request.Body = quick.MaxBytesReader(c.Response, c.Request.Body, maxBodySize)

        // Read the request body safely
        body, err := io.ReadAll(c.Request.Body)
        if err != nil {
            log.Printf("Error reading request body: %v", err)
            return c.Status(http.StatusRequestEntityTooLarge).String("Request body too large")
        }
        return c.Status(http.StatusOK).Send(body)
    })

    log.Println("Server running at http://0.0.0.0:8080")
    log.Fatal(q.Listen("0.0.0.0:8080"))
}

```
---

### 📌 Key Differences

| Implementation      | Description                                      |
|--------------------|--------------------------------------------------|
| `maxbody.New()`    | Enforces a global request body size limit.       |
| `MaxBytesReader()` | Adds an extra validation layer inside the request handler. |


---
