package main

import "github.com/jeffotoni/quick"

func main() {

	q := quick.New()

	q.Get("/", func(c *quick.Ctx) error {
		return quick.NewError(782, "Custom error message")
	})

	q.Listen(":8080")
}
