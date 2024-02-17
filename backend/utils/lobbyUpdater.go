package utils

import (
	"RPS-backend/globals"
	"RPS-backend/responses"
	"fmt"
)

func LobbyListUpdate() {
	data := make(map[string]interface{})
	for id, lobby := range globals.Lobbies {
		data[id] = map[string]interface{}{
			"clients": len(lobby.Clients),
			"max":     lobby.Max,
		}
	}

	response, err := responses.CreateResponse(responses.LobbyListUpdate, "Lobby list updated", "0", data)
	if err != nil {
		return
	}
	fmt.Println(globals.InBrowser)
	fmt.Println("Updated lobby list", response)
	for _, client := range globals.InBrowser {
		if err := client.Conn.WriteMessage(1, []byte(response)); err != nil {
			fmt.Println("LOBBY_LIST_UPDATE ERROR:", err)
			// handle error
		}
	}
}
