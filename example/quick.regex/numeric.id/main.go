package main

import (
	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New()

	// Route that accepts only numeric IDs (using regex [0-9]+)
	q.Get("/users/{id:[0-9]+}", func(c *quick.Ctx) error {
		id := c.Param("id")
		return c.JSON(map[string]string{
			"message": "User found",
			"user_id": id,
		})
	})

	// Start the server on port 8080
	q.Listen(":8080")
}

//$ curl --location 'http://localhost:8080/users/123'
