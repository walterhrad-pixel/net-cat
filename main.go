package main

import (
        "flag"
        "fmt"
        "log"
        "netcat/server"
        "os"

        "netcat/utils"
)

func main() {
        // Parse flags
        port := flag.String("p", "8989", "Port to listen on")
        room := flag.String("room", "general", "Default room name")
        flag.Parse()


        // Create and start server
        srv := server.NewServer(":" + *port)

        // Create default room
        srv.Hub.GetOrCreateRoom(*room)

        fmt.Printf("Server starting on port %s\n", *port)
        fmt.Printf("Default room: %s\n", *room)
        fmt.Println("Press Ctrl+C to stop")

                // Start server
        if err := srv.Start(); err != nil {
                log.Fatalf("Server error: %v", err)
        }

        // Display welcome banner
        fmt.Println(utils.WelcomeBanner())


        // Handle graceful shutdown
        c := make(chan struct{})
        go func() {
                <-c
                srv.Stop()
                os.Exit(0)
        }()


}
