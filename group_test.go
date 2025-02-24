package quick

import (
	"testing"
)

// TestQuick_Group verifies if a route group is correctly created with the expected prefix.
// The will test TestQuick_Group(t *testing.T)
//
// Run:
//
//	$ go test -v -run ^TestQuick_Group
func TestQuick_Group(t *testing.T) {
	q := New()

	// Create a route group with the prefix "/api"
	apiGroup := q.Group("/api")

	// Expected prefix for the group
	expectedPrefix := "/api"

	// Verify if the group was created with the correct prefix
	if apiGroup.prefix != expectedPrefix {
		t.Errorf("Expected prefix '%s', but got '%s'", expectedPrefix, apiGroup.prefix)
	}

	// Ensure at least one group exists in q.groups
	if len(q.groups) == 0 {
		t.Errorf("Expected at least one group in q.groups, but got %d", len(q.groups))
	}

	// Verify if the first group's prefix matches the expected value
	if q.groups[0].prefix != expectedPrefix {
		t.Errorf("Expected first group's prefix to be '%s', but got '%s'", expectedPrefix, q.groups[0].prefix)
	}
}

// TestGroup_Get verifies if a GET request to a route within a group returns the expected response.
// The will test TestGroup_Get(t *testing.T)

// Run:
//
//	$ go test -v -run ^TestGroup_Get
func TestGroup_Get(t *testing.T) {
	q := New()

	// Create a route group
	apiGroup := q.Group("/api")

	// Define a GET route inside the group
	apiGroup.Get("/users", func(c *Ctx) error {
		return c.Status(200).String("List of users")
	})

	// Simulate a GET request to "/api/users"
	res, err := q.QuickTest("GET", "/api/users", nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	// Validate HTTP status code
	if res.StatusCode() != 200 {
		t.Errorf("Expected status 200, but got %d", res.StatusCode())
	}

	// Validate response body
	expectedBody := "List of users"
	if res.BodyStr() != expectedBody {
		t.Errorf("Expected body '%s', but got '%s'", expectedBody, res.BodyStr())
	}
}

// TestGroup_Post verifies if a POST request creates a resource and returns the expected response.
// The will test TestGroup_Post(t *testing.T)

// Run:
//
//	$ go test -v -run ^TestGroup_Post
func TestGroup_Post(t *testing.T) {
	q := New()

	// Create a route group
	apiGroup := q.Group("/api")

	// Define a POST route inside the group
	apiGroup.Post("/users", func(c *Ctx) error {
		return c.Status(201).String("User created")
	})

	// Simulate a POST request to "/api/users"
	res, err := q.QuickTest("POST", "/api/users", nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	// Validate HTTP status code
	if res.StatusCode() != 201 {
		t.Errorf("Expected status 201, but got %d", res.StatusCode())
	}

	// Validate response body
	expectedBody := "User created"
	if res.BodyStr() != expectedBody {
		t.Errorf("Expected body '%s', but got '%s'", expectedBody, res.BodyStr())
	}
}

// TestGroup_Put verifies if a PUT request updates a resource and returns the expected response.
// The will test TestGroup_Put(t *testing.T)
//
// Run:
//
//	$ go test -v -run ^TestGroup_Put
func TestGroup_Put(t *testing.T) {
	q := New()

	// Create a route group
	apiGroup := q.Group("/api")

	// Define a PUT route inside the group
	apiGroup.Put("/users/:id", func(c *Ctx) error {
		return c.Status(200).String("User updated")
	})

	// Simulate a PUT request to "/api/users/42"
	res, err := q.QuickTest("PUT", "/api/users/42", nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	// Validate HTTP status code
	if res.StatusCode() != 200 {
		t.Errorf("Expected status 200, but got %d", res.StatusCode())
	}

	// Validate response body
	expectedBody := "User updated"
	if res.BodyStr() != expectedBody {
		t.Errorf("Expected body '%s', but got '%s'", expectedBody, res.BodyStr())
	}
}

// TestGroup_Delete verifies if a DELETE request removes a resource and returns the expected response.
// The will test TestGroup_Delete(t *testing.T)
//
// Run:
//
//	$ go test -v -run ^TestGroup_Delete
func TestGroup_Delete(t *testing.T) {
	q := New()

	// Create a route group
	apiGroup := q.Group("/api")

	// Define a DELETE route inside the group
	apiGroup.Delete("/users/:id", func(c *Ctx) error {
		return c.Status(200).String("User deleted")
	})

	// Simulate a DELETE request to "/api/users/42"
	res, err := q.QuickTest("DELETE", "/api/users/42", nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	// Validate HTTP status code
	if res.StatusCode() != 200 {
		t.Errorf("Expected status 200, but got %d", res.StatusCode())
	}

	// Validate response body
	expectedBody := "User deleted"
	if res.BodyStr() != expectedBody {
		t.Errorf("Expected body '%s', but got '%s'", expectedBody, res.BodyStr())
	}
}
