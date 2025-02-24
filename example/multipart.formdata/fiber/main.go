package main

// import (
// 	"fmt"
// 	"io"
// 	"os"

// 	"github.com/gofiber/fiber/v2"
// )

// func main() {
// 	app := fiber.New()

// 	app.Post("/upload", func(c *fiber.Ctx) error {
// 		// Get the file from the request
// 		file, err := c.FormFile("file")
// 		if err != nil {
// 			return c.Status(fiber.StatusBadRequest).SendString("Error getting file")
// 		}

// 		// Open the uploaded file
// 		src, err := file.Open()
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).SendString("Error opening file")
// 		}
// 		defer src.Close()

// 		// Create a local file
// 		dst, err := os.Create(file.Filename)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).SendString("Error saving file")
// 		}
// 		defer dst.Close()

// 		// Copy the content
// 		_, err = io.Copy(dst, src)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).SendString("Error copying file")
// 		}

// 		return c.SendString(fmt.Sprintf("Upload successful: %s", file.Filename))
// 	})

// 	app.Listen(":8080")
// }
