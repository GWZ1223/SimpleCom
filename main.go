package main

import (
	"SimpleCom/server"
)

func main() {
	newServer := server.NewServer("127.0.0.1", 8888)
	newServer.Start()
}
