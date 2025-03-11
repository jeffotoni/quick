package main

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Struct representing a user model
type My struct {
	Name string `json:"name"` // User's name
	Year int    `json:"year"` // User's birth year
}

func main() {
	e := echo.New()

	// Define a POST route at /v1/user

	e.POST("/v1/user", func(c echo.Context) error {
		var my My

		_, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Erro ao ler o body: "+err.Error())
		}

		return c.JSON(http.StatusOK, my)
	})

	e.Start(":8080")
}
