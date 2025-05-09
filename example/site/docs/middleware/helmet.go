package docs

//Signature func Helmet() func(next http.Handler) http.Handler

// import (
// 	"github.com/jeffotoni/quick"
// 	"github.com/seuusuario/helmet"
	
// )

// func main() {
// 	q := quick.New()

// 	// Use Helmet middleware with default security headers
// 	q.Use(helmet.Helmet())

// 	// Simple route to test headers
// 	q.Get("/v1/user", func(c *quick.Ctx) error {

// 		// list all headers
// 		headers := make(map[string]string)
// 		for k, v := range c.Response.Header() {
// 			if len(v) > 0 {
// 				headers[k] = v[0]
// 			}
// 		}
// 		return c.Status(200).JSONIN(headers)
// 	})

// 	q.Listen("0.0.0.0:8080")
// }