package utils

import (
	"RPS-backend/globals"
	"RPS-backend/responses"
	"log"

	"github.com/gorilla/websocket"
)

func SendGroupUpdate(groupID string) {
	data := map[string]interface{}{
		"current": len(globals.Lobbies[groupID].Clients),
		"max":     globals.Lobbies[groupID].Max,
	}
	response, err := responses.CreateResponse(responses.PlayerJoined, "Player added", groupID, data)
	if err != nil {
		// handle error
	}
	for conn := range globals.Lobbies[groupID].Clients {
		if err := conn.Conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
			log.Println("write:", err)
		}
	}
}
