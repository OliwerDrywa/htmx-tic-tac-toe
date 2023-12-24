package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan string)            // Broadcast channel

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer conn.Close()
	clients[conn] = true // Add new client to the connected clients map

	for {
		// Read message from client
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			delete(clients, conn) // Remove client from connected clients map
			break
		}

		// if i cared about security i'd sanitize the input...

		// Parse the JSON message
		var parsedMessage Message
		err = json.Unmarshal(message, &parsedMessage)
		if err != nil {
			fmt.Println("Error parsing message:", err)
			continue
		}

		// Broadcast the received message to all connected clients
		broadcast <- parsedMessage.Message
	}
}

func handleMessages() {
	for {
		// Get the next message from the broadcast channel
		message := <-broadcast

		// Create a template and interpolate it with the message
		var messageHtml bytes.Buffer
		tmpl := template.Must(template.ParseGlob("web/templates/*.html"))
		err := tmpl.ExecuteTemplate(&messageHtml, "message.fragment", struct{ Message string }{message})
		if err != nil {
			fmt.Println("Error interpolating message:", err)
			continue
		}

		// Send the interpolated message to all connected clients
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, messageHtml.Bytes())
			if err != nil {
				fmt.Println("Error writing message:", err)
				client.Close()
				delete(clients, client) // Remove client from connected clients map
			}
		}
	}
}

func init() {
	go handleMessages()
}
