package paths

import (
	"RPS-backend/structs"
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
	response, err := utils.LobbyListUpdate()
	if err != nil {
		fmt.Println(err)
	}
	ws.WriteMessage(websocket.TextMessage, []byte(response))

	for {
		var msg structs.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
	}

}
