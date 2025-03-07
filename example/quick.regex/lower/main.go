package main

import (
	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New()

	// Route that accepts only lowercase slugs (words with lowercase letters)
	// - The {slug:[a-z]+} parameter ensures that the slug consists only of lowercase letters (a-z).
	// - If uppercase letters or numbers are included, the request will not match.
	q.Get("/profile/{slug:[a-z]+}", func(c *quick.Ctx) error {
		slug := c.Param("slug")
		return c.JSON(map[string]string{
			"message": "Profile found",
			"profile": slug,
		})
	})

	// Start the server on port 8080
	q.Listen(":8080")
}

// $curl --location 'http://localhost:8080/profile/golang'
