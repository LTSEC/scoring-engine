// Main body of the program lives here
package main

import (
	"github.com/LTSEC/scoring-engine/cli"
	"github.com/LTSEC/scoring-engine/web"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", web.Root)

	go cli.Cli()

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
