package paths

import (
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
}
