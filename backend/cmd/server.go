package main

import (
	"RPS-backend/game"
	"RPS-backend/responses"
	"RPS-backend/structs"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var Clients = make(map[*websocket.Conn]bool) // connected Clients
var broadcast = make(chan Message)           // broadcast channel

type Settings struct {
}

var groups = make(map[string]*structs.Group)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Define our message object
type Message struct {
	Move string `json:"move"`
}

func main() {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)

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

	groupID := r.URL.Query().Get("groupID")
	if groupID == "" {
		log.Println("groupID not specified") // TODO: Auto assign group
		return
	}
	if _, ok := groups[groupID]; !ok {
		groups[groupID] = &structs.Group{
			Clients: make(map[*websocket.Conn]bool),
			Max:     2,
		}
	}
	if len(groups[groupID].Clients) >= groups[groupID].Max {
		log.Printf("Group %s has reached its Max connections", groupID)
		ws.WriteJSON(map[string]string{"error": "Group has reached its Max connections"})
		return
	}

	clientAddr := ws.RemoteAddr().String()
	groups[groupID].Clients[ws] = true
	clientCount := len(groups[groupID].Clients)
	MaxClients := groups[groupID].Max
	clientInfo := fmt.Sprintf("%d/%d", clientCount, MaxClients)
	log.Printf("New client( %s ) connected to group %s : %s", clientAddr, groupID, clientInfo)
	response, err := responses.CreateResponse(responses.GameFound, "Connected to group successfully", groupID)
	if err != nil {
		// handle error
	}
	ws.WriteMessage(websocket.TextMessage, []byte(response))
	Clients[ws] = true

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
			if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
				log.Println("write:", err)
			}
		}
	}

	sendGroupUpdate() // Call after client is added
	if len(groups[groupID].Clients) == groups[groupID].Max {
		game.PlayGame(groups[groupID])
	}
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(groups[groupID].Clients, ws)
			sendGroupUpdate() // Call after client is removed
			break
		}
		log.Printf("Received message from %s: %v", clientAddr, msg) // print the client's address and the received message
		for conn := range groups[groupID].Clients {
			if err := conn.WriteJSON(msg); err != nil {
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

func gameStarting(groupID string) {
	response, err := responses.CreateResponse(responses.GameStarted, "Game has started", groupID)
	if err != nil {
		// handle error
	}
	for conn := range groups[groupID].Clients {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
			log.Println("write:", err)
		}
	}
}
