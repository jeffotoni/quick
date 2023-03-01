package main

import (
	"log"

	"github.com/gojeffotoni/quick"
)

func main() {
	app := quick.New()
	app.Get("/some", QuickHandlerGet)
	log.Printf("dataaaa: %v", app.GetRoute())
	log.Fatal(app.Listen("0.0.0.0:8080"))
}

func QuickHandlerGet(c *quick.Ctx) {
	c.Set("Content-Type", "application/json")
	log.Printf("params: %v", c.Params)
}
