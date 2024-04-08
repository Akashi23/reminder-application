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
	e.Logger.Fatal(e.Start(":8000"))
}
