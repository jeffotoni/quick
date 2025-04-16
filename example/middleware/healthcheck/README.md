## 🛠️ Healthcheck Middleware in Quick ![Quick Logo](/quick.png)

**Healthcheck** is a middleware this package provides a simple way to check the health of your application.

---
### ✨ Features

- Simple healthcheck endpoint
- Customizable endpoint

---
### 🧩 Example Usage
```go
package main

import (
	"github.com/jeffotoni/quick"
	"github.com/seuusuario/healthcheck"
	
)

func main() {
	q := quick.New()

	// Use Healthcheck middleware with default healthcheck endpoint
	q.Use(healthcheck.New(
		healthcheck.Options{
			App: q,
		},
	))

	q.Get("/", func(c *quick.Ctx) error {
		return c.Status(200).String("Home page")
	})

	log.Fatalln(q.Listen(":8080"))
}
```
### 📌 cURL
```bash
$ curl -X GET 'http://localhost:8080/healthcheck'
```

### 📥 Example Output

Here's an example of the response returned:

```sh
OK
```

---
### ⚙️ Custom Configuration

You can change the endpoint by providing an Options struct:

```go
q.Use(healthcheck.New(
	healthcheck.Options{
		App: q,
		Endpoint: "/v1/health",
	},
))
```

