package server

import (
	"bufio"
 	"net"
	"strings"
	"time"
)

// Client represents a connected user
type Client struct {
	Conn     net.Conn
	Username string
	Room     *Room
	Hub      *Hub
	Send     chan []byte
	Quit     chan struct{}
}

//NewClient creates a new client
func NewClient(conn net.Conn, username string, hub *Hub) *Client {
	return &Client{
		Conn:     conn,
		Username: strings.TrimSpace(username),
		Hub:     hub,
		Send:     make(chan []byte, 256),
		Quit:     make(chan struct{}),
	}
}

//Read reads messages from client connectionSSSS
func (c *Client) Read() {
	defer func() {
		c.Quit <- struct{}{}
		c.Conn.Close()
	}()

	scanner := bufio.NewScanner(c.Conn)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		//Check for commands
		if strings.HasPrefix(text, "/") {
			c.handleCommand(text)
			continue
		}

		//Regular chat message 
		msg := NewMessage(ChatMessage, c.Username, text, c.Room.Name)
		msg.Timestamp = time.Now().Format("2006-01-02 15:04:05")
		msg.Sender = c 
		c.Room.BroadcastMessage(msg) 
	}
}

//Write sends messages to client
func (c *Client) Write() {
	defer func() {
		c.Quit <- struct{}{}
		c.Conn.Close()
	}()

	for {
		select {
		case msg := <-c.Send:
			_, err := c.Conn.Write(msg)
			if err != nil {
				return 
			}
		case <-c.Quit:
			return 
		}
	}
}

//handleCommand processes client commands
func (c *Client) handleCommand(cmd string) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return 
	}

	switch parts[0] {
	case "/name":
		if len(parts) > 1 && parts[1] != "" {
			oldName := c.Username
			c.Username = strings.TrimSpace(parts[1])
			sysMsg := NewMessage(SystemMessage, "SYSTEM", oldName+" change name to "+c.Username, c.Room.Name)
			sysMsg.Timestamp = time.Now().Format("2006-01-02 15:04:05")
			c.Room.BroadcastMessage(sysMsg)
		}
	case "/join":
		if len(parts) > 1 && parts[1] != "" {
			c.Hub.MoveClient(c, strings.TrimSpace(parts[1]))
		}
	}
}

//SendMessage sends a formatted message to client
func (c *Client) SendMessage(msg Message) {
	formatted := msg.String() + "\n"
	c.Send <- []byte(formatted)
}
