package paths

import (
	"RPS-backend/game"
	"RPS-backend/globals"
	"RPS-backend/responses"
	"RPS-backend/structs"
	"RPS-backend/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func ConnectToGame(w http.ResponseWriter, r *http.Request) {
	ws, err := utils.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()
	client := &structs.Client{
		Conn: ws,
		Move: "",
	}
	groupID := r.URL.Query().Get("groupID")
	if groupID == "" {
		log.Println("groupID not specified") // TODO: Auto assign group
		return
	}
	if _, ok := globals.Lobbies[groupID]; !ok {

		globals.Lobbies[groupID] = &structs.Group{
			Clients: make(map[*structs.Client]bool),
			Max:     2,
		}
	}
	if len(globals.Lobbies[groupID].Clients) >= globals.Lobbies[groupID].Max {
		log.Printf("Group %s has reached its Max connections", groupID)
		ws.WriteJSON(map[string]string{"error": "Group has reached its Max connections"})
		return
	}

	clientAddr := ws.RemoteAddr().String()
	globals.Lobbies[groupID].Clients[client] = true
	clientCount := len(globals.Lobbies[groupID].Clients)
	MaxClients := globals.Lobbies[groupID].Max
	clientInfo := fmt.Sprintf("%d/%d", clientCount, MaxClients)
	log.Printf("New client( %s ) connected to group %s : %s", clientAddr, groupID, clientInfo)
	response, err := responses.CreateResponse(responses.GameFound, "Connected to group successfully", groupID)
	if err != nil {
		fmt.Println(err)
		return
		// handle error
	}
	ws.WriteMessage(websocket.TextMessage, []byte(response))

	utils.LobbyListUpdate()

	fmt.Println("Updated lobby list", response)
	utils.SendGroupUpdate(groupID) // Call after client is added
	if len(globals.Lobbies[groupID].Clients) == globals.Lobbies[groupID].Max {
		if time.Since(globals.Lobbies[groupID].LastGameStarted) > 5*time.Second {
			go game.PlayGame(globals.Lobbies[groupID], groupID)
			// Update the LastGameStarted timestamp
			globals.Lobbies[groupID].LastGameStarted = time.Now()
		}
	}
	validMoves := []string{"rock", "paper", "scissors"} // temp valid move check

	for {
		var msg structs.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(globals.Lobbies[groupID].Clients, client)

			if len(globals.Lobbies[groupID].Clients) == 0 {
				delete(globals.Lobbies, groupID)
				break
			}
			utils.SendGroupUpdate(groupID) // Call after client is removed
			break
		}
		if msg.Rematch {
			fmt.Println("rematch requested")
			if len(globals.Lobbies[groupID].Clients) == globals.Lobbies[groupID].Max {
				if time.Since(globals.Lobbies[groupID].LastGameStarted) > 5*time.Second {
					go game.PlayGame(globals.Lobbies[groupID], groupID)
					// Update the LastGameStarted timestamp
					globals.Lobbies[groupID].LastGameStarted = time.Now()
				}
			}
			continue
		}

		for _, move := range validMoves {
			if msg.Move == move {
				client.Move = msg.Move // Update the client's move
				break
			}
		}
		log.Printf("Received message from %s: %v", clientAddr, msg) // print the client's address and the received message
		for conn := range globals.Lobbies[groupID].Clients {
			if err := conn.Conn.WriteJSON(msg); err != nil {
				log.Println("write:", err)
			}
		}
	}
}
