package wss

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var Server = WSServer{make(map[*websocket.Conn]WSClient), make(chan []byte)}

func init() {
	go broadcast()
}

type WSServer struct {
	Clients     map[*websocket.Conn]WSClient
	Broadcaster chan []byte
}
type WSClient struct {
	c    echo.Context
	conn *websocket.Conn
	Name string
	Role int
}

var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func (server WSServer) Connect(c echo.Context) WSClient {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		panic("failed to upgrade ws connection")
	}

	client := WSClient{
		c:    c,
		conn: conn,
		Name: "",
		Role: 0,
	}
	client = client.assignRole()
	Server.Clients[conn] = client

	return client
}

func (server WSServer) ListClients() []struct {
	Name string
	Role int
} {
	var names []struct {
		Name string
		Role int
	}
	for _, c := range server.Clients {
		names = append(
			names,
			struct {
				Name string
				Role int
			}{Name: c.Name, Role: c.Role},
		)
	}
	return names
}

var existsP1 bool = false
var existsP2 bool = false

func (client WSClient) assignRole() WSClient {
	switch {
	case !existsP1:
		client.Role = 1
		existsP1 = true

	case !existsP2:
		client.Role = 2
		existsP2 = true

	default:
		client.Role = 0
	}

	Server.Clients[client.conn] = client
	return client
}
func (client WSClient) unassignRole() WSClient {
	if client.Role == 1 {
		client.Role = 0
		existsP1 = false
	}

	if client.Role == 2 {
		client.Role = 0
		existsP2 = false
	}

	Server.Clients[client.conn] = client
	return client
}

func (client WSClient) WithName(name string) WSClient {
	client.Name = name
	Server.Clients[client.conn] = client
	return client
}

type Message struct {
	Type string `json:"_type"`

	UserName  string `json:"user-name"`
	GameInput string `json:"game-input"`
	ChatInput string `json:"chat-message"`

	Headers map[string]interface{} `json:"HEADERS"`
}

func (client WSClient) Read() (msg Message, err error) {
	// TODO - if i cared about security i'd sanitize the input...
	_, raw, err := client.conn.ReadMessage()
	if err != nil {
		return msg, err
	}

	err = json.Unmarshal(raw, &msg)
	if err != nil {
		panic(fmt.Errorf("couldn't parse message %s ", raw))
	}

	return msg, nil
}
func (client WSClient) Write(html []byte) error {
	return client.conn.WriteMessage(websocket.TextMessage, html)
}
func (client WSClient) Disconnect() {
	client = client.unassignRole()
	delete(Server.Clients, client.conn)
	err := client.conn.Close()
	if err != nil {
		panic("failed to close connection")
	}
}

func broadcast() {
	for {
		msg := <-Server.Broadcaster
		for _, client := range Server.Clients {
			client.Write(msg)
		}
	}
}
