package game

import (
	"RPS-backend/responses"
	"RPS-backend/structs"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func PlayGame(group *structs.Group) {
	gameCountdown(5, group)
	compareMoves(group)
}

func gameCountdown(seconds int, group *structs.Group) {
	for i := seconds; i > 0; i-- {
		response, err := responses.CreateResponse(responses.GameCountdown, fmt.Sprintf("Game starting in %d seconds", i), "groupID")
		if err != nil {
			// handle error
		}
		fmt.Println(response)
		for conn := range group.Clients {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
				log.Println("write:", err)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
