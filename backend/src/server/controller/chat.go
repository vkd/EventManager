package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/twinj/uuid"
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

func (c *Controller) StartChatWS(room string) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {
		c.ChatWS(room, ws)
	}
}

func (c *Controller) ChatWS(room string, ws *websocket.Conn) {
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

	email := string(uuid.NewV4().Bytes())
	broadcast, input := c.chathub.Register(RoomId(room), email)

	go func() {
		var err_rec error
		var msg string
		for {
			if err_rec = websocket.Message.Receive(ws, &msg); err_rec != nil {
				return
			}
			broadcast <- Message(msg)
		}
	}()

	for {
		select {
		case c := <-input:
			if err = websocket.Message.Send(ws, string(c)); err != nil {
				break
			}
		}
		// time.Sleep(3 * time.Second)
		// if err = websocket.Message.Send(ws, "hello"); err != nil {
		//     break
		// }
	}
}

type RoomId string
type Message string

type ChatHub struct {
	chats map[RoomId]*Chat
}

func NewChatHub() (c *ChatHub) {
	c = new(ChatHub)
	c.chats = make(map[RoomId]*Chat)
	return c
}

func (c *ChatHub) Register(room RoomId, email string) (chan<- Message, <-chan Message) {
	ch, ok := c.chats[room]
	if !ok {
		ch = NewChat()
		c.chats[room] = ch
	}

	user_chan := make(chan Message)
	ch.users[email] = user_chan
	return ch.broadcast, user_chan
}

func (c *ChatHub) Unregister(room RoomId, email string) {
	ch, ok := c.chats[room]
	if !ok {
		return
	}
	chan_user, ok := ch.users[email]
	if !ok {
		return
	}
	close(chan_user)
	delete(ch.users, email)
	if len(ch.users) == 0 {
		close(c.chats[room].broadcast)
		delete(c.chats, room)
	}
}

type Chat struct {
	broadcast chan Message
	users     map[string]chan<- Message
}

func NewChat() (c *Chat) {
	c = new(Chat)
	c.broadcast = make(chan Message)
	c.users = make(map[string]chan<- Message)
	go c.broadcasting()
	return c
}

func (c *Chat) broadcasting() {
	var m Message
	var ok bool
	for {
		select {
		case m, ok = <-c.broadcast:
			if !ok {
				return
			}
			for room := range c.users {
				c.users[room] <- m
			}
		}
	}
}
