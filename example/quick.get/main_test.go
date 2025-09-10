package main

import (
	"testing"

	"github.com/jeffotoni/quick"
)

func setupApp() *quick.Quick {
	q := quick.New(quick.Config{NoBanner: true})

	q.Get("/v1/user/:name", func(c *quick.Ctx) error {
		name := c.Param("name")
		c.Set("Content-Type", "text/plain")
		return c.Status(200).SendString("Olá " + name + "!")
	})

	q.Get("/v2/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString("Opa, funcionando!")
	})

	q.Get("/v3/user/:id", func(c *quick.Ctx) error {
		id := c.Param("id")
		c.Set("Content-Type", "application/json")
		return c.Status(200).SendString("Id:" + id)
	})

	q.Get("/v1/userx/:p1/:p2/cust/:p3/:p4", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		c.Set("X-Custom-Header", "quick-value")
		c.Set("Server", "QuickServer/1.0")
		return c.Status(200).SendString("Quick in action!")
	})

	return q
}

func TestV1UserName(t *testing.T) {
	app := setupApp()

	resp, err := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/v1/user/João",
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
	if err := resp.AssertString("Olá João!"); err != nil {
		t.Error(err)
	}
}

func TestV2User(t *testing.T) {
	app := setupApp()

	resp, err := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/v2/user",
	})

	if err != nil {
		t.Fatal(err)
	}

	// Testar status
	if err := resp.AssertStatus(200); err != nil {
		t.Error(err)
	}

	// Testar header
	if err := resp.AssertHeader("Content-Type", "application/json"); err != nil {
		t.Error(err)
	}

	// Testar response
	if err := resp.AssertString("Opa, funcionando!"); err != nil {
		t.Error(err)
	}
}

func TestV3UserId(t *testing.T) {
	app := setupApp()

	resp, err := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/v3/user/123",
	})

	if err != nil {
		t.Fatal(err)
	}

	// Testar status
	if err := resp.AssertStatus(200); err != nil {
		t.Error(err)
	}

	// Testar header
	if err := resp.AssertHeader("Content-Type", "application/json"); err != nil {
		t.Error(err)
	}

	// Testar response
	if err := resp.AssertString("Id:123"); err != nil {
		t.Error(err)
	}
}

func TestV1UserxComplex(t *testing.T) {
	app := setupApp()

	resp, err := app.Qtest(quick.QuickTestOptions{
		Method: "GET",
		URI:    "/v1/userx/p1/p2/cust/p3/p4",
	})

	if err != nil {
		t.Fatal(err)
	}

	// Testar status
	if err := resp.AssertStatus(200); err != nil {
		t.Error(err)
	}

	// Testar header
	if err := resp.AssertHeader("Content-Type", "application/json"); err != nil {
		t.Error(err)
	}

	// Testar response
	if err := resp.AssertString("Quick in action!"); err != nil {
		t.Error(err)
	}

	
	// Body() e BodyStr()
	bodyBytes := resp.Body()
	t.Logf("Body()  %d bytes", len(bodyBytes))
	
	bodyStr := resp.BodyStr()
	t.Logf("BodyStr(): %s", bodyStr)
	
	// StatusCode()
	status := resp.StatusCode()
	t.Logf("StatusCode(): %d", status)
	
	// Response()
	httpResp := resp.Response()
	t.Logf("Response().Status: %s", httpResp.Status)
	
	// AssertNoHeader()
	if err := resp.AssertNoHeader("No Header"); err != nil {
		t.Logf("AssertNoHeader funcionou: %v", err)
	} else {
		t.Log("AssertNoHeader: header não existe")
	}
	
	// AssertBodyContains()
	if err := resp.AssertBodyContains("Quick"); err != nil {
		t.Error(err)
	} else {
		t.Log("AssertBodyContains: 'Quick'")
	}
	
	// AssertHeaderContains()
	if err := resp.AssertHeaderContains("Server", "Quick"); err != nil {
		t.Error(err)
	} else {
		t.Log("AssertHeaderContains: 'Quick'")
	}
	
	// AssertHeaderHasPrefix()
	if err := resp.AssertHeaderHasPrefix("X-Custom-Header", "quick"); err != nil {
		t.Error(err)
	} else {
		t.Log("AssertHeaderHasPrefix: 'quick'")
	}
	
	// AssertHeaderHasValueInSet()
	allowedValues := []string{"quick-value", "other-value", "test-value"}
	if err := resp.AssertHeaderHasValueInSet("X-Custom-Header", allowedValues); err != nil {
		t.Error(err)
	} else {
		t.Log("AssertHeaderHasValueInSet: header  está no conjunto permitido")
	}
}