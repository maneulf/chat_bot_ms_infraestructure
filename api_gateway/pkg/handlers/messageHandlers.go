package handlers

import (
	"github.com/gorilla/websocket"
	m "github.com/maneulf/messages_models/models"
)

type MessageHandler struct {
	base baseHandler
}

func NewMessageHandler() MessageHandler {

	webAdapter := webAdapterHandler{}
	whatsaapAdapter := whatsaapAdapterHandler{}
	telegramAdapterHandler := telegramAdapterHandler{}

	telegramAdapterHandler.setNext(&whatsaapAdapter)
	whatsaapAdapter.setNext(&webAdapter)

	mh := MessageHandler{
		base: &telegramAdapterHandler,
	}
	return mh
}

func (mh *MessageHandler) Handler(message m.ApiGateWayMessage, ws *websocket.Conn) {
	mh.base.handle(message, ws)
}
