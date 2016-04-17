package controller

import (
	"net/http"
)

type json_login struct {
	Email string `json:"email" binding:"required"`
	Pass  string `json:"pass" binding:"required"`
}

var (
	db_users = []struct {
		Login string
		Pass  string
	}{
		{"admin", "admin"},
		{"user", "user"},
		{"guest", "guest"},
	}
)

func (c *Controller) Login(b Binder) (int, *R) {
	var json json_login
	err := b.BindJSON(&json)
	if err != nil {
		return errorJson(err)
	}

	is_exists := false
	for _, u := range db_users {
		if u.Login == json.Email && u.Pass == json.Pass {
			is_exists = true
			break
		}
	}

	if !is_exists {
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
