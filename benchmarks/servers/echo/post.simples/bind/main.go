package main

// import (
// 	"net/http"

// 	"github.com/labstack/echo/v4"
// )

// // Struct representing a user model
// type My struct {
// 	Name string `json:"name"` // User's name
// 	Year int    `json:"year"` // User's birth year
// }

// func main() {
// 	e := echo.New()

// 	// Define a POST route at /v1/user

// 	e.POST("/v1/user", func(c echo.Context) error {
// 		var my My

// 		if err := c.Bind(&my); err != nil {

// 			return c.String(http.StatusBadRequest, err.Error())
// 		}

// 		return c.JSON(http.StatusOK, my)
// 	})

// 	e.Start(":8080")
// }

// //curl --location 'http://localhost:8080/v1/user' \
// // --header 'Content-Type: application/json' \
// // --data '{"name": "Alice", "year": 20}'
