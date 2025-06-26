package main

import (
	"fmt"

	"github.com/jeffotoni/quick"
)

type P struct {
	Code string `json:"code"`
}

func main() {
	q := quick.New()

	q.Get("/v1/user/:code", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		code := c.Param("code")
		fmt.Println("Code: ", code)
		return c.Status(quick.StatusOK).JSON(P{
			Code: code,
		})
	})

	q.Post("/v1/user/:code", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		code := c.Param("code")
		fmt.Println("Code: ", code)
		return c.Status(quick.StatusOK).JSON(P{
			Code: code,
		})
	})

	/// Start the server on port 8080
	q.Listen("0.0.0.0:8080")
}
