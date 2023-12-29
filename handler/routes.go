package handler

import (
	"hackathon23/view/routes"

	"github.com/labstack/echo"
)

func IndexHandler(c echo.Context) error {
	return routes.Index().Render(c.Request().Context(), c.Response())
}
