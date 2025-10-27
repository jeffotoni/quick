package main


import (
	"testing"

	"github.com/jeffotoni/quick"
)

func setupApp() *quick.Quick {
	q := quick.New(quick.Config{NoBanner: true})

	// PUT route to update a user by ID
	q.Put("/users/:id", func(c *quick.Ctx) error {
		userID := c.Param("id")
		c.Set("Content-Type", "text/plain")
		c.Set("X-User-ID", userID)
		return c.Status(200).SendString("User " + userID + " updated successfully!")
	})

	// PUT route to update a specific type by ID
	q.Put("/tipos/:id", func(c *quick.Ctx) error {
		tiposID := c.Param("id")
		c.Set("Content-Type", "application/json")
		c.Set("X-Type-ID", tiposID)
		return c.Status(200).SendString("User " + tiposID + " type updated successfully!")
	})

	return q
}

func TestPutUsers(t *testing.T) {
	app := setupApp()

	resp, err := app.Qtest(quick.QuickTestOptions{
		Method: "PUT",
		URI:    "/users/123",
	})

	if err != nil {
		t.Fatal(err)
	}

	// Testar status
	if err := resp.AssertStatus(200); err != nil {
		t.Error(err)
	}

	// Testar header
	if err := resp.AssertHeader("Content-Type", "text/plain"); err != nil {
		t.Error(err)
	}

	if err := resp.AssertString("User 123 updated successfully!"); err != nil {
		t.Error(err)
	}
}

func TestPutTypes(t *testing.T) {
	app := setupApp()

	resp, err := app.Qtest(quick.QuickTestOptions{
		Method: "PUT",
		URI:    "/tipos/456",
	})

	if err != nil {
		t.Fatal(err)
	}

	// Test status
	if err := resp.AssertStatus(200); err != nil {
		t.Error(err)
	}

	// Test header
	if err := resp.AssertHeader("Content-Type", "application/json"); err != nil {
		t.Error(err)
	}

	// Test response
	if err := resp.AssertString("User 456 type updated successfully!"); err != nil {
		t.Error(err)
	}
}

func TestBodyMethods(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "PUT",
		URI:    "/users/789",
	})

	body := resp.Body()
	if len(body) == 0 {
		t.Error("Body vazio")
	}

	bodyStr := resp.BodyStr()
	if bodyStr != "User 789 updated successfully!" {
		t.Error("BodyStr incorreto")
	}
}

func TestStatusCode(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "PUT",
		URI:    "/tipos/999",
	})

	if resp.StatusCode() != 200 {
		t.Error("StatusCode incorreto")
	}
}

func TestResponse(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "PUT",
		URI:    "/users/111",
	})

	if resp.Response() == nil {
		t.Error("Response nil")
	}
}

func TestAssertNoHeader(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "PUT",
		URI:    "/users/222",
	})

	if err := resp.AssertNoHeader("X-Invalid-Header"); err != nil {
		t.Error(err)
	}
}

func TestAssertBodyContains(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "PUT",
		URI:    "/users/333",
	})

	if err := resp.AssertBodyContains("updated"); err != nil {
		t.Error(err)
	}
}

func TestAssertHeaderContains(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "PUT",
		URI:    "/users/444",
	})

	if err := resp.AssertHeaderContains("Content-Type", "text"); err != nil {
		t.Error(err)
	}
}

func TestAssertHeaderHasPrefix(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "PUT",
		URI:    "/tipos/555",
	})

	if err := resp.AssertHeaderHasPrefix("Content-Type", "application"); err != nil {
		t.Error(err)
	}
}

func TestAssertHeaderHasValueInSet(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "PUT",
		URI:    "/tipos/666",
	})

	allowed := []string{"application/json", "text/plain"}
	if err := resp.AssertHeaderHasValueInSet("Content-Type", allowed); err != nil {
		t.Error(err)
	}
}

func TestAssertCustomHeaders(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "PUT",
		URI:    "/users/777",
	})

	if err := resp.AssertHeader("X-User-ID", "777"); err != nil {
		t.Error(err)
	}
}

