package main

// import (
// 	"fmt"
// 	"io"
// 	"os"

// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	r := gin.Default()

// 	r.POST("/upload", func(c *gin.Context) {
// 		// Get the file from the request
// 		file, err := c.FormFile("file")
// 		if err != nil {
// 			c.String(400, "Error getting file")
// 			return
// 		}

// 		// Open the uploaded file
// 		src, err := file.Open()
// 		if err != nil {
// 			c.String(500, "Error opening file")
// 			return
// 		}
// 		defer src.Close()

// 		// Create a local file
// 		dst, err := os.Create(file.Filename)
// 		if err != nil {
// 			c.String(500, "Error saving file")
// 			return
// 		}
// 		defer dst.Close()

// 		// Copy the content
// 		_, err = io.Copy(dst, src)
// 		if err != nil {
// 			c.String(500, "Error copying file")
// 			return
// 		}

// 		c.String(200, fmt.Sprintf("Upload successful: %s", file.Filename))
// 	})

// 	r.Run(":8080")
// }
