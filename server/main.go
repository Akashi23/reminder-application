package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		return key == "123", nil
	}))

	e.GET("/", func(c echo.Context) error {
		data := GetNotes()
		return c.JSON(http.StatusOK, data)
	})

	e.GET("/notes", func(c echo.Context) error {
		data := GetNotes()
		return c.JSON(http.StatusOK, data)
	})

	e.POST("/notes", func(c echo.Context) error {
		note := new(Note)
		if err := c.Bind(note); err != nil {
			return err
		}

		note.CreatedAt = time.Now().String()
		note.UpdatedAt = time.Now().String()

		CreateNote(*note)

		return c.JSON(http.StatusCreated, note)
	})

	e.GET("/notes/:id", func(c echo.Context) error {
		id := c.Param("id")
		data := GetNote(id)
		return c.JSON(http.StatusOK, data)
	})

	e.PUT("/notes/:id", func(c echo.Context) error {
		id := c.Param("id")
		note := new(Note)
		if err := c.Bind(note); err != nil {
			return err
		}

		note.Id = id

		note.UpdatedAt = time.Now().String()

		UpdateNote(*note)

		return c.JSON(http.StatusOK, note)

	})

	e.DELETE("/notes/:id", func(c echo.Context) error {
		id := c.Param("id")
		DeleteNote(id)
		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start(":8000"))
}
