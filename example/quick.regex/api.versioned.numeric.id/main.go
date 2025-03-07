package main

import (
	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New()

	// Route that accepts an API version (v1, v2, etc.) and a numeric user ID
	// - The {version:v[0-9]+} parameter ensures that the version starts with "v" followed by digits.
	// - The {id:[0-9]+} parameter ensures that the user ID consists of only numbers.
	q.Get("/api/{version:v[0-9]+}/users/{id:[0-9]+}", func(c *quick.Ctx) error {
		version := c.Param("version")
		id := c.Param("id")
		return c.JSON(map[string]string{
			"message": "API Versioned User",
			"version": version,
			"user_id": id,
		})
	})

	// Start the server on port 8080
	q.Listen(":8080")
}

// $ curl --location 'http://localhost:8080/api/v1/users/456'
