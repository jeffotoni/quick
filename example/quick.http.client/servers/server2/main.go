//go:build !exclude_test

package main

import (
	"github.com/jeffotoni/quick"
)

func main() {

	h := func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		if c.Body() == nil {
			return c.Status(200).Send([]byte(`{"data":"quick is awesome!"}`))
		}
		return c.Status(200).Send(c.Body())
	}

	forms := func(c *quick.Ctx) error {

		// Get individual values from the form
		name := c.FormValue("name")
		email := c.FormValue("email")

		// Get all values from the form
		allValues := c.FormValues()

		// Respond with the received data
		return c.Status(200).JSON(map[string]any{
			"message": "Form received",
			"name":    name,
			"email":   email,
			"data":    allValues,
		})
	}

	q := quick.New()

	q.Get("/v1/user/:id", h)
	q.Post("/v1/user", h)
	q.Put("/v1/user/:id", h)
	q.Delete("/v1/user/:id", h)
	q.Post("/postform", forms)

	// show run server
	// export PRINT_SERVER=true
	q.Listen("0.0.0.0:3000")
}
