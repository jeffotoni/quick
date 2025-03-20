// Package quick provides a fast and flexible web framework with built-in
// HTTP testing utilities.
package quick

import (
	"fmt"
	"net/http"
	"testing"
)

// TestQTest_Options_GET verifies if the response body contains the expected substring.
//
// This test performs a GET request to "/v1/user" and checks if:
//   - The response status is 200 (OK).
//   - The response body contains "Success".
func TestQTest_Options_GET(t *testing.T) {
	q := New()

	q.Get("/v1/user", func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(StatusOK).String("Success")
	})

	opts := QuickTestOptions{
		Method:     "GET",
		URI:        "/v1/user",
		Headers:    map[string]string{"Accept": "application/json"},
		LogDetails: true,
	}

	result, err := q.Qtest(opts)
	if err != nil {
		t.Fatalf("Error in test: %v", err)
	}

	if err := result.AssertStatus(StatusOK); err != nil {
		t.Errorf("Status assertion failed: %v", err)
	}

	fmt.Println("DEBUG Body (QTest):", result.BodyStr())

	if err := result.AssertBodyContains("Success"); err != nil {
		t.Errorf("Body assertion failed: %v", err)
	}
}

// TestQTest_Options_POST verifies if the response body contains the expected JSON message.
//
// This test performs a POST request to "/v1/user/api" and checks if:
//   - The response status is 200 (OK).
//   - The "Content-Type" header is "application/json".
//   - The response body contains `"message":"Success"`.
func TestQTest_Options_POST(t *testing.T) {
	q := New()

	q.Post("/v1/user/api", func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		c.Response.Header().Set("Content-Type", "application/json")
		return c.Status(StatusOK).String(`{"message":"Success"}`)
	})

	opts := QuickTestOptions{
		Method: "POST",
		URI:    "/v1/user/api",
		QueryParams: map[string]string{
			"param1": "value1",
			"param2": "value2",
		},
		Body: []byte(`{"key":"value"}`),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Cookies: []*http.Cookie{
			{Name: "session", Value: "abc123"},
		},
		LogDetails: true,
	}

	result, err := q.Qtest(opts)
	if err != nil {
		t.Fatalf("Error in Qtest: %v", err)
	}

	if err := result.AssertStatus(StatusOK); err != nil {
		t.Errorf("Status assertion failed: %v", err)
	}

	if err := result.AssertHeader("Content-Type", "application/json"); err != nil {
		t.Errorf("Header assertion failed: %v", err)
	}

	if err := result.AssertBodyContains("Success"); err != nil {
		t.Errorf("Body assertion failed: %v", err)
	}
}

// TestQTest_Options_PUT verifies if the response body contains the expected JSON message.
//
// This test performs a PUT request to "/v1/user/update" and checks if:
//   - The response status is 200 (OK).
//   - The response body contains `"message":"User updated successfully"`.
func TestQTest_Options_PUT(t *testing.T) {
	q := New()

	q.Put("/v1/user/update", func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(StatusOK).String(`{"message":"User updated successfully"}`)
	})

	opts := QuickTestOptions{
		Method: "PUT",
		URI:    "/v1/user/update",
		Body:   []byte(`{"name":"Jeff Quick","age":30}`),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		LogDetails: true,
	}

	result, err := q.Qtest(opts)
	if err != nil {
		t.Fatalf("Error in Qtest: %v", err)
	}

	if err := result.AssertStatus(StatusOK); err != nil {
		t.Errorf("Status assertion failed: %v", err)
	}

	if err := result.AssertBodyContains(`"message":"User updated successfully"`); err != nil {
		t.Errorf("Body assertion failed: %v", err)
	}
}

// TestQTest_Options_DELETE verifies if the response body contains the expected JSON message.
//
// This test performs a DELETE request to "/v1/user/delete" and checks if:
//   - The response status is 200 (OK).
//   - The response body contains `"message":"User deleted successfully"`.
func TestQTest_Options_DELETE(t *testing.T) {
	q := New()

	q.Delete("/v1/user/delete", func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(StatusOK).String(`{"message":"User deleted successfully"}`)
	})

	opts := QuickTestOptions{
		Method:     "DELETE",
		URI:        "/v1/user/delete",
		LogDetails: true,
	}

	result, err := q.Qtest(opts)
	if err != nil {
		t.Fatalf("Error in Qtest: %v", err)
	}

	if err := result.AssertStatus(StatusOK); err != nil {
		t.Errorf("Status assertion failed: %v", err)
	}

	if err := result.AssertBodyContains(`"message":"User deleted successfully"`); err != nil {
		t.Errorf("Body assertion failed: %v", err)
	}
}

// TestQTest_Options_PATCH verifies if the response body contains the expected JSON message.
//
// This test performs a PATCH request to "/v1/user/patch" and checks if:
//   - The response status is 200 (OK).
//   - The response body contains `"message":"User patched successfully"`.
func TestQTest_Options_PATCH(t *testing.T) {
	q := New()

	q.Patch("/v1/user/patch", func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(StatusOK).String(`{"message":"User patched successfully"}`)
	})

	opts := QuickTestOptions{
		Method: "PATCH",
		URI:    "/v1/user/patch",
		Body:   []byte(`{"nickname":"Johnny"}`),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		LogDetails: true,
	}

	result, err := q.Qtest(opts)
	if err != nil {
		t.Fatalf("Error in Qtest: %v", err)
	}

	if err := result.AssertStatus(StatusOK); err != nil {
		t.Errorf("Status assertion failed: %v", err)
	}

	if err := result.AssertBodyContains(`"message":"User patched successfully"`); err != nil {
		t.Errorf("Body assertion failed: %v", err)
	}
}

// TestQTest_Options_OPTIONS verifies if the response status matches 204 (No Content).
//
// This test performs an OPTIONS request to "/v1/user/options" and checks if:
//   - The response status is 204 (No Content).
func TestQTest_Options_OPTIONS(t *testing.T) {
	q := New()

	q.Options("/v1/user/options", func(c *Ctx) error {
		return c.Status(StatusNoContent).String("")
	})

	opts := QuickTestOptions{
		Method:     "OPTIONS",
		URI:        "/v1/user/options",
		LogDetails: true,
	}

	result, err := q.Qtest(opts)
	if err != nil {
		t.Fatalf("Error in Qtest: %v", err)
	}

	if err := result.AssertStatus(StatusNoContent); err != nil {
		t.Errorf("Status assertion failed: %v", err)
	}
}
