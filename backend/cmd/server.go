package main

import (
	"RPS-backend/globals"
	"RPS-backend/paths"
	"RPS-backend/structs"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var Clients = make(map[*websocket.Conn]bool) // connected Clients
var broadcast = make(chan structs.Message)   // broadcast channel

// Define our message object

func createLobbies() {
	globals.Lobbies = make(map[string]*structs.Group)
	globals.InBrowser = []*structs.Client{}
}

func main() {
	fs := http.FileServer(http.Dir("../public"))
	createLobbies()
	http.Handle("/", fs)
	http.HandleFunc("/game", paths.ConnectToGame)
	http.HandleFunc("/lobby", paths.ConnectToLobby)
	go handleMessages()

	host := "localhost"
	port := "8000"
	log.Printf("HTTP server started on http://%s:%s", host, port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
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
