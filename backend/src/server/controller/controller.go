package controller

import (
	"net/http"
)

type Controller struct {
	chathub *ChatHub
}

func NewController() (c *Controller) {
	c = new(Controller)
	c.chathub = NewChatHub()
	return c
}

type Binder interface {
	BindJSON(interface{}) error
}

type R map[string]interface{}

func errorMsg(msg string) (int, *R) {
	return http.StatusBadRequest, &R{
		"status":    "error",
		"error_msg": msg,
	}
}

func errorJson(err error) (int, *R) {
	return errorMsg("Error json: " + err.Error())
}

func errorUnauth() (int, *R) {
	return http.StatusUnauthorized, &R{
		"status":    "error",
		"error_msg": "Wrong token",
	}
}
