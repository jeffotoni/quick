package quick

import (
	"fmt"
	"net/http"
	"testing"
)

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

	// Debug the response body for troubleshooting
	fmt.Println("DEBUG Body (QTest):", result.BodyStr())

	if err := result.AssertBodyContains("Success"); err != nil {
		t.Errorf("Body assertion failed: %v", err)
	}
}
func TestQTest_Options_POST(t *testing.T) {

	q := New()

	q.Post("/v1/user/api", func(c *Ctx) error {
		c.Set("Content-Type", "application/json")
		c.Response.Header().Set("Content-Type", "application/json")
		return c.Status(StatusOK).String(`{"message":"Success"}`)
	})

	opts := QuickTestOptions{
		Method: "POST", // Método correto
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
		LogDetails: true, // Ativa o log para debug
	}

	// Executa o teste
	result, err := q.Qtest(opts)
	if err != nil {
		t.Fatalf("Error in Qtest: %v", err)
	}

	// Verificar se o status está correto
	if err := result.AssertStatus(StatusOK); err != nil {
		t.Errorf("Status assertion failed: %v", err)
	}

	// Verificar se o Header esperado está presente
	if err := result.AssertHeader("Content-Type", "application/json"); err != nil {
		t.Errorf("Header assertion failed: %v", err)
	}

	// Verificar se a resposta contém "Success"
	if err := result.AssertBodyContains("Success"); err != nil {
		t.Errorf("Body assertion failed: %v", err)
	}
}

// TestQTest_Options_PUT checks if the response body contains a specific substring.
// The result will TestQTest_Options_PUT(expected any) error
func TestQTest_Options_PUT(t *testing.T) {
	q := New()

	// Define the PUT route
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

// TestQTest_Options_DELETE checks if the response body contains a specific substring.
// The result will TestQTest_Options_DELETE(expected any) error
func TestQTest_Options_DELETE(t *testing.T) {
	q := New()

	// Define the DELETE route
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

// TestQTest_Options_PATCH checks if the response body contains a specific substring.
// The result will TestQTest_Options_PATCH(expected any) error
func TestQTest_Options_PATCH(t *testing.T) {
	q := New()

	// Define the PATCH route
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

// TestQTest_Options_OPTIONS checks if the response body contains a specific substring.
// The result will TestQTest_Options_OPTIONS(expected any) error
func TestQTest_Options_OPTIONS(t *testing.T) {
	q := New()

	// Define the OPTIONS route
	q.Options("/v1/user/options", func(c *Ctx) error {
		c.Set("Allow", "GET, POST, PUT, DELETE, OPTIONS")
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

	if err := result.AssertHeader("Allow", "GET, POST, PUT, DELETE, OPTIONS"); err != nil {
		t.Errorf("Header assertion failed: %v", err)
	}
}
