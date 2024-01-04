package html

import (
	"bytes"
	"hackathon23/handler/game"
	"hackathon23/html/components"
	"hackathon23/html/views"

	"github.com/a-h/templ"
	"github.com/labstack/echo"
)

type HTMLBuilder struct {
	c echo.Context
}

func NewBuilder(c echo.Context) HTMLBuilder {
	return HTMLBuilder{c}
}

/*
renders templ.Component to a []byte
*/
func render(c echo.Context, comp templ.Component) []byte {
	var raw bytes.Buffer

	err := comp.Render(c.Request().Context(), &raw)
	if err != nil {
		panic(err)
	}

	return raw.Bytes()
}

func (html HTMLBuilder) SignInForm() []byte {
	return render(html.c, views.SignInForm())
}
func (html HTMLBuilder) GameScreen(state game.State, role int) []byte {
	return render(html.c, views.YouJoinedRoom(state, role))
}

func (html HTMLBuilder) UserJoinedTheRoomMessage(username string, role int) []byte {
	return render(html.c, components.UserJoinedTheRoomMessage(username, role))
}
func (html HTMLBuilder) UserLeftTheRoomMessage(username string, role int) []byte {
	return render(html.c, components.UserLeftTheRoomMessage(username, role))
}
func (html HTMLBuilder) NewChatMessage(name string, variant int, text string) []byte {
	return render(html.c, components.NewChatMessage(name, variant, text))
}
func (html HTMLBuilder) CurrentlyOnline(names []struct {
	Name string
	Role int
}) []byte {
	return render(html.c, components.CurrentlyOnline(names))
}
func (html HTMLBuilder) EmptyChatInput() []byte {
	return render(html.c, components.EmptyChatInput())
}
