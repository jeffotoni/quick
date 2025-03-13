package main

import (
	"fmt"
	"net/http"

	"github.com/jeffotoni/quick"
)

// This example demonstrates advanced usage of Quick with custom configurations and middlewares
func main() {
	// Initialize Quick with a custom configuration
	config := quick.Config{
		RouteCapacity: 500, // Custom route capacity
	}
	q := quick.New(config)

	// Middleware example: Logging requests
	q.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("Received request: %s %s\n", r.Method, r.URL.Path)
			h.ServeHTTP(w, r)
		})
	})

	// Define multiple routes
	q.Get("/", func(c *quick.Ctx) error {
		return c.Status(200).String("Welcome to Quick!")
	})

	q.Post("/data", func(c *quick.Ctx) error {
		return c.Status(201).String("Data received!")
	})

	q.Put("/update", func(c *quick.Ctx) error {
		return c.Status(200).String("Data updated!")
	})

	q.Delete("/delete", func(c *quick.Ctx) error {
		return c.Status(204).String("") // No content response
	})

	// Start the server
	fmt.Println("Server running at http://localhost:8080")
	q.Listen(":8080")
}

// $ curl -X GET http://localhost:8080/

// $ curl -X POST http://localhost:8080/data

// $ curl -X PUT http://localhost:8080/update

// $ curl -X DELETE http://localhost:8080/delete
