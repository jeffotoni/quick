package docs

// func New(config Config) func(next http.Handler) http.Handler

// import (
// 	"time"
// 	"github.com/jeffotoni/quick"
// 	"github.com/jeffotoni/quick/middleware/limiter"
// )



// func main() {
// 	q := quick.New()

// 	// Apply the rate limiter middleware
// 	q.Use(limiter.New(limiter.Config{
// 		// Maximum 10 requests allowed per IP
// 		Max: 10,
// 		// The limit resets after 5 seconds
// 		Expiration: 5 * time.Second,
// 		KeyGenerator: func(c *quick.Ctx) string {
// 			// Uses the client's IP address as the key
// 			return c.RemoteIP()
// 		},
// 		LimitReached: func(c *quick.Ctx) error {
// 			c.Set("Content-Type", "application/json")
// 			// The client should wait 10 seconds before retrying
// 			c.Set("Retry-After", "10")
// 			return c.Status(quick.StatusTooManyRequests).JSON(map[string]string{
// 				"error":   "Too many requests",
// 				"message": "You have exceeded the request limit. 
// 				Please wait 1 second and try again.",
// 				"retry_after": "10s",
// 			})
// 		},
// 	}))

// 	// Define a simple GET route
// 	q.Get("/", func(c *quick.Ctx) error {
// 		return c.Status(200).JSON(map[string]string{"msg": "Quick in action ❤️!"})
// 	})

// 	// Start the server on port 8080
// 	q.Listen(":8080")
// }

