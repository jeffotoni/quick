package main

// import (
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"

// 	"github.com/labstack/echo/v4"
// )

// func main() {
// 	e := echo.New()

// 	e.POST("/upload", func(c echo.Context) error {
// 		// Get the file from the request
// 		file, err := c.FormFile("file")
// 		if err != nil {
// 			return c.String(http.StatusBadRequest, "Error getting file")
// 		}

// 		// Open the uploaded file
// 		src, err := file.Open()
// 		if err != nil {
// 			return c.String(http.StatusInternalServerError, "Error opening file")
// 		}
// 		defer src.Close()

// 		// Create a local file
// 		dst, err := os.Create(file.Filename)
// 		if err != nil {
// 			return c.String(http.StatusInternalServerError, "Error saving file")
// 		}
// 		defer dst.Close()

// 		// Copy the content
// 		_, err = io.Copy(dst, src)
// 		if err != nil {
// 			return c.String(http.StatusInternalServerError, "Error copying file")
// 		}

// 		return c.String(http.StatusOK, fmt.Sprintf("Upload successful: %s", file.Filename))
// 	})

// 	e.Start(":8080")
// }
