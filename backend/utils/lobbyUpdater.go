package utils

import (
	"RPS-backend/globals"
	"RPS-backend/responses"
	"fmt"
)

func LobbyListUpdate() (string, error) {
	data := make(map[string]interface{})
	for id, lobby := range globals.Lobbies {
		data[id] = map[string]interface{}{
			"clients": len(lobby.Clients),
			"max":     lobby.Max,
		}
	}

	response, err := responses.CreateResponse(responses.LobbyListUpdate, "Lobby list updated", "0", data)
	if err != nil {
		return "", err
		// handle error
		fmt.Println(err)
	}
	return response, nil
}
