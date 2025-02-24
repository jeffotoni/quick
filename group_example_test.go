// Package quick provides a minimalistic and high-performance web framework for Go.
//
// This file contains example implementations demonstrating how to use route groups (Group)
// in the Quick framework. Groups allow better organization of routes by prefixing related
// endpoints under a common path.
//
// These examples showcase different HTTP methods (GET, POST, PUT, DELETE) within a route group.
package quick

import "fmt"

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
