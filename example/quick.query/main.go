package main

import (
	"fmt"

	"github.com/jeffotoni/quick"
)

// type P struct {
//     Code string `json:"code"`
// }

type Q struct {
	Query map[string]string `json:"query"`
}

func main() {
	q := quick.New()

	q.Get("/v1/user", func(c *quick.Ctx) error {
		q := c.Query
		fmt.Println("query: ", q)
		return c.JSON(Q{
			Query: q,
		})
	})

	q.Post("/v1/user", func(c *quick.Ctx) error {
		q := c.Query
		fmt.Println("query: ", q)
		return c.JSON(Q{
			Query: q,
		})
	})

	/// Start the server on port 8080
	q.Listen("0.0.0.0:8080")
}
