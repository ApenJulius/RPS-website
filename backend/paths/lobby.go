package paths

import (
	"RPS-backend/globals"
	"RPS-backend/structs"
	"RPS-backend/utils"
	"log"
	"net/http"
)

func fastRemove(slice []*structs.Client, i int) []*structs.Client {
	slice[i] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}
func findClientIndex(client *structs.Client, clients []*structs.Client) int {
	for i, c := range clients {
		if c == client {
			return i
		}
	}
	return -1 // return -1 if the client is not found
}

func ConnectToLobby(w http.ResponseWriter, r *http.Request) {
	ws, err := utils.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()
	client := &structs.Client{
		Conn: ws,
		Move: "",
	}
	globals.InBrowser = append(globals.InBrowser, client)
	utils.LobbyListUpdate()

	for {
		var msg structs.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			clientIndex := findClientIndex(client, globals.InBrowser)
			if clientIndex != -1 { // make sure the client is in the slice
				globals.InBrowser = fastRemove(globals.InBrowser, clientIndex)
			}
			break
		}
	}

}
