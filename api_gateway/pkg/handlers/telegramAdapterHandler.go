package handlers

import (
	"log"

	"github.com/gorilla/websocket"
	m "github.com/maneulf/messages_models/models"
)

type telegramAdapterHandler struct {
	next baseHandler
}

func (h *telegramAdapterHandler) handle(body m.ApiGateWayMessage, ws *websocket.Conn) {
	log.Println("Telegram adapter handler executed")
	/*
		if provider == telegram {

		}
	*/

	h.next.handle(body, ws)
}

func (h *telegramAdapterHandler) setNext(nh baseHandler) {
	h.next = nh
}

func (h *telegramAdapterHandler) fordwareMessage(body m.ApiGateWayMessage, ws *websocket.Conn) {

}
