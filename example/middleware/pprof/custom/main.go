package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/pprof"
)

func main() {
	q := quick.New()

	q.Use(pprof.New())

	q.Get("/busy", func(c *quick.Ctx) error {
		// Simulates a load
		sum := 0
		for i := 0; i < 1e7; i++ {
			sum += i
		}
		return c.String("done")
	})

	// Mandatory route for pprof to work correctly
	q.Get("/debug/pprof*", func(c *quick.Ctx) error {
		return c.Next()
	})

	q.Listen(":8080")
}
