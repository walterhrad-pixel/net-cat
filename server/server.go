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
}

func handleConnection(conn net.conn, hub *Hub) {
	defer conn.Close()

	conn.Write([]byte("Welcome to TCP-chat!\n"))
	conn.Write([]byte("[ENTER YOUR NAME]: "))

	scanner := bufio.NewScanner(conn)

	if !scanner.Scan() {
		return
	}
	username := strings.TrimSpace(scanner.Text())

	if username == "" {
		conn.Write([]byte("Invalid name\n"))
		return
	}
	client := NewClient(conn, username, hub)

	room := hub.GetOrCreateRoom("main")

	room.AddClient(client)

	go client.Write()
	go client.Read()

	<-client.Quit

	room.RemoveClient(client)
}
