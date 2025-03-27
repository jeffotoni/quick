package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/pprof"
)

func main() {
	// Initialize a new Quick application
	q := quick.New()

	// Enable the pprof middleware for runtime profiling
	// This provides access to routes like /debug/pprof/heap, /goroutine, /profile, etc.
	q.Use(pprof.New())

	// Simulated workload endpoint to help generate profiling data
	q.Get("/busy", func(c *quick.Ctx) error {
		// Simulate CPU load for profiling
		sum := 0
		for i := 0; i < 1e7; i++ {
			sum += i
		}
		return c.String("done")
	})

	// Mandatory route to forward matching requests to the pprof middleware
	// This allows the pprof handler to respond properly under /debug/pprof*
	q.Get("/debug/pprof*", func(c *quick.Ctx) error {
		return c.Next()
	})

	// Start the server on port 8080
	q.Listen(":8080")
}

// $ curl http://localhost:8080/busy
// $ curl http://localhost:8080/debug/pprof/
