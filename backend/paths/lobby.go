package paths

import (
	"RPS-backend/globals"
	"RPS-backend/responses"
	"RPS-backend/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func ConnectToLobby(w http.ResponseWriter, r *http.Request) {
	ws, err := utils.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()
	log.Println(globals.Lobbies)
	data := make(map[string]interface{})

	for id, lobby := range globals.Lobbies {
		data[id] = map[string]interface{}{
			"clients": len(lobby.Clients),
			"max":     lobby.Max,
		}
	}

	response, err := responses.CreateResponse(responses.LobbyConnect, "Connected to lobby", "0", data)
	if err != nil {
		// handle error
		fmt.Println(err)
	}
	ws.WriteMessage(websocket.TextMessage, []byte(response))

}
