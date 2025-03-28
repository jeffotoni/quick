package pprof

import (
	"fmt"

	"github.com/jeffotoni/quick"
)

// This function is named ExampleNew()
//
//	it with the Examples type.
func ExampleNew() {

	// Create a new Quick application instance
	q := quick.New()

	// Apply the pprof middleware to enable runtime profiling
	// This allows access to profiling endpoints like /debug/pprof/heap, /goroutine, etc.
	q.Use(New())

	// Define a test route that matches /debug/pprof*
	// This is required so that the Quick router delegates the request to the pprof middleware
	q.Get("/debug/pprof*", func(c *quick.Ctx) error {
		return c.Next()
	})

	// Simulate a GET request with headers
	res, _ := q.Qtest(quick.QuickTestOptions{
		Method: quick.MethodGet,
		URI:    "/debug/pprof",
	})

	if err := res.AssertStatus(200); err != nil {
		fmt.Println("Status error:", err)
	}

	fmt.Println("Status:", res.StatusCode())

	// Output: Status: 200

}
