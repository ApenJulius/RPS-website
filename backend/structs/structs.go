package structs

import "github.com/gorilla/websocket"

type Group struct {
	Clients map[*websocket.Conn]bool
	Max     int
}
