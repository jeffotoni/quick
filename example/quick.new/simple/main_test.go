package main

import (
	"testing"

	"github.com/jeffotoni/quick"
)

func setupApp() *quick.Quick {
	q := quick.New(quick.Config{NoBanner: true})

	// Define a simple GET route at the root path
	q.Get("/", func(c *quick.Ctx) error {
		// Set response header to indicate plain text response
		c.Set("Content-Type", "text/plain")
		c.Set("Server", "QuickServer")
		c.Set("X-Version", "1.0.0")

		// Return a 200 OK response with a message
		return c.Status(200).String("Quick in action!")
	})

	return q
}

func TestRootRoute(t *testing.T) {
	app := setupApp()

	resp, err := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/",
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

	// Testar response
	if err := resp.AssertString("Quick in action!"); err != nil {
		t.Error(err)
	}
}

func TestBody(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/",
	})

	body := resp.Body()
	if len(body) == 0 {
		t.Error("Body vazio")
	}
}

func TestBodyStr(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/",
	})

	bodyStr := resp.BodyStr()
	if bodyStr != "Quick in action!" {
		t.Error("BodyStr incorreto")
	}
}

func TestStatusCode(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/",
	})

	if resp.StatusCode() != 200 {
		t.Error("StatusCode incorreto")
	}
}

func TestResponse(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/",
	})

	if resp.Response() == nil {
		t.Error("Response nil")
	}
}

func TestAssertNoHeader(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/",
	})

	if err := resp.AssertNoHeader("X-Invalid"); err != nil {
		t.Error(err)
	}
}

func TestAssertBodyContains(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/",
	})

	if err := resp.AssertBodyContains("action"); err != nil {
		t.Error(err)
	}
}

func TestAssertHeaderContains(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/",
	})

	if err := resp.AssertHeaderContains("Server", "Quick"); err != nil {
		t.Error(err)
	}
}

func TestAssertHeaderHasPrefix(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/",
	})

	if err := resp.AssertHeaderHasPrefix("Content-Type", "text"); err != nil {
		t.Error(err)
	}
}

func TestAssertHeaderHasValueInSet(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/",
	})

	allowed := []string{"text/plain", "application/json"}
	if err := resp.AssertHeaderHasValueInSet("Content-Type", allowed); err != nil {
		t.Error(err)
	}
}

func TestAllAssertions(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/",
	})

	// Testar tudo junto
	resp.AssertStatus(200)
	resp.AssertHeader("Content-Type", "text/plain")
	resp.AssertNoHeader("X-Fake-Header")
	resp.AssertString("Quick in action!")
	resp.AssertBodyContains("Quick")
	resp.AssertHeaderContains("Server", "Quick")
	resp.AssertHeaderHasPrefix("X-Version", "1.0")
	
	allowed := []string{"1.0.0", "2.0.0"}
	resp.AssertHeaderHasValueInSet("X-Version", allowed)
}