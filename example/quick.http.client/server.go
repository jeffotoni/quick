package main

import "github.com/jeffotoni/quick"

func main() {

	h := func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		if c.Body() == nil {
			return c.Status(200).Send([]byte(`{"data":"quick is awesome!"}`))
		}
		return c.Status(200).Send(c.Body())
	}

	q := quick.New()

	q.Get("/v1/user/:id", h)
	q.Post("/v1/user", h)
	q.Put("/v1/user/:id", h)
	q.Delete("/v1/user/:id", h)

	q.Listen("0.0.0.0:3000")
}
