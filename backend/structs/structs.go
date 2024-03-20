package structs

import (
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn  *websocket.Conn
	Move  string
	Score int
}

type Group struct {
	Clients         map[*Client]bool
	Max             int
	LastGameStarted time.Time
	GameInProgress  bool
}
type Message struct {
	Move    string `json:"move"`
	Rematch bool   `json:"rematch,omitempty"`
}
type Settings struct {
}
