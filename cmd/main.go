package main

import (
	"fmt"
	"hackathon23/handler"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.Static("/static", "web")
	e.GET("/", handler.IndexHandler)
	e.GET("/ws", handler.WSHandler)

	fmt.Println("listening on http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}
