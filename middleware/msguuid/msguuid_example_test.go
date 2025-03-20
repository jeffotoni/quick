package msguuid

import (
	"fmt"

	"github.com/jeffotoni/quick"
)

// This function is named ExampleNew()
//
//	it with the Examples type.
func ExampleNew() {
	q := quick.New()

	// Apply MsgUUID Middleware globally
	q.Use(New())

	// Define an endpoint that responds with a UUID
	q.Get("/v1/msguuid/default", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// Log headers to validate UUID presence
		fmt.Println("Headers:", c.Response.Header())

		// Return a 200 OK status
		return c.Status(200).JSON(nil)
	})

	// Send test request using Quick's built-in test utility
	_, _ = q.Qtest(quick.QuickTestOptions{
		Method:  quick.MethodGet,
		URI:     "/v1/msguuid/default",
		Headers: map[string]string{"Content-Type": "application/json"},
	})

	// Out put: null
	// Console: Headers: map[Content-Type:[application/json] Msguuid:[f299b00d-875e-4502-966e-22e16767eb13]]
}
