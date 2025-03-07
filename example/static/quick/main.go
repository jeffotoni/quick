//go:build !exclude_test

package main

import "github.com/jeffotoni/quick"

func main() {

	// start Quick
	q := quick.New()

	// start dir files
	q.Static("/static", "./static")

	// server files
	q.Get("/", func(c *quick.Ctx) error {
		c.File("./static/index.html")
		return nil
	})

	q.Listen("0.0.0.0:8080")
}

// $ curl --location 'http://localhost:8080/'
