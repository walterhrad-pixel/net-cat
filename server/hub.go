package server

import (
        "sync"
)

// Hub manages all rooms and client coordination
type Hub struct {
        Rooms    map[string]*Room
        mu       sync.RWMutex
}

// NewHub creates a new hub
func NewHub() *Hub {
        return &Hub{
                Rooms: make(map[string]*Room),
        }
}

// GetOrCreateRoom returns existing room or creates new one
func (h *Hub) GetOrCreateRoom(name string) *Room {
        h.mu.Lock()
        defer h.mu.Unlock()

        if room, exists := h.Rooms[name]; exists {
                return room
        }

        room := NewRoom(name)
        h.Rooms[name] = room
        go room.Run()
        return room
}

// GetRoom returns a room by name
func (h *Hub) GetRoom(name string) (*Room, bool) {
        h.mu.RLock()
        defer h.mu.RUnlock()
        room, exists := h.Rooms[name]
        return room, exists
}

// MoveClient moves client to a different room
func (h *Hub) MoveClient(client *Client, roomName string) {
        // Remove from current room
        if client.Room != nil {
                client.Room.RemoveClient(client)
        }

        // Get or create target room
        room := h.GetOrCreateRoom(roomName)
        room.AddClient(client)
}

// RoomCount returns total number of rooms
func (h *Hub) RoomCount() int {
        h.mu.RLock()
        defer h.mu.RUnlock()
        return len(h.Rooms)
}
