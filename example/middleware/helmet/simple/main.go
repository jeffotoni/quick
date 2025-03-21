package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/helmet"
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

	q.Listen(":8080")
}

// $ curl -X GET 'http://localhost:8080/v1/user'
