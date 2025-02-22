package goquick

import (
	"fmt"
	"testing"
)

// This function is named ExampleQuick_Group()
// it with the Examples type.
func ExampleQuick_Group() {
	q := New()

	apiGroup := q.Group("/api")

	fmt.Println(apiGroup.prefix)

	// Out put: /api
}

// This function is named ExampleGroup_Get()
// it with the Examples type.
func ExampleGroup_Get() {
	q := New()

	apiGroup := q.Group("/api")

	apiGroup.Get("/users", func(c *Ctx) error {
		return c.Status(200).String("List of users")
	})

	res, _ := q.QuickTest("GET", "/api/users", nil)

	fmt.Println(res.BodyStr())

	// Out put: List of users
}

// This function is named ExampleGroup_Post()
// it with the Examples type.
func ExampleGroup_Post() {
	q := New()

	apiGroup := q.Group("/api")

	apiGroup.Post("/users", func(c *Ctx) error {
		return c.Status(201).String("User created")
	})

	res, _ := q.QuickTest("POST", "/api/users", nil)

	fmt.Println(res.BodyStr())

	// Out put: User created
}

// This function is named ExampleGroup_Put()
// it with the Examples type.
func ExampleGroup_Put() {
	q := New()

	apiGroup := q.Group("/api")

	apiGroup.Put("/users/:id", func(c *Ctx) error {
		return c.Status(200).String("User updated")
	})

	res, _ := q.QuickTest("PUT", "/api/users/42", nil)

	fmt.Println(res.BodyStr())

	// Out put: User updated
}

// This function is named ExampleGroup_Delete()
// it with the Examples type.
func ExampleGroup_Delete() {
	q := New()

	apiGroup := q.Group("/api")

	apiGroup.Delete("/users/:id", func(c *Ctx) error {
		return c.Status(200).String("User deleted")
	})

	res, _ := q.QuickTest("DELETE", "/api/users/42", nil)

	fmt.Println(res.BodyStr())

	// Out put: User deleted
}

// go test -v -run ^TestQuick_Group
func TestQuick_Group(t *testing.T) {
	q := New()

	apiGroup := q.Group("/api")

	expectedPrefix := "/api"
	if apiGroup.prefix != expectedPrefix {
		t.Errorf("Expected prefix '%s', but got '%s'", expectedPrefix, apiGroup.prefix)
	}

	if len(q.groups) == 0 {
		t.Errorf("Expected at least one group in q.groups, but got %d", len(q.groups))
	}

	if q.groups[0].prefix != expectedPrefix {
		t.Errorf("Expected first group's prefix to be '%s', but got '%s'", expectedPrefix, q.groups[0].prefix)
	}
}

// go test -v -run ^TestGroup_Get
func TestGroup_Get(t *testing.T) {
	q := New()

	apiGroup := q.Group("/api")

	apiGroup.Get("/users", func(c *Ctx) error {
		return c.Status(200).String("List of users")
	})

	res, err := q.QuickTest("GET", "/api/users", nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	if res.StatusCode() != 200 {
		t.Errorf("Expected status 200, but got %d", res.StatusCode())
	}

	expectedBody := "List of users"
	if res.BodyStr() != expectedBody {
		t.Errorf("Expected body '%s', but got '%s'", expectedBody, res.BodyStr())
	}
}

// go test -v -run ^TestGroup_Post
func TestGroup_Post(t *testing.T) {
	q := New()

	apiGroup := q.Group("/api")

	apiGroup.Post("/users", func(c *Ctx) error {
		return c.Status(201).String("User created")
	})

	res, err := q.QuickTest("POST", "/api/users", nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	if res.StatusCode() != 201 {
		t.Errorf("Expected status 201, but got %d", res.StatusCode())
	}

	expectedBody := "User created"
	if res.BodyStr() != expectedBody {
		t.Errorf("Expected body '%s', but got '%s'", expectedBody, res.BodyStr())
	}
}

// go test -v -run ^TestGroup_Put
func TestGroup_Put(t *testing.T) {
	q := New()

	apiGroup := q.Group("/api")

	apiGroup.Put("/users/:id", func(c *Ctx) error {
		return c.Status(200).String("User updated")
	})

	res, err := q.QuickTest("PUT", "/api/users/42", nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	if res.StatusCode() != 200 {
		t.Errorf("Expected status 200, but got %d", res.StatusCode())
	}

	expectedBody := "User updated"
	if res.BodyStr() != expectedBody {
		t.Errorf("Expected body '%s', but got '%s'", expectedBody, res.BodyStr())
	}
}

// go test -v -run ^TestGroup_Delete
func TestGroup_Delete(t *testing.T) {
	q := New()

	apiGroup := q.Group("/api")

	apiGroup.Delete("/users/:id", func(c *Ctx) error {
		return c.Status(200).String("User deleted")
	})

	res, err := q.QuickTest("DELETE", "/api/users/42", nil)
	if err != nil {
		t.Fatalf("QuickTest failed: %v", err)
	}

	if res.StatusCode() != 200 {
		t.Errorf("Expected status 200, but got %d", res.StatusCode())
	}

	expectedBody := "User deleted"
	if res.BodyStr() != expectedBody {
		t.Errorf("Expected body '%s', but got '%s'", expectedBody, res.BodyStr())
	}
}
