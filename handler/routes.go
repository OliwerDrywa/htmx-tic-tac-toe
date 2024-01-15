package handler

import (
	"fmt"
	"hackathon23/handler/TicTacToe"
	"hackathon23/handler/WebSocketServer"
	"hackathon23/html"
	"hackathon23/html/views"

	"github.com/labstack/echo"
)

func IndexHandler(c echo.Context) error {
	return views.Index().Render(c.Request().Context(), c.Response())
}

var ttt = TicTacToe.New()
var wss = WebSocketServer.New()

func WSHandler(c echo.Context) (err error) {
	templ := html.NewBuilder(c)
	client := wss.
		Connect(c).
		AssignRole()

	err = client.Write(templ.SignInForm())
	if err != nil {
		fmt.Println("		failed to write")
		client.
			UnassignRole().
			Disconnect()

		return nil
	}

	msg, err := client.Read()
	if err != nil || msg.Type != "sign-in" {
		client.
			UnassignRole().
			Disconnect()

		return nil
	}
	// TODO - validate the UserName isn't empty on server too
	// TODO - and that usernames don't repeat
	client.SetName(msg.UserName)
	defer (func() {
		name := client.Name
		role := client.Role
		client.
			UnassignRole().
			Disconnect()

		for _, c := range wss.Clients {
			c.Write(templ.UserLeftTheRoomMessage(name, role))
			c.Write(templ.CurrentlyOnline(ListClients(wss)))
		}

		if len(ListClients(wss)) < 2 {
			fmt.Println("less then 2 users")
		}
	})()

	err = client.Write(templ.WaitingForGame(client.Role, ttt.GetState()))
	if err != nil {
		return nil
	}

	for _, c := range wss.Clients {
		c.Write(templ.UserJoinedTheRoomMessage(client.Name, client.Role))
		c.Write(templ.CurrentlyOnline(ListClients(wss)))
	}

	// there's now 2 players
	// inform all
	if len(ListClients(wss)) == 2 {
		fmt.Println("exactly 2 users")
		for _, c := range wss.Clients {
			c.Write(templ.Game(client.Role, ttt.GetState()))
		}
	}

	// there was already 2 players
	// inform only the one joining
	if len(ListClients(wss)) > 2 {
		fmt.Println("more than 2 users")
		err = client.Write(templ.Game(client.Role, ttt.GetState()))
		if err != nil {
			return nil
		}
	}

	for {
		msg, err := client.Read()
		if err != nil {
			return nil
		}

		switch msg.Type {
		case "chat-message":
			err = client.Write(templ.EmptyChatInput())
			if err != nil {
				return nil
			}
			for _, c := range wss.Clients {
				c.Write(templ.NewChatMessage(client.Name, client.Role, msg.ChatInput))
			}

		case "game-input":
			row := int(msg.GameInput[0] - '0')
			col := int(msg.GameInput[2] - '0')

			ttt.MakeMove(row, col, client.Role)
			for _, c := range wss.Clients {
				c.Write(templ.Game(c.Role, ttt.GetState()))
			}
		}
	}
}

func ListClients(s *WebSocketServer.Server) (names []html.NameRole) {
	for _, c := range s.Clients {
		names = append(names, html.NameRole{Name: c.Name, Role: c.Role})
	}
	return names
}
