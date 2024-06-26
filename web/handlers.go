package web

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/LTSEC/scoring-engine/score_holder"
	"github.com/a-h/templ"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
	clients  = make(map[*websocket.Conn]bool)
	mu       sync.Mutex
)

// Handler for serving the table page
func TableHandler(c echo.Context) error {
	allTeams := score_holder.GetMap()
	if len(allTeams) > 0 {
		teams := make([]string, len(allTeams))
		for i := range allTeams {
			teams[i] = "Team " + strconv.Itoa(i)
		}
		services := []string{"ftp", "http", "ssh"}
		return render(c, Table(teams, services))
	} else {
		return errors.New("not enough services/teams provided")
	}
}

func render(ctx echo.Context, cmp templ.Component) error {
	return cmp.Render(ctx.Request().Context(), ctx.Response())
}

// WebSocket handler
func WebSocketHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	clients[ws] = true

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			delete(clients, ws)
			break
		}
	}

	return nil
}

// Function to broadcast updates to all connected WebSocket clients
func BroadcastUpdates() {
	for {
		time.Sleep(1 * time.Second)

		mu.Lock()
		// getting data from the score holder
		data := score_holder.GetMap()
		mu.Unlock()

		message, _ := json.Marshal(data)

		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}
