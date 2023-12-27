package handler

import (
	"hackathon23/view"

	"github.com/labstack/echo"
)

func IndexHandler(c echo.Context) error {
	return view.Index().Render(c.Request().Context(), c.Response())
}
