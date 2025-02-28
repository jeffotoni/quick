// Package quick provides a minimalistic and high-performance web framework for Go.
//
// This file contains example implementations demonstrating how to use route groups (Group)
// in the Quick framework. Groups allow better organization of routes by prefixing related
// endpoints under a common path.
//
// These examples showcase different HTTP methods (GET, POST, PUT, DELETE) within a route group.
package quick

import (
	"fmt"
	"net/http"
)

// This function is named ExampleQuick_Group()
// it with the Examples type.
func ExampleQuick_Group() {
	q := New()

	// Create a route group with prefix "/api"
	apiGroup := q.Group("/api")

	// Print the prefix of the group
	fmt.Println(apiGroup.prefix)
	// Out put: /api
}

// This function is named ExampleGroup_Get()
// it with the Examples type.
func ExampleGroup_Get() {
	q := New()

	// Create a route group with prefix "/api"
	apiGroup := q.Group("/api")

	// Define a GET route inside the group
	apiGroup.Get("/users", func(c *Ctx) error {
		// Return a success message
		return c.Status(200).String("List of users")
	})

	// Simulate a GET request to "/api/users"
	res, _ := q.QuickTest("GET", "/api/users", nil)

	// Print the response body
	fmt.Println(res.BodyStr())

	// Out put: List of users
}

// This function is named ExampleGroup_Post()
// it with the Examples type.
func ExampleGroup_Post() {
	q := New()

	// Create a route group with prefix "/api"
	apiGroup := q.Group("/api")

	// Define a POST route inside the group
	apiGroup.Post("/users", func(c *Ctx) error {
		// Return a success message
		return c.Status(201).String("User created")
	})

	// Simulate a POST request to "/api/users"
	res, _ := q.QuickTest("POST", "/api/users", nil)

	// Print the response body
	fmt.Println(res.BodyStr())

	// Out put: User created
}

// This function is named ExampleGroup_Put()
// it with the Examples type.
func ExampleGroup_Put() {
	q := New()

	// Create a route group with prefix "/api"
	apiGroup := q.Group("/api")

	// Define a PUT route inside the group
	apiGroup.Put("/users/:id", func(c *Ctx) error {
		// Return a success message
		return c.Status(200).String("User updated")
	})

	// Simulate a PUT request to "/api/users/42"
	res, _ := q.QuickTest("PUT", "/api/users/42", nil)

	// Print the response body
	fmt.Println(res.BodyStr())

	// Out put: User updated
}

// This function is named ExampleGroup_Delete()
// it with the Examples type.
func ExampleGroup_Delete() {
	q := New()

	// Create a route group with prefix "/api"
	apiGroup := q.Group("/api")

	// Define a DELETE route inside the group
	apiGroup.Delete("/users/:id", func(c *Ctx) error {
		// Return a success message
		return c.Status(200).String("User deleted")
	})

	// Simulate a DELETE request to "/api/users/42"
	res, _ := q.QuickTest("DELETE", "/api/users/42", nil)

	// Print the response body
	fmt.Println(res.BodyStr())

	// Out put: User deleted
}

// This function is named ExampleGroup_Delete()
// it with the Examples type.
func ExampleGroup_Use() {
	// Create a new Quick instance
	q := New()

	// Create a new group with a common prefix
	api := q.Group("/api")

	// Define a simple middleware that logs requests
	logMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Middleware activated for:", r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}

	// Apply middleware to the group
	api.Use(logMiddleware)

	// Define a GET route inside the group
	api.Get("/hello", func(c *Ctx) error {
		return c.Status(200).String("Hello from API group Quick")
	})

	// Simulate a request to test middleware activation
	res, _ := q.QuickTest("GET", "/api/hello", nil)

	// Print the response body
	fmt.Println(res.BodyStr())

	// Out put:
	// Middleware activated for: /api/hello
	// Hello from API group Quick
}

// This function is named ExampleGroup_Patch()
// it with the Examples type.
func ExampleGroup_Patch() {
	// Create a new Quick instance
	q := New()

	// Create a new group with a common prefix
	api := q.Group("/api")

	// Register a PATCH route dynamically
	api.Patch("/update", func(c *Ctx) error {
		return c.Status(200).String("PATCH request received")
	})

	// Simulate a PATCH request
	res, _ := q.QuickTest("PATCH", "/api/update", nil)

	// Print the response body
	fmt.Println(res.BodyStr())

	// Out put:
	// PATCH request received
}

// This function is named ExampleGroup_Options()
// it with the Examples type.
func ExampleGroup_Options() {
	// Create a new Quick instance
	q := New()

	// Create a new group with a common prefix
	api := q.Group("/api")

	// Register an OPTIONS route dynamically
	api.Options("/resource", func(c *Ctx) error {
		c.Set("Allow", "GET, POST, OPTIONS")
		return c.Status(204).Send(nil) // No Content response
	})

	// Simulate an OPTIONS request
	res, _ := q.QuickTest("OPTIONS", "/api/resource", nil)

	// Print the response status
	fmt.Println("Status:", res.StatusCode())

	// Out put:
	// Status: 204
}
