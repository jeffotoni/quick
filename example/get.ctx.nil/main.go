package main

import "github.com/jeffotoni/quick"
import "github.com/jeffotoni/quick/middleware/logger"
import "github.com/jeffotoni/quick/middleware/compress"
import "github.com/jeffotoni/quick/middleware/cors"


func main() {

	q := quick.New()

	q.Use(cors.New())
	q.Use(logger.New())
	q.Use(compress.Gzip())
	
	q.Get("/v1/ping", func(c *quick.Ctx) error {
		return c.Status(200).String("pong")
	})

	q.Listen(":8080")
}
