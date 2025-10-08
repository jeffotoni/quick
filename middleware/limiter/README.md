
## 🚦 Rate Limiter ![Quick Logo](/quick.png) 

The **Rate Limiter** is a middleware for the Quick framework that controls the number of requests allowed in a given time period. It helps prevent API abuse and improves system stability by preventing server overload.

### 🚀 Features

| Feature                         | Description |
|----------------------------------|-------------|
| 🎯 **Request Rate Limiting**     | Configurable maximum number of requests per client within a time window. |
| ⏳ **Automatic Expiration**      | Resets the request counter automatically after the configured time. |
| 🔑 **Custom Client Identification** | Uses a `KeyGenerator` function to define a unique client key (e.g., IP-based). |
| ⚠️ **Custom Response on Limit Exceeded** | Allows defining a custom response when the request limit is reached. |
| ⚡ **Efficient Performance**     | Implements sharding and optimizations to reduce concurrency issues. |

### 🌍 Global Rate Limiter 

The example below shows how to apply the Rate Limiter as global middleware.
```go
package main

import (
	"time"
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/limiter"
)

func main() {
	q := quick.New()

	// Apply the rate limiter middleware
	q.Use(limiter.New(limiter.Config{
		// Maximum 10 requests allowed per IP
		Max: 10,
		// The limit resets after 5 seconds
		Expiration: 5 * time.Second,
		KeyGenerator: func(c *quick.Ctx) string {
			// Uses the client's IP address as the key
			return c.RemoteIP()
		},
		LimitReached: func(c *quick.Ctx) error {
			c.Set("Content-Type", "application/json")
			// The client should wait 10 seconds before retrying
			c.Set("Retry-After", "10")
			return c.Status(quick.StatusTooManyRequests).JSON(map[string]string{
				"error":   "Too many requests",
				"message": "You have exceeded the request limit. 
				Please wait 1 second and try again.",
				"retry_after": "10s",
			})
		},
	}))

	// Define a simple GET route
	q.Get("/", func(c *quick.Ctx) error {
		return c.Status(200).JSON(map[string]string{"msg": "Quick in action ❤️!"})
	})

	// Start the server on port 8080
	q.Listen(":8080")
}


```
### 📌 cURL

```bash
$ curl -i -X GET http://localhost:8080/
```
If the same IP makes more than 10 requests in 5 seconds, the middleware returns:

```bash
{
    "error": "Too many requests",
    "message": "You have exceeded the request limit. 
	 Please wait 1 second and try again.",
    "retry_after": "10s"
}