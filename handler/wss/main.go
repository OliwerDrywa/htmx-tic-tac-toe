package wss

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var Server = WSServer{make(map[*websocket.Conn]*WSClient)}

type WSServer struct {
	Clients map[*websocket.Conn]*WSClient
}
type WSClient struct {
	c      echo.Context
	conn   *websocket.Conn
	server *WSServer
	Name   string
	Role   int
}

var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func (server *WSServer) Connect(c echo.Context) *WSClient {
	fmt.Println("connecting...")
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		panic("failed to upgrade ws connection")
	}

	client := WSClient{
		c:      c,
		conn:   conn,
		server: server,
		Name:   "anonymous",
		Role:   0,
	}
	client.assignRole()
	Server.Clients[conn] = &client

	return &client
}

func (server *WSServer) ListClients() []struct {
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

func (client *WSClient) assignRole() *WSClient {
	if !existsP1 {
		client.Role = 1
		existsP1 = true
	} else if !existsP2 {
		client.Role = 2
		existsP2 = true
	}

	return client
}
func (client *WSClient) unassignRole() *WSClient {
	if client.Role == 1 {
		client.Role = 0
		existsP1 = false
	} else if client.Role == 2 {
		client.Role = 0
		existsP2 = false
	}

	return client
}
func (client *WSClient) SetName(name string) *WSClient {
	client.Name = name
	return client
}

type Message struct {
	Type string `json:"_type"`

	UserName  string `json:"user-name"`
	GameInput string `json:"game-input"`
	ChatInput string `json:"chat-message"`

	Headers map[string]interface{} `json:"HEADERS"`
}

func (client *WSClient) Read() (msg Message, err error) {
	fmt.Printf("		%s.Read()\n\n", client.Name)
	// TODO - if i cared about security i'd sanitize the input...
	_, raw, err := client.conn.ReadMessage()
	if err != nil {
		fmt.Println("		failed to read")
		return msg, err
	}

	err = json.Unmarshal(raw, &msg)
	if err != nil {
		panic(fmt.Errorf("couldn't parse message %s ", raw))
	}

	return msg, nil
}
func (client *WSClient) Write(html []byte) (err error) {
	fmt.Printf("		%s.Write()\n", client.Name)
	return client.conn.WriteMessage(websocket.TextMessage, html)
}
func (client *WSClient) Disconnect() {
	fmt.Println("	disconnecting...")
	client.unassignRole()
	delete(Server.Clients, client.conn)
	err := client.conn.Close()
	if err != nil {
		panic("failed to close connection")
	}
}
