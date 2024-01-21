// All handler functions should be located in this file
package web

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

// Handler
func Root(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
