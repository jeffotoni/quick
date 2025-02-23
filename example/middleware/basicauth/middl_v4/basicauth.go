package main

import (
	"log"

	"github.com/jeffotoni/quick"
	middleware "github.com/jeffotoni/quick/middleware/basicauth"
)

// curl -u admin:1234 http://localhost:8080/protected
// or
// curl -H "Authorization: Basic YWRtaW46MTIzNA==" http://localhost:8080/protected
func main() {

	q := quick.New()

	// using group to isolate routes and middlewares
	gr := q.Group("/")

	// middleware BasicAuth
	gr.Use(middleware.BasicAuth("admin", "1234"))

	// route public
	q.Get("/v1/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString("Public quick route")
	})

	// protected route
	gr.Get("/protected", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString("You have accessed a protected route!")
	})

	// Start server
	log.Fatal(q.Listen("0.0.0.0:8080"))
}
