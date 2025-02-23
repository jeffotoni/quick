package main

// import (
// 	"fmt"

// 	"github.com/gofiber/fiber/v2"
// )

// func main() {
// 	app := fiber.New()

// 	app.Static("/static", "./static")

// 	app.Get("/", func(c *fiber.Ctx) error {
// 		return c.SendFile("./static/index.html")
// 	})

// 	port := 8080
// 	fmt.Printf("Server Run http://localhost:%d\n", port)
// 	app.Listen(fmt.Sprintf(":%d", port))
// }
