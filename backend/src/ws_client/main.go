package main

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

func main() {
	host := "176.112.197.64"
	// host := "localhost"
	ws, err := websocket.Dial("ws://"+host+":19888/room/1/ws_chat", "", "http://"+host)
	if err != nil {
		log.Fatal(err)
	}

	message := []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjExMDcyMDV9.5XaP2409hAHi5duAtU1KnkMN1JWuj20hTlYxrY6d6lQ")
	_, err = ws.Write(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sended\n")

	var msg = make([]byte, 512)
	var n int
	for {
		n, err = ws.Read(msg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Receive: '%s'\n", msg[:n])
	}
}
