package docs
//SIGNATURE func Gzip() func(next http.Handler) http.Handler

//import
// import (
// 	"log"

// 	"github.com/jeffotoni/quick"
// 	"github.com/jeffotoni/quick/middleware/compress"
// )	

// func main() {
// 	q := quick.New()

// 	// Enable GZIP compression
// 	q.Use(compress.Gzip())

// 	// Define a compressed response route
// 	q.Get("/v1/compress", func(c *quick.Ctx) error {
// 		c.Set("Content-Type", "application/json")
// 		c.Set("Accept-Encoding", "gzip")

// 		type response struct {
// 			Msg     string              `json:"msg"`
// 			Headers map[string][]string `json:"headers"`
// 		}

// 		return c.Status(200).JSON(&response{
// 			Msg:     "Quick ❤️",
// 			Headers: c.Headers,
// 		})
// 	})

// 	log.Fatal(q.Listen("0.0.0.0:8080"))
// }