package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Struct representing a user model
type My struct {
	Name string `json:"name"` // User's name
	Year int    `json:"year"` // User's birth year
}

func main() {

	gin.SetMode(gin.ReleaseMode)

	// r := gin.Default()
	r := gin.New()

	r.POST("/v1/user", func(c *gin.Context) {
		var my My // Create a variable to store incoming user data

		// Parse the request body into the struct
		err := c.ShouldBind(&my)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
			return
		}

		// Return the parsed JSON data as a response with 200 OK
		c.JSON(http.StatusOK, my)
	})

	r.Run(":8080")
}

//curl --location 'http://localhost:8080/v1/user' \
// --header 'Content-Type: application/json' \
// --data '{"name": "Alice", "year": 20}'
