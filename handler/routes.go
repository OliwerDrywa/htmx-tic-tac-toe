package handler

import (
	"hackathon23/handler/tic_tac_toe"
	"hackathon23/handler/web_socket_server"
	"hackathon23/html"
	"hackathon23/html/views"

	"github.com/labstack/echo"
)

func IndexHandler(c echo.Context) error {
	return views.Index().Render(c.Request().Context(), c.Response())
}

var ttt = tic_tac_toe.New()
var wss = web_socket_server.New()

func WSHandler(c echo.Context) (err error) {
	templ := html.NewBuilder(c)
	client := wss.Connect(c)

	err = client.Write(templ.SignInForm())
	if err != nil {
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
						c.Write(templ.Game(c.Role, ttt.State, ttt.NowPlaying))
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
			// ttt = nil // Set ttt to nil if there are less than 2 players
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

	// there's now 2 players
	// inform all
	if len(wss.GetClients()) == 2 {
		for _, c := range wss.GetClients() {
			c.Write(templ.Game(c.Role, ttt.State, ttt.NowPlaying))
		}
	}

	// there was already 2 players
	// inform only the one joining
	if len(wss.GetClients()) > 2 {
		err = client.Write(templ.Game(client.Role, ttt.State, ttt.NowPlaying))
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

			err = ttt.MakeMove(row, col, client.Role)
			if err != nil {
				break
			}
			for _, c := range wss.GetClients() {
				c.Write(templ.Game(c.Role, ttt.State, ttt.NowPlaying))
			}
		}
	}
}

func ListClients(wss *web_socket_server.Server) (names []html.NameRole) {
	for _, c := range wss.GetClients() {
		names = append(names, html.NameRole{Name: c.Name, Role: c.Role})
	}
	return names
}
