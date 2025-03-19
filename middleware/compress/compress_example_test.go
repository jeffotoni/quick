package compress

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick"
)

// This function is named ExampleGzip()
// it with the Examples type.
func ExampleGzip() {
	// Starting Quick framework instance
	q := quick.New()

	// Enable Gzip middleware
	// This will automatically compress responses for clients that support Gzip
	q.Use(Gzip())

	// Define a route that returns a compressed JSON response
	q.Get("/v1/compress", func(c *quick.Ctx) error {
		// Setting response headers
		c.Set("Content-Type", "application/json")

		// Defining the response structure
		type response struct {
			Msg     string              `json:"msg"`
			Headers map[string][]string `json:"headers"`
		}

		// Returning a JSON response with headers
		return c.Status(200).JSON(&response{
			Msg:     "Quick ❤️",
			Headers: c.Headers,
		})
	})

	// Simulate a GET request with headers using Quick's testing functionality
	res, err := q.Qtest(quick.QuickTestOptions{
		Method:  quick.MethodGet,
		URI:     "/v1/compress",
		Headers: map[string]string{"Accept-Encoding": "gzip"},
	})
	if err != nil {
		log.Fatalf("Error running test request: %v", err)
	}

	// Print the response body
	fmt.Println(res.BodyStr())

	// Out put:
	// {"msg":"Quick ❤️","headers":{"Accept-Encoding":["gzip"]}}
}
