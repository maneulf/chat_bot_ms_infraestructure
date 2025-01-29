package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/api_gateway/pkg/handlers"
	"github.com/gorilla/websocket"
	m "github.com/maneulf/messages_models/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("failed to upgrade websocket connection: %v", err)
		return
	}

	messageHandler := handlers.NewMessageHandler()

	go func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.Printf("error reading message: %v", err)
				break
			}
			log.Printf("received message: %s", msg)

			var message m.ApiGateWayMessage

			err = json.Unmarshal(msg, &message)

			if err != nil {
				log.Printf("Error: %s\n", err.Error())
			}

			log.Println(message)
			messageHandler.Handler(message, ws)

		}
	}(ws)
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	log.Println("Starting WebSocket server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
