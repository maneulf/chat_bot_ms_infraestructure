package handlers

import (
	"log"

	"github.com/gorilla/websocket"
	m "github.com/maneulf/messages_models/models"
)

type whatsaapAdapterHandler struct {
	next baseHandler
}

func (h *whatsaapAdapterHandler) handle(body m.ApiGateWayMessage, ws *websocket.Conn) {
	log.Println("Whatsapp adapter handler executed")

	/*
		if provider == whatsapp {

		}
	*/
	if h.next != nil {
		h.next.handle(body, ws)
	}

}

func (h *whatsaapAdapterHandler) setNext(nh baseHandler) {
	h.next = nh
}

func (h *whatsaapAdapterHandler) fordwareMessage(body m.ApiGateWayMessage, ws *websocket.Conn) {

}
