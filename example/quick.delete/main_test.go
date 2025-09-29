package main

import (
	"net/http"
	"testing"

	"github.com/jeffotoni/quick"
)

func setupApp() *quick.Quick {
	q := quick.New(quick.Config{NoBanner: true})

	// Simulating a "database" with pre-registered users
	users := map[string]User{
		"1": {Name: "Maria", Year: 2000}, // Fixed user with ID 1
		"2": {Name: "João", Year: 1995},  // Additional user for testing
	}

	// DELETE route to remove a user by ID
	q.Delete("/v1/user/:id", func(c *quick.Ctx) error {
		userID := c.Params["id"] // Retrieve user ID from URL parameter

		// Check if the user exists in the "database"
		if _, exists := users[userID]; !exists {
			return c.Status(http.StatusNotFound).JSON(map[string]string{"error": "User not found"})
		}

		// Delete the user from the "database"
		delete(users, userID)

		// Return a success response
		return c.Status(http.StatusOK).JSON(map[string]string{"msg": "User deleted successfully!"})
	})

	return q
}

func TestDeleteUserSuccess(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "DELETE",
		URI:    "/v1/user/1",
	})

	resp.AssertStatus(200)
	resp.AssertBodyContains("User deleted successfully!")
}

func TestDeleteUserNotFound(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "DELETE",
		URI:    "/v1/user/999",
	})

	resp.AssertStatus(404)
	resp.AssertBodyContains("User not found")
}

func TestBodyMethods(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "DELETE",
		URI:    "/v1/user/2",
	})

	body := resp.Body()
	if len(body) == 0 {
		t.Error("Body vazio")
	}

	bodyStr := resp.BodyStr()
	if bodyStr == "" {
		t.Error("BodyStr vazio")
	}
}

func TestStatusCode(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "DELETE",
		URI:    "/v1/user/1",
	})

	if resp.StatusCode() != 200 {
		t.Error("StatusCode incorreto")
	}
}

func TestResponse(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "DELETE",
		URI:    "/v1/user/1",
	})

	if resp.Response() == nil {
		t.Error("Response nil")
	}
}

func TestAssertNoHeader(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "DELETE",
		URI:    "/v1/user/1",
	})

	resp.AssertNoHeader("X-Invalid-Header")
}

func TestAssertBodyContains(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "DELETE",
		URI:    "/v1/user/1",
	})

	resp.AssertBodyContains("deleted")
}

func TestAssertHeaderContains(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "DELETE",
		URI:    "/v1/user/1",
	})

	resp.AssertHeaderContains("Content-Type", "application/json")
}

func TestAssertHeaderHasPrefix(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "DELETE",
		URI:    "/v1/user/1",
	})

	resp.AssertHeaderHasPrefix("Content-Type", "application/json")
}

func TestAssertHeaderHasValueInSet(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "DELETE",
		URI:    "/v1/user/1",
	})

	allowed := []string{"application/json", "text/plain"}
	resp.AssertHeaderHasValueInSet("Content-Type", allowed)
}

func TestDeleteMultipleScenarios(t *testing.T) {
	app := setupApp()
	
	// Test successful deletion
	resp1, _ := app.Qtest(quick.QuickTestOptions{
		Method: "DELETE",
		URI:    "/v1/user/1",
	})
	resp1.AssertStatus(200)
	resp1.AssertBodyContains("successfully")

	// Test not found
	resp2, _ := app.Qtest(quick.QuickTestOptions{
		Method: "DELETE",
		URI:    "/v1/user/999",
	})
	resp2.AssertStatus(404)
	resp2.AssertBodyContains("not found")
}