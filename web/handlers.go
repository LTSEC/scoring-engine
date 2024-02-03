// All handler functions should be located in this file
package web

import (
	// "net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// Handler
func TableHandler(c echo.Context) error {
	services := []string{"ftp", "http", "etc."}
	teams := []string{"team a", "team b", "team c"}
	return render(c, Table(teams, services, true))
}
// This is necessary to render the templ component and display the html
func render(ctx echo.Context, cmp templ.Component) error {
	return cmp.Render(ctx.Request().Context(), ctx.Response())
}

//c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
//	return cmp.Render(c.Request().Context(), c.Response().Writer)

//Todo - arrows and scoring
// use a map, black box takes in state and score