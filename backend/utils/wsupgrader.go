package utils

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Configure the upgrader
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
