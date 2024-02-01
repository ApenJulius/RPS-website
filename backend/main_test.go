package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestHandleConnections(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleConnections))
	defer server.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.1
	u := url.URL{Scheme: "ws", Host: server.URL[7:], Path: "/ws"}

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer ws.Close()

	// Send message
	err = ws.WriteJSON(Message{Move: "test move"})
	if err != nil {
		t.Fatalf("unable to write JSON to ws: %v", err)
	}

	// Give the server a chance to receive the message
	time.Sleep(2 * time.Second)

	// Check if the message was received
	if len(broadcast) == 0 {
		t.Fatalf("message not received by server")
	}

	// Read message
	var msg Message
	err = ws.ReadJSON(&msg)
	if err != nil {
		t.Fatalf("unable to read JSON from ws: %v", err)
	}

	// Check if the message was broadcasted correctly
	if msg.Move != "test move" {
		t.Fatalf("message not broadcasted correctly: got %v, want %v", msg.Move, "test move")
	}
}
