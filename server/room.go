package server

import (
	"sync"
	"time"
)

const maxClientsPerRoom = 10

// Room represents a chat room
type Room struct {
	Name      string
	Clients   map[*Client]bool
	Broadcast chan Message
	History   []Message
	HistoryMu sync.RWMutex
	mu        sync.RWMutex
}

// NewRoom creates a new room
func NewRoom(name string) *Room {
	return &Room{
		Name:      name,
		Clients:   make(map[*Client]bool),
		Broadcast: make(chan Message, 100),
		History:   make([]Message, 0, 100),
	}
}

// AddClient adds a client to the room
func (r *Room) AddClient(client *Client) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.Clients) >= maxClientsPerRoom {
		return false
	}

	r.Clients[client] = true
	client.Room = r

	// Send join message to others (including sender for confirmation)
	joinMsg := NewMessage(JoinMessage, client.Username, "", r.Name)
	joinMsg.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	r.Broadcast <- joinMsg

	// Send history to new client
	r.sendHistory(client)

	return true
}

// RemoveClient removes a client from the room
func (r *Room) RemoveClient(client *Client) {
	r.mu.Lock()
	delete(r.Clients, client)
	r.mu.Unlock()

	// Send leave message
	leaveMsg := NewMessage(LeaveMessage, client.Username, "", r.Name)
	leaveMsg.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	r.Broadcast <- leaveMsg
}

// BroadcastMessage sends a message to all clients in the room
func (r *Room) BroadcastMessage(msg Message) {
	r.HistoryMu.Lock()
	r.History = append(r.History, msg)
	// Keep last 100 messages
	if len(r.History) > 100 {
		r.History = r.History[len(r.History)-100:]
	}
	r.HistoryMu.Unlock()

	r.Broadcast <- msg
}

// Run processes room messages
func (r *Room) Run() {
	for msg := range r.Broadcast {
		r.mu.RLock()
		for client := range r.Clients {
			// Skip sender for chat messages to avoid echo
			if msg.Sender != nil && client == msg.Sender {
				continue
			}
			select {
			case client.Send <- []byte(msg.String() + "\n"):
			default:
				// Client buffer full, skip
			}
		}
		r.mu.RUnlock()
	}
}

// sendHistory sends chat history to a client (without extra labels)
func (r *Room) sendHistory(client *Client) {
	r.HistoryMu.RLock()
	defer r.HistoryMu.RUnlock()

	for _, msg := range r.History {
		client.SendMessage(msg)
	}
}

// ClientCount returns number of clients in room
func (r *Room) ClientCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.Clients)
}
