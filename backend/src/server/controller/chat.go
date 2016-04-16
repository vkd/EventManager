package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type json_get_chat struct {
	Token string `json:"token" binding:"required"`
}

func (c *Controller) GetChat(b Binder) (int, *R) {
	var json json_get_chat
	if err := b.BindJSON(&json); err != nil {
		return errorJson(err)
	}

	if !validToken(json.Token) {
		return errorUnauth()
	}

	return http.StatusOK, &R{
		"status": "ok",
	}
}

func (c *Controller) ChatWS(ws *websocket.Conn) {
	var err error

	var token string
	if err = websocket.Message.Receive(ws, &token); err != nil {
		fmt.Println("Can't receive")
		return
	}

	if !validToken(token) {
		err = websocket.Message.Send(ws, "Wrong token")
		if err != nil {
			log.Printf("Error validate token: %v", err)
		}
		return
	}

	for {
		time.Sleep(3 * time.Second)
		if err = websocket.Message.Send(ws, "hello"); err != nil {
			fmt.Println("Error send to websocket: %v", err)
			break
		}
	}
}
