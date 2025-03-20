package msgid

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

		msgId := c.Request.Header.Get("Msgid")

		// Return 200 OK status
		return c.Status(200).JSON(map[string]string{"msgid": msgId})
	})

	// Send a test request using Quick's built-in test utility
	resp, err := q.Qtest(quick.QuickTestOptions{
		Method:  quick.MethodGet,
		URI:     "/v1/msguuid/default",
		Headers: map[string]string{"Content-Type": "application/json"},
	})

	// Handle potential errors in test execution
	if err != nil {
		fmt.Println("Test execution error:", err)
		return
	}

	// Print the response body for verification
	fmt.Println("Response:", string(resp.Body()))

	// Out put:
	// Response: {"msgid":"f299b00d-875e-4502-966e-22e16767eb13"}
}
