package server

import (
	"sync"
)

type Room struct {
	Name       string
	Clients    map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	History    []Message
	mu         sync.Mutex
}

func NewRoom(name string) *Room {
	return &Room{
		Name:       name,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		History:    []Message{},
	}
}
func (r *Room) Run() {
	for {
		select {
		case client := <-r.Register:
			r.mu.Lock()
			r.Clients[client] = true
			client.Room = r

			for _, msg := range r.History {
				client.Send <- []byte(msg.String() + "\n")
			}

			r.mu.Unlock

			joinMsg := NewMessage(JoinMessage, client.Username, "", r.Name)
			joinMsg.Timestamp = time.Now().Format("2006-01-02 15:04:05")
			r.BroadcastMessage(joinMsg)
