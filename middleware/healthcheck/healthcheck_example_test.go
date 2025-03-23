package healthcheck

import (
	"log"

	"github.com/jeffotoni/quick"
)

func Example() {
	q := quick.New()
	q.Use(New(
		Options{
			App: q,
		},
	))

	q.Get("/", func(c *quick.Ctx) error {
		return c.Status(200).String("Home page")
	})

	log.Fatalln(q.Listen(":8080"))
}
