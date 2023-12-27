package handler

import (
	"encoding/json"
	"fmt"
	"hackathon23/view"

	"bytes"
	"html/template"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

type messageWS struct {
	Text string `json:"message"`
}

type message struct {
	Icon string
	Name string
	Text string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]*User) // Connected clients
var broadcast = make(chan message)            // Broadcast channel

// when player first joins
func WebSocketHandler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	user, err := RandomUser()
	if err != nil {
		return err
	}

	clients[conn] = user // Add new client to the connected clients map

	var html bytes.Buffer
	err = view.ApiJoinRoom("foo", "test").Render(c.Request().Context(), &html)
	if err != nil {
		return err
	}

	// Send the HTML message to all connected clients
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, html.Bytes())
		if err != nil {
			fmt.Println("Error writing message:", err)
			client.Close()
			delete(clients, client) // Remove client from connected clients map
		}
	}

	for {
		// Read msg from client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			delete(clients, conn) // Remove client from connected clients map
			break
		}

		// if i cared about security i'd sanitize the input...

		// Parse the JSON message
		var parsedMessage messageWS
		err = json.Unmarshal(msg, &parsedMessage)
		if err != nil {
			fmt.Println("Error parsing message:", err)
			continue
		}

		// Broadcast the received message to all connected clients
		broadcast <- message{
			user.Icon,
			user.Name, // + user.Results[0].Name.Last
			parsedMessage.Text,
		}
	}

	return nil
}

func renderMessageSent(data message) (msg []byte, err error) {
	// Create a template and interpolate it with the message
	var html bytes.Buffer
	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))
	err = tmpl.ExecuteTemplate(&html, "api.message-sent", data)
	if err != nil {
		fmt.Println("Error interpolating message:", err)
		return nil, err
	}

	return html.Bytes(), nil
}

func init() {
	go broadcastMessagesToAllClient()
}

func broadcastMessagesToAllClient() {
	for {
		// Get the next msg from the broadcast channel
		msg := <-broadcast

		// Create a template and interpolate it with the message
		html, err := renderMessageSent(msg)
		if err != nil {
			fmt.Println("Error rendering message:", err)
			continue
		}

		// Send the HTML message to all connected clients
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, html)
			if err != nil {
				fmt.Println("Error writing message:", err)
				client.Close()
				delete(clients, client) // Remove client from connected clients map
			}
		}
	}
}
