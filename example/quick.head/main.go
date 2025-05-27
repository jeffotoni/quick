package main

import (
	"log"

	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New()

	// Register a HEAD route
	q.Head("/ping", func(c *quick.Ctx) error {
		c.Set("Content-Type", "text/plain")
		return c.String("pong!") // Will not be included in response body for HEAD
	})

	// Optional: Register a GET to show difference
	q.Get("/ping", func(c *quick.Ctx) error {
		c.Set("Content-Type", "text/plain")
		return c.String("pong!")
	})

	log.Println("Server running at http://localhost:8080")
	if err := q.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}

// curl -I -XHEAD http://localhost:8080/ping
// curl -i http://localhost:8080/ping
