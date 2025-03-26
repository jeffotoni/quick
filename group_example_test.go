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

	// As the prefix field is unexported, we test by registering a route and calling it
	apiGroup.Get("/check", func(c *Ctx) error {
		return c.Status(200).String("Prefix OK")
	})

	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/api/check",
	})

	if err := res.AssertStatus(200); err != nil {
		fmt.Println("Status error:", err)
	}

	if err := res.AssertString("Prefix OK"); err != nil {
		fmt.Println("Body error:", err)
	}

	fmt.Println(res.BodyStr())

	// Output: Prefix OK
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
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/api/users",
	})

	if err := res.AssertStatus(200); err != nil {
		fmt.Println("Status error:", err)
	}

	if err := res.AssertString("List of users"); err != nil {
		fmt.Println("Body error:", err)
	}

	fmt.Println(res.BodyStr())

	// Output: List of users
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
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodPost,
		URI:    "/api/users",
	})

	if err := res.AssertStatus(201); err != nil {
		fmt.Println("Status error:", err)
	}

	if err := res.AssertString("User created"); err != nil {
		fmt.Println("Body error:", err)
	}

	// Print the response body
	fmt.Println(res.BodyStr())

	// Output: User created
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
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodPut,
		URI:    "/api/users/42",
	})

	if err := res.AssertStatus(200); err != nil {
		fmt.Println("Status error:", err)
	}

	if err := res.AssertString("User updated"); err != nil {
		fmt.Println("Body error:", err)
	}

	// Print the response body
	fmt.Println(res.BodyStr())

	// Output: User updated
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
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodDelete,
		URI:    "/api/users/42",
	})

	if err := res.AssertStatus(200); err != nil {
		fmt.Println("Status error:", err)
	}

	if err := res.AssertString("User deleted"); err != nil {
		fmt.Println("Body error:", err)
	}

	// Print the response body
	fmt.Println(res.BodyStr())

	// Output: User deleted
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
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodGet,
		URI:    "/api/hello",
	})

	if err := res.AssertStatus(200); err != nil {
		fmt.Println("Status error:", err)
	}

	if err := res.AssertString("Hello from API group Quick"); err != nil {
		fmt.Println("Body error:", err)
	}

	// Print the response body
	fmt.Println(res.BodyStr())

	// Output:
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
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodPatch,
		URI:    "/api/update",
	})

	if err := res.AssertStatus(200); err != nil {
		fmt.Println("Status error:", err)
	}

	if err := res.AssertString("PATCH request received"); err != nil {
		fmt.Println("Body error:", err)
	}

	// Print the response body
	fmt.Println(res.BodyStr())

	// Output: PATCH request received
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
	res, _ := q.Qtest(QuickTestOptions{
		Method: MethodOptions,
		URI:    "/api/resource",
	})

	if err := res.AssertStatus(204); err != nil {
		fmt.Println("Status error:", err)
	}

	// Print the response status
	fmt.Println("Status:", res.StatusCode())

	// Output: Status: 204
}
