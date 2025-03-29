package maxbody

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jeffotoni/quick"
)

// This function is named ExampleNew()
//
//	it with the Examples type.
func ExampleNew() {
	// Define a maximum body size of 1KB (1024 bytes)
	const maxBodySize = 1024

	// Create a new Quick instance
	q := quick.New()

	// Apply the maxbody middleware to limit request size
	q.Use(New(maxBodySize))

	// Define a test route
	q.Post("/v1/user/maxbody", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		// Read the request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			return c.Status(http.StatusRequestEntityTooLarge).String("Request body too large")
		}

		// Return the received body
		return c.Status(http.StatusOK).Send(body)
	})

	// Simulate an HTTP request exceeding the 1KB limit
	oversizedBody := make([]byte, 2048) // 2KB payload to exceed the limit
	for i := range oversizedBody {
		oversizedBody[i] = 'A'
	}

	// Send test request using Quick's built-in test utility
	res, _ := q.Qtest(quick.QuickTestOptions{
		Method:  quick.MethodPost,
		URI:     "/v1/user/maxbody",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    oversizedBody, // Exceeding body limit
	})

	// Print the response status and body
	fmt.Println(res.StatusCode()) // Expected: 413 (Payload Too Large)
	fmt.Println(res.BodyStr())    // Expected: "Request body too large"

	// Output:
	// 413
	// Request body too large
}
