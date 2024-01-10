package handler

import (
	"fmt"
	"hackathon23/handler/game"
	"hackathon23/handler/wss"
	"hackathon23/html"
	"hackathon23/html/views"

	"github.com/labstack/echo"
)

func IndexHandler(c echo.Context) error {
	return views.Index().Render(c.Request().Context(), c.Response())
}

func WSHandler(c echo.Context) error {
	templ := html.NewBuilder(c)
	client := wss.Server.Connect(c)

	err := client.Write(templ.SignInForm())
	if err != nil {
		client.Disconnect()
		return nil
	}

	msg, err := client.Read()
	if err != nil || msg.Type != "sign-in" {
		client.Disconnect()
		return nil
	}
	// TODO - validate the UserName isn't empty on server too
	// TODO - and that usernames don't repeat
	client = client.WithName(msg.UserName)
	defer (func() {
		client.Disconnect()
		wss.Server.Broadcaster <- templ.UserLeftTheRoomMessage(client.Name, client.Role)
		wss.Server.Broadcaster <- templ.CurrentlyOnline(wss.Server.ListClients())
	})()

	err = client.Write(templ.GameScreen(game.TicTacToe.GetState(), client.Role))
	if err != nil {
		return nil
	}

	wss.Server.Broadcaster <- templ.UserJoinedTheRoomMessage(client.Name, client.Role)
	wss.Server.Broadcaster <- templ.CurrentlyOnline(wss.Server.ListClients())

	for {
		msg, err := client.Read()
		if err != nil {
			break
		}

		switch msg.Type {
		case "chat-message":
			err = client.Write(templ.EmptyChatInput())
			if err != nil {
				break
			}
			wss.Server.Broadcaster <- templ.NewChatMessage(client.Name, client.Role, msg.ChatInput)

		case "game-input":
			col := int(msg.GameInput[0] - '0')
			row := int(msg.GameInput[2] - '0')
			fmt.Println(row, col, client.Role)
			// game.TicTacToe.MakeMove(row, col, client.Role)
		}
	}

	return nil
}
