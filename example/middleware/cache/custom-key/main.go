// Example of cache middleware with custom key generation in Quick
package main

import (
	"fmt"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/cache"
)

func main() {
	// Create a new Quick q
	q := quick.New()

	// Use the cache middleware with custom key generation
	q.Use(cache.New(cache.Config{
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *quick.Ctx) string {
			return c.Path() + "?lang=" + c.Query["lang"] + "&user=" + c.Query["user"]
		},
		CacheHeader:          "X-Cache-Status",
		StoreResponseHeaders: true,
		Methods:              []string{quick.MethodGet},
		CacheInvalidator: func(c *quick.Ctx) bool {
			return c.Query["clear"] == "1"
		},
	}))

	// Route that returns a greeting in the requested language
	q.Get("/greeting", func(c *quick.Ctx) error {
		lang := c.Query["lang"]
		user := c.Query["user"]
		if user == "" {
			user = "Guest"
		}

		greeting := "Hello"
		switch lang {
		case "pt":
			greeting = "Olá"
		case "es":
			greeting = "Hola"
		case "fr":
			greeting = "Bonjour"
		case "it":
			greeting = "Ciao"
		case "de":
			greeting = "Hallo"
		}

		return c.String(fmt.Sprintf("%s, %s! (Generated at %s)",
			greeting, user, time.Now().Format(time.RFC3339)))
	})

	// Start the server
	fmt.Println("Server running on http://localhost:3000")
	fmt.Println("Try these examples:")
	fmt.Println("  - GET /greeting?lang=en&user=John")
	fmt.Println("  - GET /greeting?lang=pt&user=Maria")
	fmt.Println("  - GET /greeting?lang=es&user=Carlos")
	fmt.Println("  - GET /greeting?lang=en&user=John (should be cached)")
	fmt.Println("  - GET /greeting?lang=en&user=John&clear=1 (invalidates cache)")
	fmt.Println("Check the X-Cache-Status header in the response to see if it's a HIT or MISS")
	q.Listen(":3000")
}
