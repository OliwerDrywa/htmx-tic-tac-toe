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

var ttt *TicTacToe.Game // Declare a pointer to TicTacToe.Game
var wss = WebSocketServer.New()

func WSHandler(c echo.Context) (err error) {
	templ := html.NewBuilder(c)
	client := wss.Connect(c)

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

	// TODO - validate the UserName isn't empty on server too and that usernames don't repeat
	client.
		SetName(msg.UserName).
		AssignRole()

	defer (func() {
		client.Disconnect()

		// Reassign role if user with Role = 1 or 2 disconnects
		if client.Role == 1 || client.Role == 2 {
			for _, c := range wss.GetClients() {
				if c.Role == 0 {
					c.Role = client.Role
					if ttt != nil {
						c.Write(templ.Game(c.Role, ttt.State))
					}
					break
				}
			}
		}

		for _, c := range wss.GetClients() {
			c.Write(templ.UserLeftTheRoomMessage(client.Name, client.Role))
			c.Write(templ.CurrentlyOnline(ListClients(wss)))
		}

		client.UnassignRole()

		if len(wss.GetClients()) < 2 {
			fmt.Println("less then 2 users")
			ttt = nil // Set ttt to nil if there are less than 2 players
			for _, c := range wss.GetClients() {
				c.Write(templ.NoGame())
			}
		}
	})()

	err = client.Write(templ.WaitingForGame())
	if err != nil {
		return nil
	}

	for _, c := range wss.GetClients() {
		c.Write(templ.UserJoinedTheRoomMessage(client.Name, client.Role))
		c.Write(templ.CurrentlyOnline(ListClients(wss)))
	}

	if len(wss.GetClients()) >= 2 && ttt == nil {
		ttt = TicTacToe.New() // Create a new game of TicTacToe
	}

	// there's now 2 players
	// inform all
	if len(wss.GetClients()) == 2 {
		fmt.Println("exactly 2 users")
		for _, c := range wss.GetClients() {
			c.Write(templ.Game(client.Role, ttt.State))
		}
	}

	// there was already 2 players
	// inform only the one joining
	if len(wss.GetClients()) > 2 {
		fmt.Println("more than 2 users")
		err = client.Write(templ.Game(client.Role, ttt.State))
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
			for _, c := range wss.GetClients() {
				c.Write(templ.NewChatMessage(client.Name, client.Role, msg.ChatInput))
			}

		case "game-input":
			row := int(msg.GameInput[0] - '0')
			col := int(msg.GameInput[2] - '0')

			ttt.MakeMove(row, col, client.Role)
			for _, c := range wss.GetClients() {
				c.Write(templ.Game(c.Role, ttt.State))
			}
		}
	}
}

func ListClients(wss *WebSocketServer.Server) (names []html.NameRole) {
	for _, c := range wss.GetClients() {
		names = append(names, html.NameRole{Name: c.Name, Role: c.Role})
	}
	return names
}
