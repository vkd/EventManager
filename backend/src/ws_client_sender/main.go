package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

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

	archive := make(map[string]int)

	go func() {
		var msg = make([]byte, 512)
		var ok bool
		var count, n int
		var line string
		for {
			n, err = ws.Read(msg)
			if err != nil {
				return
				// log.Fatal(err)
			}
			line = string(msg[:n])
			if count, ok = archive[line]; ok {
				if count <= 1 {
					delete(archive, line)
				} else {
					archive[line] = count - 1
				}
			} else {
				fmt.Printf(">>%s\n", line)
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		res, ok := archive[string(line)]
		if ok {
			archive[string(line)] = res + 1
		} else {
			archive[string(line)] = 1
		}

		// time.Sleep(3 * time.Second)
		_, err = ws.Write([]byte(`{"author":"system", "message":"` + string(line) + `"}`))
		if err != nil {
			log.Fatal(err)
		}
	}
}
