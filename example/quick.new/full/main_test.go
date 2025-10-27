package main

import (
	"net/http"
	"testing"

	"github.com/jeffotoni/quick"
)

func setupApp() *quick.Quick {
	config := quick.Config{
		RouteCapacity: 500,
		NoBanner:      true,
	}
	q := quick.New(config)

	// Middleware correto
	q.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware", "executed")
			next.ServeHTTP(w, r)
		})
	})

	// Define multiple routes
	q.Get("/", func(c *quick.Ctx) error {
		c.Set("Content-Type", "text/plain")
		return c.Status(200).String("Welcome to Quick!")
	})

	q.Post("/data", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(201).String("Data received!")
	})

	q.Put("/update", func(c *quick.Ctx) error {
		c.Set("Content-Type", "text/plain")
		return c.Status(200).String("Data updated!")
	})

	q.Delete("/delete", func(c *quick.Ctx) error {
		return c.Status(204).String("") // No content response
	})

	return q
}

func TestGetRoot(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/",
	})

	resp.AssertStatus(200)
	resp.AssertHeader("Content-Type", "text/plain")
	resp.AssertString("Welcome to Quick!")
}

func TestPostData(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "POST",
		URI:    "/data",
	})

	resp.AssertStatus(201)
	resp.AssertHeader("Content-Type", "application/json")
	resp.AssertString("Data received!")
}

func TestPutUpdate(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "PUT",
		URI:    "/update",
	})

	resp.AssertStatus(200)
	resp.AssertHeader("Content-Type", "text/plain")
	resp.AssertString("Data updated!")
}

func TestDelete(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "DELETE",
		URI:    "/delete",
	})

	resp.AssertStatus(204)
	resp.AssertString("")
}

func TestBodyMethods(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/",
	})

	body := resp.Body()
	if len(body) == 0 {
		t.Error("Body vazio")
	}

	bodyStr := resp.BodyStr()
	if bodyStr != "Welcome to Quick!" {
		t.Error("BodyStr incorreto")
	}
}

func TestStatusCode(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "POST",
		URI:    "/data",
	})

	if resp.StatusCode() != 201 {
		t.Error("StatusCode incorreto")
	}
}

func TestResponse(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "PUT",
		URI:    "/update",
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

	resp.AssertNoHeader("X-Invalid-Header")
}

func TestAssertBodyContains(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/",
	})

	resp.AssertBodyContains("Welcome")
}

func TestAssertHeaderContains(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/",
	})

	resp.AssertHeaderContains("Content-Type", "text")
}

func TestAssertHeaderHasPrefix(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "POST",
		URI:    "/data",
	})

	resp.AssertHeaderHasPrefix("Content-Type", "application")
}

func TestAssertHeaderHasValueInSet(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/",
	})

	allowed := []string{"text/plain", "application/json"}
	resp.AssertHeaderHasValueInSet("Content-Type", allowed)
}

func TestMiddleware(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/",
	})

	resp.AssertHeader("X-Middleware", "executed")
}