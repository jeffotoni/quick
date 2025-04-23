// The BasicAuth middleware implements HTTP Basic Authentication
// to secure specific routes on an HTTP server.
// Example of how to use middleware in Quick
// $ curl -H "Authorization: Basic $(echo -n 'wronguser:wrongpass' | base64)" http://localhost:8080/protected
// $ curl -H "Authorization: Basic $(echo -n 'admin:1234' | base64)" http://localhost:8080/protected
// $ curl http://localhost:8080/protected
package basicauth

import (
	"log"

	"github.com/jeffotoni/quick"
)

// This function is named ExampleBasicAuth()
// it with the Examples type.
func ExampleBasicAuth() {
	//starting Quick
	q := quick.New()

	// calling middleware
	q.Use(BasicAuth("admin", "1234"))

	// everything below Use will apply the middleware
	q.Get("/protected", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString("You have accessed a protected route!")
	})

	// Start server
	log.Fatal(q.Listen("0.0.0.0:8080"))
}

//This function is named ExampleBasicAuth_withGroup
func ExampleBasicAuth_withGroup() {
	//starting Quick
	q := quick.New()

	// using group to isolate routes and middlewares
	gr := q.Group("/")

	// middleware BasicAuth
	gr.Use(BasicAuth("admin", "1234"))

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
