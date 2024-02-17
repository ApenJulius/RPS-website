package game

import (
	"RPS-backend/responses"
	"RPS-backend/structs"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func PlayGame(group *structs.Group, groupID string) {
	gameCountdown(5, group, groupID)
	compareMoves(group, groupID)
}

func gameCountdown(seconds int, group *structs.Group, groupID string) {
	for i := seconds; i > 0; i-- {
		response, err := responses.CreateResponse(responses.GameCountdown, fmt.Sprintf("Choices locked in %d", i), groupID)
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
func compareMoves(group *structs.Group, groupID string) { // 2player max

	if len(group.Clients) != 2 {
		fmt.Println("Not enough players to compare moves")
		return
	}
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
		clients[0].Conn.WriteMessage(websocket.TextMessage, []byte("It's a draw!"))
		clients[1].Conn.WriteMessage(websocket.TextMessage, []byte("It's a draw!"))
	} else if (moves[0] == "rock" && moves[1] == "scissors") ||
		(moves[0] == "scissors" && moves[1] == "paper") ||
		(moves[0] == "paper" && moves[1] == "rock") {
		fmt.Println("Client", clients[0], "wins!")
		clients[0].Conn.WriteMessage(websocket.TextMessage, []byte("You win!"))
		clients[1].Conn.WriteMessage(websocket.TextMessage, []byte("You lose!"))
	} else {
		fmt.Println("Client", clients[1], "wins!")
		clients[0].Conn.WriteMessage(websocket.TextMessage, []byte("You lose!"))
		clients[1].Conn.WriteMessage(websocket.TextMessage, []byte("You win!"))
	}
}
