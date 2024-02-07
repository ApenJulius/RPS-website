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
		response, err := responses.CreateResponse(responses.GameCountdown, fmt.Sprintf("Choices locked in %d", i), "groupID")
		if err != nil {
			// handle error
		}
		fmt.Println(response)
		for conn := range group.Clients {
			if err := conn.Conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
				log.Println("write:", err)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
func compareMoves(group *structs.Group) { // 2player max
	fmt.Println("Comparing moves")
	for client := range group.Clients {
		fmt.Println("Client Moves:", client.Move)
	}
	var moves []string
	var clients []*structs.Client

	for client := range group.Clients {
		moves = append(moves, client.Move)
		clients = append(clients, client)
	}

	if moves[0] == moves[1] {
		fmt.Println("It's a draw!")
	} else if (moves[0] == "rock" && moves[1] == "scissors") ||
		(moves[0] == "scissors" && moves[1] == "paper") ||
		(moves[0] == "paper" && moves[1] == "rock") {
		fmt.Println("Client", clients[0], "wins!")
	} else {
		fmt.Println("Client", clients[1], "wins!")
	}
}
