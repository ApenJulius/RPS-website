package structs

import "github.com/gorilla/websocket"

type Client struct {
	Conn *websocket.Conn
	Move string
}

type Group struct {
	Clients map[*Client]bool
	Max     int
}
type Message struct {
	Move string `json:"move"`
}
