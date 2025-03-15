package main

import (
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/limiter"
)

func main() {
	q := quick.New()

	q.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Second,
		KeyGenerator: func(c *quick.Ctx) string {
			return c.RemoteIP()
		},
		LimitReached: func(c *quick.Ctx) {
			c.Set("Content-Type", "application/json")
			c.Status(quick.StatusTooManyRequests).String(`{"error":"Too many requests"}`)
		},
	}))

	q.Get("/", func(c *quick.Ctx) error {
		return c.Status(200).JSON(map[string]string{"msg": "Hello, Quick!"})
	})

	q.Listen(":8080")
}
