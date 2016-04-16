package controller

import (
	"net/http"
)

type json_login struct {
	Email string `json:"email" binding:"required"`
	Pass  string `json:"pass" binding:"required"`
}

func (c *Controller) Login(b Binder) (int, *R) {
	var json json_login
	err := b.BindJSON(&json)
	if err != nil {
		return errorJson(err)
	}

	if json.Email != "admin" || json.Pass != "admin" {
		return http.StatusUnauthorized, &R{
			"status":    "error",
			"error_msg": "Wrong password",
		}
	}

	token, err := createToken()
	if err != nil {
		return errorMsg("Wrong create token")
	}

	return http.StatusOK, &R{
		"status": "ok",
		"token":  token,
	}
}
