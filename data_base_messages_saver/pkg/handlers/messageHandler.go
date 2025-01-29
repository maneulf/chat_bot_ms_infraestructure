package handlers

import (
	"encoding/json"
	"log"

	db "github.com/data_base_messages_saver/database"
	"github.com/gin-gonic/gin"
	m "github.com/maneulf/messages_models/models"
	"gorm.io/gorm"
)

type MessageHandler struct {
	db *gorm.DB
}

func NewMessageHandler() *MessageHandler {
	db, err := db.ConnectToMariaDB()
	if err != nil {
		panic(err.Error())
	}
	return &MessageHandler{
		db: db,
	}

}

func (mh *MessageHandler) CsmlMessageHandler(c *gin.Context) {

	src := c.Query("src")

	var message string
	var userID string
	var requestID string

	if src == "request" {
		log.Println("Saving request message")
		var requestMessage m.CsmlRequestMessage
		err := c.ShouldBindJSON(&requestMessage)

		if err != nil {
			log.Printf("Error: %s", err.Error())
		}
		userID = requestMessage.Client.UserID
		requestID = requestMessage.RequestID
		bytes, err := json.Marshal(&requestMessage)

		if err != nil {
			log.Printf("Error: %s", err.Error())
		}
		message = string(bytes)

	} else if src == "response" {
		log.Println("Saving response message")
		var responseMessage m.CsmlResponseMessage
		err := c.ShouldBindJSON(&responseMessage)
		userID = responseMessage.Client.UserID
		requestID = responseMessage.RequestID

		if err != nil {
			log.Printf("Error: %s", err.Error())
		}
		bytes, err := json.Marshal(&responseMessage)

		if err != nil {
			log.Printf("Error: %s", err.Error())
		}
		message = string(bytes)

	}

	c.JSON(200, gin.H{"message": "message saved succefull"})

	go mh.CsmlMessageSaver(message, userID, requestID, src)

}

func (mh *MessageHandler) CsmlMessageSaver(message string, user_id string, request_id string, src string) {
	dbmessage := m.CsmlDataBaseMessageModelDB{
		Message:   message,
		UserID:    user_id,
		RequestID: request_id,
		Source:    src,
	}

	mh.db.Create(&dbmessage)

}
