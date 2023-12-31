package handler

import (
	"fmt"
	"hackathon23/web/components"
	"hackathon23/web/views"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var clients = make(map[*websocket.Conn]string) // Connected clients
var broadcast = make(chan []byte)              // Broadcast channel

var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

// when player first joins
func WebSocketHandler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer onWebSocketClosed(c, conn)

	// event - ws_connected
	html, err := render(c, views.WebSocketConnected())
	if err != nil {
		return err
	}
	err = send(conn, html)
	if err != nil {
		return err
	}

	// event - user_joined
	var msg0 struct {
		UserName string `json:"user-name"`
	}
	err = receive(conn, &msg0)
	if err != nil {
		return err
	}
	username := msg0.UserName
	// TODO - validate the UserName isn't empty on server too

	// Add new client to the connected clients map
	clients[conn] = username

	html, err = render(c, views.YouJoinedRoom())
	if err != nil {
		return err
	}
	err = send(conn, html)
	if err != nil {
		return err
	}

	// Send the "X user joined" message to all connected clients
	html, err = render(c, components.UserJoinedTheRoomMessage(username))
	if err != nil {
		return err
	}
	broadcast <- html
	html, err = render(c, components.GuestList(clients))
	if err != nil {
		return err
	}
	broadcast <- html

	for {
		// Read msg from client
		var msg struct {
			ChatMessage string `json:"chat-message"`
		}
		err = receive(conn, &msg)
		if err != nil {
			return err
		}

		html, err = render(c, components.NewMessage(username, msg.ChatMessage))
		if err != nil {
			return err
		}

		// Broadcast the received message to all connected clients
		broadcast <- html

		// clear chat box
		html, err = render(c, components.EmptyChatInput())
		if err != nil {
			return err
		}

		err = send(conn, html)
		if err != nil {
			return err
		}
	}
}

func onWebSocketClosed(c echo.Context, conn *websocket.Conn) error {
	conn.Close()
	delete(clients, conn) // Remove client from connected clients map

	html, err := render(c, components.UserLeftTheRoomMessage(clients[conn]))
	if err != nil {
		return err
	}
	broadcast <- html
	html, err = render(c, components.GuestList(clients))
	if err != nil {
		return err
	}
	broadcast <- html

	return nil
}

func init() {
	go (func() {
		for {
			html := <-broadcast

			for conn := range clients {
				err := send(conn, html)
				if err != nil {
					fmt.Println("Error writing message:", err)
					conn.Close()
					delete(clients, conn) // Remove client from connected clients map
				}
			}
		}
	})()
}
