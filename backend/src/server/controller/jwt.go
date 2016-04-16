package controller

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func createToken() (string, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	// token.Claims["foo"] = "bar"
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	// Sign and get the complete encoded token as a string
	return token.SignedString([]byte("MY_SUPER_SECRET_KEY"))
}

func validToken(token_str string) bool {
	token, err := jwt.Parse(token_str, func(token *jwt.Token) (interface{}, error) {
		return []byte("MY_SUPER_SECRET_KEY"), nil
	})

	if err != nil {
		log.Printf("Error token: %v", err)
		return false
	}
	return token.Valid
}
