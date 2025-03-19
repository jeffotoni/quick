package main

import (
	"log"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/middleware/maxbody"
)

func main() {
	q := quick.New()

	q.Use(maxbody.New(0))

	q.Post("/v1/user/maxbody/large", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")

		log.Printf("body: %s", c.BodyString())
		return c.Status(200).Send(c.Body())
	})

	log.Fatal(q.Listen("0.0.0.0:8080"))
}

// $ curl -i -XPOST http://0.0.0.0:8080/v1/user/maxbody/large -d '{"data":"quick is awesome!"}'
