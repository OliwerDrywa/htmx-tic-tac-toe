package web_socket_server

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

type Server struct {
	clients  map[*websocket.Conn]*client
	existsP1 bool
	existsP2 bool
}
type client struct {
	ctx    echo.Context
	conn   *websocket.Conn
	server *Server
	Name   string
	Role   int
}

func New() *Server {
	return &Server{
		make(map[*websocket.Conn]*client),
		false,
		false,
	}
}

var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func (s *Server) Connect(ctx echo.Context) *client {
	conn, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		panic("failed to upgrade ws connection")
	}

	client := client{ctx: ctx, conn: conn, server: s}
	s.clients[conn] = &client

	fmt.Println(client.Name)
	fmt.Println(client.Role)
	return &client
}

type Message struct {
	Type string `json:"_type"`

	UserName  string `json:"user-name"`
	GameInput string `json:"game-input"`
	ChatInput string `json:"chat-message"`

	Headers map[string]interface{} `json:"HEADERS"`
}

func (c *client) Read() (msg Message, err error) {
	// TODO - if i cared about security i'd sanitize the input...
	_, raw, err := c.conn.ReadMessage()
	if err != nil {
		return msg, err
	}

	err = json.Unmarshal(raw, &msg)
	if err != nil {
		panic(fmt.Errorf("couldn't parse message %s ", raw))
	}

	return msg, nil
}
func (c *client) Write(html []byte) (err error) {
	return c.conn.WriteMessage(websocket.TextMessage, html)
}
func (c *client) Disconnect() {
	delete(c.server.clients, c.conn)
	err := c.conn.Close()
	if err != nil {
		panic("failed to close connection")
	}
}
