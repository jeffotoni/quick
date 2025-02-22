package main

import "github.com/jeffotoni/goquick"

type My struct {
	Name string `json:"name"`
	Year int    `json:"year"`
}

func main() {
	q := quick.New()
	q.Post("/v1/user", func(c *quick.Ctx) error {
		var my My
		err := c.BodyParser(&my)
		if err != nil {
			c.Status(400).SendString(err.Error())
		}

		return c.Status(200).JSON(&my)
		// ou
		//c.Status(200).String(c.BodyString())
	})

	q.Listen("0.0.0.0:8080")
}
