## 🎗️ Recover Middleware in Quick ![Quick Logo](/quick.png)

**Recover** is a middleware this package provides a simple way to handle panics in your application and prints stack trace.

---
### 🧩 Example Usage
```go
package main

import (
	"errors"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/recover"
)

func main() {
	q := quick.New()

	// Apply the Recover middleware
	q.Use(recover.New(recover.Config{
		App: q,
	}))

	// Define a test route
	q.Get("/v1/recover", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// halt the server
		panic(errors.New("Panicking!"))
	})

	// Start the server
	q.Listen("0.0.0.0:8080")
}
```

---
### 📌 cURL
```bash
$ curl -i -XGET http://localhost:8080/v1/recover
```

### 📥 Example Output

Here's an example of the response returned:

```sh
Internal Server Error
```
