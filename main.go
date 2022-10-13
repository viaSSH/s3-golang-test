package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {

	fmt.Println("hi")

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "health",
		})
	})

	e.GET("/:key", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "health",
		})
	})

	e.Logger.Debug(e.Start(":8090"))
}
