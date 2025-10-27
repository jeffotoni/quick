package main

import (
	"testing"

	"github.com/jeffotoni/quick"
)

func setupApp() *quick.Quick {
	q := quick.New(quick.Config{NoBanner: true})

	// Register a HEAD route
	q.Head("/ping", func(c *quick.Ctx) error {
		c.Set("Content-Type", "text/plain")
		c.Set("X-Head-Route", "true")
		return c.String("pong!") // Will not be included in response body for HEAD
	})

	// Optional: Register a GET to show difference
	q.Get("/ping", func(c *quick.Ctx) error {
		c.Set("Content-Type", "text/plain")
		c.Set("X-Get-Route", "true")
		return c.String("pong!")
	})

	return q
}

func TestHeadPing(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "HEAD",
		URI:    "/ping",
	})

	resp.AssertStatus(200)
	resp.AssertHeader("Content-Type", "text/plain")
	resp.AssertHeader("X-Head-Route", "true")
}

func TestGetPing(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/ping",
	})

	resp.AssertStatus(200)
	resp.AssertHeader("Content-Type", "text/plain")
	resp.AssertHeader("X-Get-Route", "true")
	resp.AssertString("pong!")
}

func TestBodyMethods(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/ping",
	})

	body := resp.Body()
	if len(body) == 0 {
		t.Error("Body vazio")
	}

	bodyStr := resp.BodyStr()
	if bodyStr != "pong!" {
		t.Error("BodyStr incorreto")
	}
}

func TestStatusCode(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "HEAD",
		URI:    "/ping",
	})

	if resp.StatusCode() != 200 {
		t.Error("StatusCode incorreto")
	}
}

func TestResponse(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/ping",
	})

	if resp.Response() == nil {
		t.Error("Response nil")
	}
}

func TestAssertNoHeader(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/ping",
	})

	resp.AssertNoHeader("X-Invalid-Header")
}

func TestAssertBodyContains(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/ping",
	})

	resp.AssertBodyContains("pong")
}

func TestAssertHeaderContains(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/ping",
	})

	resp.AssertHeaderContains("Content-Type", "text")
}

func TestAssertHeaderHasPrefix(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/ping",
	})

	resp.AssertHeaderHasPrefix("Content-Type", "text")
}

func TestAssertHeaderHasValueInSet(t *testing.T) {
	app := setupApp()
	resp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/ping",
	})

	allowed := []string{"text/plain", "application/json"}
	resp.AssertHeaderHasValueInSet("Content-Type", allowed)
}

func TestHeadVsGetDifference(t *testing.T) {
	app := setupApp()
	
	// Test HEAD route
	headResp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "HEAD",
		URI:    "/ping",
	})

	// Test GET route  
	getResp, _ := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/ping",
	})

	// HEAD não deve ter corpo
	if len(headResp.Body()) > 0 {
		t.Error("HEAD response should not have body")
	}

	// GET deve ter corpo
	if len(getResp.Body()) == 0 {
		t.Error("GET response should have body")
	}

	// Verificar headers específicos de cada rota
	headResp.AssertHeader("X-Head-Route", "true")
	getResp.AssertHeader("X-Get-Route", "true")
}