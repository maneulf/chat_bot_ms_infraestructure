package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	m "github.com/maneulf/messages_models/models"
)

type waitingForButtons struct {
	message m.CsmlResponseMessage
}

func main() {
	// Define the WebSocket server URL
	serverURL := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("connecting to %s\n", serverURL.String())

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(serverURL.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	ch := make(chan waitingForButtons)

	go func(ch chan waitingForButtons) {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			handleCsmlMessage(message, ch)
		}
	}(ch)

	user_id := uuid.NewString()
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">>>>> ")

		var text string
		text, err = reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		select {
		case wfb := <-ch:
			if err != nil {
				log.Printf("Error: %s\n", err.Error())
			}
			buttons := wfb.message.Messages[0].Payload.Content.Buttons
			n, err := strconv.Atoi(text)
			if err == nil {
				text = buttons[n-1].Content.Accepts[0]
			}
		default:
			if err != nil {
				log.Printf("Error: %s\n", err.Error())
			}
		}

		message := m.ApiGateWayMessage{
			Client: m.Client{
				UserID: user_id,
			},
			Metadata: m.Metadata{
				Firstname: "Manuel",
				Lastname:  "Echeverry",
			},
			RequestID: uuid.New().String(),
			Payload: m.Payload{
				Content: m.Content{
					Text: text,
				},
				ContentType: "text",
			},
			Provider: m.Provider{
				Name: "webAdapter",
			},
		}

		m, _ := json.Marshal(message)
		//fmt.Println(string(m))

		err = conn.WriteMessage(websocket.TextMessage, m)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}

func handleCsmlMessage(message []byte, ch chan waitingForButtons) {
	var msg m.CsmlResponseMessage

	json.Unmarshal(message, &msg)

	if len(msg.Messages) == 0 {
		return
	}

	if len(msg.Messages) > 0 {
		for _, m := range msg.Messages {
			if m.Payload.ContentType == "text" {
				fmt.Println(m.Payload.Content.Text)

			} else if m.Payload.ContentType == "question" {
				fmt.Println(m.Payload.Content.Title)

				if len(m.Payload.Content.Buttons) > 0 {
					i := 1
					for _, b := range m.Payload.Content.Buttons {
						fmt.Printf("%d) %s\n", i, b.Content.Title)
						i++
					}
					go func() {
						ch <- waitingForButtons{message: msg}
					}()
				}
			}
		}
		fmt.Print(">>>>> ")
	}
}
