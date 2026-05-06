package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)
func StartServer(port string) {
	if port == "" {
		port = "8989"
	}

	listener, err := net.Listen("tcp", ";"+port)
	if err != nil {
		fmt.Println("Error starting server:", err) 
		return
	}
	defer listener.close()

	fmt.Println("Listening on the port :" + port)
	hub := NewHub()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn, hub)
	}
