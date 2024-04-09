package main

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		data := SheetProcess()
		var raw []json.RawMessage
		err := json.Unmarshal(data, &raw)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, raw)
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
