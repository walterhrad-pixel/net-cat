package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"netcat/utils"
)

// Server represents the TCP chat server
type Server struct {
	Addr string
	Hub  *Hub
	Quit chan struct{}
}

// NewServer creates a new server
func NewServer(addr string) *Server {
	return &Server{
		Addr: addr,
		Hub:  NewHub(),
		Quit: make(chan struct{}),
	}
}

// Start begins listening for connections
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Printf("Server listening on %s\n", s.Addr)

	for {
		select {
		case <-s.Quit:
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go s.handleConnection(conn)
		}
	}
}

// Stop shuts down the server
func (s *Server) Stop() {
	close(s.Quit)
}

// handleConnection manages a new client connection
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Send ASCII welcome banner
	conn.Write([]byte(utils.WelcomeBanner()))

	// Get username
	username, err := s.getUsername(conn)
	if err != nil {
		fmt.Printf("Error getting username: %v\n", err)
		return
	}

	// Create client with hub reference
	client := NewClient(conn, username, s.Hub)

	// Add to default room
	defaultRoom := s.Hub.GetOrCreateRoom("general")
	if !defaultRoom.AddClient(client) {
		conn.Write([]byte("Room is full. Try again later.\n"))
		return
	}

	// Send welcome message
	welcome := NewMessage(SystemMessage, "SYSTEM",
		fmt.Sprintf("Welcome %s! You joined room 'general'.", username), "general")
	welcome.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	client.SendMessage(welcome)

	// Start client goroutines
	go client.Write()
	go client.Read()

	// Wait for client to quit
	<-client.Quit

	// Cleanup
	if client.Room != nil {
		client.Room.RemoveClient(client)
	}
}

// getUsername reads username (banner already provided prompt)
func (s *Server) getUsername(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)

	for {
		username, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		username = strings.TrimSpace(username)
		if username != "" {
			return username, nil
		}

		conn.Write([]byte("Username cannot be empty. Please try again.\n"))
	}
}
