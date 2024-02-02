package main

import (
	"RPS-backend/responses"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

type Group struct {
	clients map[*websocket.Conn]bool
	max     int
}

var groups = make(map[string]*Group)

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
		groups[groupID] = &Group{
			clients: make(map[*websocket.Conn]bool),
			max:     2,
		}
	}

	if len(groups[groupID].clients) >= groups[groupID].max {
		log.Printf("Group %s has reached its max connections", groupID)
		ws.WriteJSON(map[string]string{"error": "Group has reached its max connections"})
		return
	}

	clientAddr := ws.RemoteAddr().String()
	groups[groupID].clients[ws] = true
	log.Printf("New client( %s ) connected to group %s", clientAddr, groupID)
	response, err := responses.CreateResponse(responses.GameFound, "Connected to group successfully", groupID)
	if err != nil {
		// handle error
	}
	ws.WriteMessage(websocket.TextMessage, []byte(response))
	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(groups[groupID].clients, ws)
			break
		}
		log.Printf("Received message from %s: %v", clientAddr, msg) // print the client's address and the received message
		for conn := range groups[groupID].clients {
			if err := conn.WriteJSON(msg); err != nil {
				log.Println("write:", err)
			}
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
