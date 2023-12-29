package handler

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/a-h/templ"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

func render(c echo.Context, component templ.Component) ([]byte, error) {
	var html bytes.Buffer

	err := component.Render(c.Request().Context(), &html)
	if err != nil {
		return html.Bytes(), err
	}

	return html.Bytes(), nil
}

func send(conn *websocket.Conn, html []byte) (err error) {
	err = conn.WriteMessage(websocket.TextMessage, html)
	if err != nil {
		return err
	}

	return nil
}

func receive(conn *websocket.Conn, v interface{}) (err error) {
	// Read msg from client
	_, rawMsg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Error reading message:", err)
		return err
	}

	// TODO - if i cared about security i'd sanitize the input...

	// Parse the JSON message
	err = json.Unmarshal(rawMsg, v)
	if err != nil {
		fmt.Println("Error parsing message:", err)
		return err
	}

	return nil
}
