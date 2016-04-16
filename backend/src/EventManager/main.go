package main

import (
	"demonize"
	"server"
)

func main() {
	srv := &server.Server{}

	demonize.Demonize(srv, "EventManager")
}
