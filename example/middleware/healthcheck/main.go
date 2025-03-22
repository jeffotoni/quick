package main

import (
	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/healthcheck"
)

func main() {
	q := quick.New()

	// register healthcheck middleware
	q.Use(healthcheck.New())

	q.Listen(":8080")
}
