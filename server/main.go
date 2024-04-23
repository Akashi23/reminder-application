package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "query:api-key",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == GetEnv("API_KEY"), nil
		},
	}))
	e.GET("/ws", sendRemind)

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

		newNote := CreateNote(*note)

		return c.JSON(http.StatusCreated, newNote)
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

		updatedNote := UpdateNote(*note)

		return c.JSON(http.StatusOK, updatedNote)

	})

	e.DELETE("/notes/:id", func(c echo.Context) error {
		id := c.Param("id")
		DeleteNote(id)
		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start(":" + GetEnv("PORT")))
}

var (
	upgrader = websocket.Upgrader{}
)

func sendRemind(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Connected to server!"))
		if err != nil {
			c.Logger().Error(err)
		}

		go CheckRemind(ws, c)

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}

func CheckRemind(ws *websocket.Conn, c echo.Context) {
	for {
		notes := GetNotes()
		for _, note := range notes {
			if note.RemindDate == time.Now().Format("2006/01/02") {
				sendNote, err := json.Marshal(note)
				if err != nil {
					fmt.Println(err)
				} else {
					err := ws.WriteMessage(websocket.TextMessage, []byte(string(sendNote)))
					if err != nil {
						c.Logger().Error(err)
					}
				}
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
