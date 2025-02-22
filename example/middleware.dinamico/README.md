# Middleware Dinâmico

```go
package main

import (
	"net/http"

	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New()

	q.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Block") == "" {
				w.WriteHeader(400)
				w.Write([]byte("Sua chamada não irá continuar, preciso definir seu Header com Block: false para passar"))
				return
			}

			if r.Header.Get("Block") == "true" {
				w.WriteHeader(200)
				w.Write([]byte("Sua messgem está bloqueada, coloque false em seu parametro Block"))
				return
			}
			h.ServeHTTP(w, r)
		})
	})

	q.Get("/greet/:name", func(c *quick.Ctx) error {
		name := c.Param("name")
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString("Olá " + name + "!")
	})

	q.Listen("0.0.0.0:8080")
}
```
```go
curl --location --request GET 'http://localhost:8080/greet/:name' \
--header 'Content-Type: application/json/' \
--data 'Sua chamada não irá continuar, preciso definir seu Header com Block: false para passar'
```