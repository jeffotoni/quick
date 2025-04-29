// Example of basic cache middleware usage in Quick
package main

import (
	"fmt"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/cache"
)

func main() {
	// Create a new Quick app
	q := quick.New()

	// Use the cache middleware with default settings
	q.Use(cache.New())

	// Route 1: Returns the current time
	q.Get("/time", func(c *quick.Ctx) error {
		return c.String("Current time: " + time.Now().Format(time.RFC1123))
	})

	// Route 2: Returns a random number
	q.Get("/random", func(c *quick.Ctx) error {
		return c.String("Random value: " + time.Now().Format("15:04:05.000"))
	})

	// Route 3: Returns JSON data
	q.Get("/profile", func(c *quick.Ctx) error {
		return c.JSON(quick.M{
			"user":  "jeffotoni",
			"since": time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	// Start the server
	fmt.Println("Server running on http://localhost:3000")
	fmt.Println("Try these endpoints:")
	fmt.Println("  - GET /time (cached for 1 minute)")
	fmt.Println("  - GET /random (cached for 1 minute)")
	fmt.Println("  - GET /profile (cached for 1 minute)")
	fmt.Println("Check the X-Cache-Status header in the response to see if it's a HIT or MISS")
	q.Listen(":3000")
}
