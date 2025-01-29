package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/api_gateway/configs"
	"github.com/gorilla/websocket"
	m "github.com/maneulf/messages_models/models"
)

type webAdapterHandler struct {
	next baseHandler
}

func (h *webAdapterHandler) handle(body m.ApiGateWayMessage, ws *websocket.Conn) {
	log.Println("Web adapter handler executed")

	if body.Provider.Name == "webAdapter" {
		h.fordwareMessage(body, ws)
	} else {
		if h.next != nil {
			h.next.handle(body, ws)
		}
	}
}

func (h *webAdapterHandler) setNext(nh baseHandler) {
	h.next = nh
}

func (h *webAdapterHandler) fordwareMessage(body m.ApiGateWayMessage, ws *websocket.Conn) {
	serverConf := configs.ConfigFromEnv("")
	url := serverConf.Service.WebAdapterMSPath
	log.Printf("Forwarding message to %s", url)

	payload := m.CsmlRequestMessage{
		Client: struct {
			UserID string "json:\"user_id\""
		}{
			UserID: body.Client.UserID,
		},
		Metadata: struct {
			Firstname string "json:\"firstname\""
			Lastname  string "json:\"lastname\""
		}{
			Firstname: body.Metadata.Firstname,
			Lastname:  body.Metadata.Lastname,
		},
		RequestID: body.RequestID,
		Payload: struct {
			Content struct {
				Text string "json:\"text\""
			} "json:\"content\""
			ContentType string "json:\"content_type\""
		}{
			Content: struct {
				Text string "json:\"text\""
			}{
				Text: body.Payload.Content.Text,
			},
			ContentType: body.Payload.ContentType,
		},
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)

	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, &buf)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	var responseBOdy map[string]interface{}
	json.NewDecoder(res.Body).Decode(&responseBOdy)

	log.Println(responseBOdy)

	ws.WriteJSON(responseBOdy)

}
