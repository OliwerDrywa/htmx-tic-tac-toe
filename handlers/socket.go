package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bytes"
	"html/template"

	"github.com/gorilla/websocket"
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
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer conn.Close()

	user, err := RandomUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusFailedDependency)
		return
	}

	clients[conn] = user // Add new client to the connected clients map
	html, err := renderUserJoined(templateUserJoined{user.Icon, user.Name})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendToAllClients(html)

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
}

func sendToAllClients(bytes []byte) {
	// Send the interpolated message to all connected clients
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			fmt.Println("Error writing message:", err)
			client.Close()
			delete(clients, client) // Remove client from connected clients map
		}
	}
}

type templateUserJoined struct {
	Icon string
	Name string
}

func renderUserJoined(data templateUserJoined) (msg []byte, err error) {
	var html bytes.Buffer
	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))
	err = tmpl.ExecuteTemplate(&html, "api.user-joined", data)
	if err != nil {
		fmt.Println("Error interpolating message:", err)
		return nil, err
	}

	return html.Bytes(), nil
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
	go (func() {
		for {
			// Get the next msg from the broadcast channel
			msg := <-broadcast

			// Create a template and interpolate it with the message
			html, err := renderMessageSent(msg)
			if err != nil {
				fmt.Println("Error rendering message:", err)
				continue
			}

			sendToAllClients(html)
		}
	})()
}
