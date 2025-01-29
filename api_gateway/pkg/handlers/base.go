package handlers

import (
	"github.com/gorilla/websocket"
	m "github.com/maneulf/messages_models/models"
)

type baseHandler interface {
	setNext(nh baseHandler)
	handle(body m.ApiGateWayMessage, ws *websocket.Conn)
	fordwareMessage(body m.ApiGateWayMessage, ws *websocket.Conn)
}
