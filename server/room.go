package server

import (
	"sync"
	"time"
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
				client.Send <- FormatMessage(msg)
			}

			r.mu.Unlock()

			joinMsg := NewMessage(JoinMessage, client.Username, "", r.Name)
			joinMsg.Timestamp = time.Now().Format("2006-01-02 15:04:05")
			r.BroadcastMessage(joinMsg)

		case client := <-r.Unregister:
			r.mu.Lock()
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
			}
			r.mu.Unlock()

			leaveMsg := NewMessage(LeaveMessage, client.Username, "", r.Name)
			leaveMsg.Timestamp = time.Now().Format("2006-01-02 15:04:05")
			r.BroadcastMessage(leaveMsg)

		case msg := <-r.Broadcast:
			r.mu.Lock()

			r.History = append(r.History, msg)

			for client := range r.Clients {
				if client != msg.Sender {
					select {
					case client.Send <- FormatMessage(msg)
					default:
						close(client.Send)
						delete(r.Clients, client)
					}
				}
			}

			r.mu.Unlock()
		}
	}
}

func (r *Room) AddClient(c *Client) {
	r.Register <- c
}

func (r *Room) RemoveClient(c *Client) {
	r.Unregister <- c
}

func (r *Room) BroadcastMessage(msg Message) {
	r.Broadcast <- msg
}
