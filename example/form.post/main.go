package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New()

	q.Post("/postform", func(c *quick.Ctx) error {

		// Get individual values ​​from the form
		name := c.FormValue("name")
		email := c.FormValue("email")

		// Get all values ​​from the form
		allValues := c.FormValues()

		// Respond with the received data
		return c.Status(200).JSON(map[string]any{
			"message": "Form received",
			"name":    name,
			"email":   email,
			"data":    allValues,
		})
	})

	fmt.Println("Server running on :3000")
	log.Fatal(q.Listen("0.0.0.0:3000"))
}
