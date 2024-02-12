package paths

import (
	"RPS-backend/globals"
	"RPS-backend/responses"
	"RPS-backend/utils"
	"log"
	"net/http"
)

func ConnectToLobby(w http.ResponseWriter, r *http.Request) {
	ws, err := utils.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()
	log.Println(globals.Lobbies)
	_, err = responses.CreateResponse(responses.LobbyConnect, "Connected to lobby", "0", {globals.Lobbies})
	if err != nil {
		// handle error
	}

}
