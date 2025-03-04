package main

import (
	"log"
	"os"

	"github.com/jeffotoni/quick"
	middleware "github.com/jeffotoni/quick/middleware/basicauth"
)

// export USER=admin
// export PASSWORD=1234

var (
	User     = os.Getenv("USER")
	Password = os.Getenv("PASSORD")
)

// curl -u admin:1234 http://localhost:8080/protected
func main() {

	q := quick.New()

	q.Use(middleware.BasicAuth(User, Password))

	q.Get("/protected", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString("You have accessed a protected route!")
	})

	// Start server
	log.Fatal(q.Listen("0.0.0.0:8080"))
}
