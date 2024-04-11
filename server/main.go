package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		data := GetNotes()
		return c.JSON(http.StatusOK, data)
	})

	e.GET("/create-test", func(c echo.Context) error {
		// note := Note{
		// 	id:         "1",
		// 	title:      "Title",
		// 	content:    "description",
		// 	remindDate: "2021-01-01",
		// 	createdAt:  "2021-01-01",
		// 	updatedAt:  "2021-01-01",
		// }
		note := GetNote("2")
		return c.JSON(http.StatusOK, note)
	})

	e.Logger.Fatal(e.Start(":8000"))
}
