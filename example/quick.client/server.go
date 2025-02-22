package main

import "github.com/jeffotoni/goquick"

func main() {

	h := func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		if c.Body() == nil {
			return c.Status(200).Send([]byte(`{"data":"quick is awesome!"}`))
		}
		return c.Status(200).Send(c.Body())
	}

	q := quick.New()

	q.Get("/get", h)
	q.Post("/post", h)
	q.Put("/put", h)
	q.Delete("/delete", h)

	q.Listen(":8000")
}
