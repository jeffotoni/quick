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

		// Retrieve the MsgUUID from the request headers
		_ = c.Request.Header.Get("Msguuid")

		// Return the MsgUUID in the JSON response
		return c.Status(200).JSON(map[string]string{"message": "generated msgId"})
		//return c.Status(200).JSON(map[string]string{"msguuid": msguuid})
	})

	// Send test request using Quick's built-in test utility
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

	if err := resp.AssertStatus(200); err != nil {
		fmt.Println("status error:", err)
	}

	// Print response body to verify the MsgUUID
	fmt.Println(string(resp.Body()))
	//alternative print
	//{"msguuid":"8bf38aea-d27f-4217-bea0-0549181cc26a"}

	// Output:
	// {"message":"generated msgId"}
}
