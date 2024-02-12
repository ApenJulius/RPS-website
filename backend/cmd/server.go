package main

import (
	"RPS-backend/game"
	"RPS-backend/paths"
	"RPS-backend/responses"
	"RPS-backend/structs"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var Clients = make(map[*websocket.Conn]bool) // connected Clients
var broadcast = make(chan Message)           // broadcast channel

// Define our message object

func main() {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	http.HandleFunc("/game", paths.ConnectToGame)
	http.handleFunc("/lobby", paths.ConnectToLobby)
	go handleMessages()

	host := "localhost"
	port := "8000"
	log.Printf("HTTP server started on http://%s:%s", host, port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
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
	if _, ok := groups[groupID]; !ok {

		groups[groupID] = &structs.Group{
			Clients: make(map[*structs.Client]bool),
			Max:     2,
		}
	}
	if len(groups[groupID].Clients) >= groups[groupID].Max {
		log.Printf("Group %s has reached its Max connections", groupID)
		ws.WriteJSON(map[string]string{"error": "Group has reached its Max connections"})
		return
	}

	clientAddr := ws.RemoteAddr().String()
	groups[groupID].Clients[client] = true
	clientCount := len(groups[groupID].Clients)
	MaxClients := groups[groupID].Max
	clientInfo := fmt.Sprintf("%d/%d", clientCount, MaxClients)
	log.Printf("New client( %s ) connected to group %s : %s", clientAddr, groupID, clientInfo)
	response, err := responses.CreateResponse(responses.GameFound, "Connected to group successfully", groupID)
	if err != nil {
		// handle error
	}
	ws.WriteMessage(websocket.TextMessage, []byte(response))

	sendGroupUpdate := func() {
		data := map[string]interface{}{
			"current": len(groups[groupID].Clients),
			"max":     groups[groupID].Max,
		}
		response, err := responses.CreateResponse(responses.PlayerJoined, "Player added", groupID, data)
		if err != nil {
			// handle error
		}
		for conn := range groups[groupID].Clients {
			if err := conn.Conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
				log.Println("write:", err)
			}
		}
	}

	sendGroupUpdate() // Call after client is added
	if len(groups[groupID].Clients) == groups[groupID].Max {
		go game.PlayGame(groups[groupID])
	}
	validMoves := []string{"rock", "paper", "scissors"} // temp valid move check

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(groups[groupID].Clients, client)
			sendGroupUpdate() // Call after client is removed
			break
		}
		for _, move := range validMoves {
			if msg.Move == move {
				client.Move = msg.Move // Update the client's move
				break
			}
		}
		log.Printf("Received message from %s: %v", clientAddr, msg) // print the client's address and the received message
		for conn := range groups[groupID].Clients {
			if err := conn.Conn.WriteJSON(msg); err != nil {
				log.Println("write:", err)
			}
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(Clients, client)
			}
		}
	}
}
