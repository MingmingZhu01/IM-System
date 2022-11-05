package main

import (
	"example.com/IM-System/server"
)

func main() {
	/*
		server := NewServer("127.0.0.1", 8888)
		server.Start()*/
	server1 := server.NewServer("127.0.0.1", 8888)
	server1.Start()
}
