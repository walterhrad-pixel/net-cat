package main

import (
	"fmt"
	"log"
	"netcat/server"
	"os"
)

func main() {
	// 1. Handle Arguments manually to match [USAGE] requirements
	args := os.Args[1:]
	port := "8989" // Default port

	if len(args) == 1 {
		port = args[0]
	} else if len(args) > 1 {
		// This handles the "./TCPChat 2525 localhost" case
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	// 2. Create and start server with the provided port
	srv := server.NewServer(":" + port)

	// 3. Create default room
	srv.Hub.GetOrCreateRoom("general")

	// 4. Print startup info BEFORE starting (Start() is blocking)
	fmt.Printf("Listening on the port :%s\n", port)
	// Optional project-specific prints:
	// fmt.Println("Default room: general")
	// fmt.Println("Press Ctrl+C to stop")

	// 5. Handle graceful shutdown in a goroutine
	go func() {
		// This will catch a manual stop if you implement signal handling, 
		// but for now, it allows the server to run.
		<-srv.Quit
		os.Exit(0)
	}()

	// 6. Start server (This blocks until the server stops)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
