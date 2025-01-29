package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	m "github.com/maneulf/messages_models/models"
	"github.com/web_adapter_ms/configs"
	"gorm.io/gorm"
)

type MessageHandler struct {
	//ch       chan m.CsmlMessage
	ch       chan Body
	isInit   bool
	CsmlPath string
	DBcon    *gorm.DB
}

type Body struct {
	c       *gin.Context
	payload m.CsmlRequestMessage
}

func (mh *MessageHandler) doCsmlRequest(message Body) (*m.CsmlResponseMessage, error) {
	serverConf := configs.ConfigFromEnv("")
	log.Printf("Forwarding message to %s", serverConf.Service.CsmlPath)
	url := serverConf.Service.CsmlPath
	csmlXApiKey := serverConf.Service.CsmlXApiKey
	method := "POST"

	requestMessage := m.CsmlRequestMessage{
		Client: struct {
			UserID string "json:\"user_id\""
		}{
			UserID: message.payload.Client.UserID,
		},
		Metadata: struct {
			Firstname string "json:\"firstname\""
			Lastname  string "json:\"lastname\""
		}{
			Firstname: message.payload.Metadata.Firstname,
			Lastname:  message.payload.Metadata.Lastname,
		},
		RequestID: message.payload.RequestID,
		Payload: struct {
			Content struct {
				Text string "json:\"text\""
			} "json:\"content\""
			ContentType string "json:\"content_type\""
		}{
			Content: struct {
				Text string "json:\"text\""
			}{
				Text: message.payload.Payload.Content.Text,
			},
			ContentType: message.payload.Payload.ContentType,
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(requestMessage)

	if err != nil {
		log.Fatalf("Error: %s", err.Error())
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, &buf)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-api-key", csmlXApiKey)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var responseMessage m.CsmlResponseMessage
	err = json.Unmarshal(body, &responseMessage)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, err
	}
	log.Println(responseMessage)

	message.c.JSON(200, responseMessage)

	return &responseMessage, nil
}

func (mh *MessageHandler) forwardMessage(message Body) {
	responseBody, err := mh.doCsmlRequest(message)

	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	go mh.DataBaseRequestMessage(message.payload, "request")
	go mh.DataBaseResponseMessage(*responseBody, "response")
}

func (mh *MessageHandler) Init() {
	if !mh.isInit {
		log.Println("Initialising messages handler ...")
		mh.isInit = true
	}
}

func (mh *MessageHandler) MessageHandler(c *gin.Context) {
	mh.Init()
	var body m.CsmlRequestMessage

	err := c.ShouldBindJSON(&body)
	if err != nil {
		log.Printf("Error on request: %s", err.Error())
	}

	b := Body{
		c:       c,
		payload: body,
	}

	mh.forwardMessage(b)
}

func (mh *MessageHandler) DataBaseRequestMessage(message m.CsmlRequestMessage, src string) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(message)

	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	serverConf := configs.ConfigFromEnv("")
	url := serverConf.Service.DataBaseMessageSaverMS
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?src=%s", url, src), &buf)

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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var responseMessage map[string]interface{}
	err = json.Unmarshal(body, &responseMessage)

	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	log.Println(responseMessage)

}

func (mh *MessageHandler) DataBaseResponseMessage(message m.CsmlResponseMessage, src string) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(message)

	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	serverConf := configs.ConfigFromEnv("")
	url := serverConf.Service.DataBaseMessageSaverMS
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?src=%s", url, src), &buf)

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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var responseMessage map[string]interface{}
	err = json.Unmarshal(body, &responseMessage)

	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	log.Println(responseMessage)

}

/*
func (mh *MessageHandler) DataBaseMessageSaverHandler(message interface{}, src string) {

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(message)

	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	serverConf := configs.ConfigFromEnv("")
	url := serverConf.Service.DataBaseMessageSaverMS
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?src=%s", url, src), &buf)

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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var responseMessage map[string]interface{}
	err = json.Unmarshal(body, &responseMessage)

	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	log.Println(responseMessage)
}
*/
