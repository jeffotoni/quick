package healthcheck

import (
	"fmt"

	"github.com/jeffotoni/quick"
)

// ExampleNew demonstrates how to use the healthcheck middleware
// with default options, registering a /healthcheck endpoint.
func ExampleNew() {
	q := quick.New()

	// Use healthcheck middleware with default options
	q.Use(New(Options{
		App: q,
	}))

	resp, _ := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/healthcheck",
	})

	fmt.Println("Status:", resp.StatusCode())
	fmt.Println("Body:", resp.BodyStr())

	// Output:
	// Status: 200
	// Body: OK
}

// ExampleNew_customEndpoint shows how to configure a custom healthcheck
// endpoint and define a custom health Probe function.
func ExampleNew_customEndpoint() {
	q := quick.New()

	// Use healthcheck with a custom endpoint and probe logic
	q.Use(New(Options{
		App:      q,
		Endpoint: "/v1/health",
		Probe: func(c *quick.Ctx) bool {
			return true // simulate healthy state
		},
	}))

	resp, _ := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/v1/health",
	})

	fmt.Println("Status:", resp.StatusCode())
	fmt.Println("Body:", resp.BodyStr())

	// Output:
	// Status: 200
	// Body: OK
}

// ExampleNew_unhealthyProbe simulates an unhealthy service
// by returning false in the Probe function.
func ExampleNew_unhealthyProbe() {
	q := quick.New()

	// Use healthcheck where Probe simulates an unhealthy service
	q.Use(New(Options{
		App: q,
		Probe: func(c *quick.Ctx) bool {
			return false // simulate failure
		},
	}))

	resp, _ := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/healthcheck",
	})

	fmt.Println("Status:", resp.StatusCode())
	fmt.Println("Body:", resp.BodyStr())

	// Output:
	// Status: 503
	// Body: Service Unavailable
}
