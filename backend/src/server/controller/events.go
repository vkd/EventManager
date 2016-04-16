package controller

import (
	"net/http"
)

type json_events struct {
	Token string `json:"token" binding:"required"`
}

func (c *Controller) GetEvents(b Binder) (int, *R) {
	var json json_events
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
