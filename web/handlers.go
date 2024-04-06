// All handler functions should be located in this file
package web

import (
	// "net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/LTSEC/scoring-engine/config"
)

var frontEndConfig *config.Yaml

// function to store yaml pointer
func StoreYaml(p *config.Yaml) {
	frontEndConfig = p
}

// Handler
func TableHandler(c echo.Context) error {
	services := []string{"ftp", "http", "etc."}
	teams := []string{"team a", "team b", "team c"}
	return render(c, Table(teams, services, false))
}
// This is necessary to render the templ component and display the html
func render(ctx echo.Context, cmp templ.Component) error {
	return cmp.Render(ctx.Request().Context(), ctx.Response())
}

//Todo - scoring
// use a map, black box takes in state and score