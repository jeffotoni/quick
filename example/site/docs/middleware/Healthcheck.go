package docs

//Signature func New(opts Options) func(next http.Handler) http.Handler

// import (
// 	"github.com/jeffotoni/quick"
// 	"github.com/seuusuario/healthcheck"
// )

// func main() {
// 	q := quick.New()

// 	 Use Healthcheck middleware with default healthcheck endpoint
// 	q.Use(healthcheck.New(
// 		healthcheck.Options{
// 			App: q,
// 		},
// 	))

// 	q.Get("/", func(c *quick.Ctx) error {
// 		return c.Status(200).String("Home page")
// 	})

// 	log.Fatalln(q.Listen(":8080"))
// }