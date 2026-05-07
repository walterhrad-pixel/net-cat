package server

import (
	"net"
	"testing"
	"time"
)

func newTestClientForHub(t *testing.T, username string, hub *Hub) (*Client, net.Conn) {
	t.Helper()
	s, c := net.Pipe()
	return NewClient(s, username, hub), c
}

func TestNewHub(t *testing.T) {
	hub := NewHub()
	if hub.Rooms == nil || len(hub.Rooms) != 0 {
		t.Error("expected initialized empty Rooms map")
	}
}

func TestGetOrCreateRoom(t *testing.T) {
	hub := NewHub()
	room := hub.GetOrCreateRoom("general")
	if room == nil || room.Name != "general" || len(hub.Rooms) != 1 {
		t.Error("expected room 'general' to be created")
	}
	room2 := hub.GetOrCreateRoom("general")
	if room != room2 || len(hub.Rooms) != 1 {
		t.Error("expected same room instance for same room name")
	}
}

func TestGetRoom(t *testing.T) {
	hub := NewHub()
	hub.GetOrCreateRoom("general")
	room, exists := hub.GetRoom("general")
	if !exists || room == nil || room.Name != "general" {
		t.Error("expected to find room 'general'")
	}
	room, exists = hub.GetRoom("nonexistent")
	if exists || room != nil {
		t.Error("expected no room for non-existent name")
	}
}

func TestMoveClient(t *testing.T) {
	hub := NewHub()
	client, remote := newTestClientForHub(t, "alice", hub)
	defer remote.Close()
	hub.MoveClient(client, "general")
	if client.Room == nil || client.Room.Name != "general" {
		t.Error("expected client to be in room 'general'")
	}
}

func TestRoomCount(t *testing.T) {
	hub := NewHub()
	if hub.RoomCount() != 0 {
		t.Error("expected 0 rooms for empty hub")
	}
	hub.GetOrCreateRoom("general")
	hub.GetOrCreateRoom("random")
	if hub.RoomCount() != 2 {
		t.Error("expected 2 rooms")
	}
}

func TestHub_Concurrent(t *testing.T) {
	hub := NewHub()
	done := make(chan bool, 30)
	for i := 0; i < 10; i++ {
		go func() { hub.GetOrCreateRoom("general"); done <- true }()
	}
	for i := 0; i < 5; i++ {
		go func() { hub.GetOrCreateRoom(string(rune('A' + i))); done <- true }()
	}
	for i := 0; i < 10; i++ {
		go func() { _ = hub.RoomCount(); done <- true }()
	}
	for i := 0; i < 10; i++ {
		go func() { _, _ = hub.GetRoom("general"); done <- true }()
	}
	for i := 0; i < 30; i++ {
		select { case <-done: case <-time.After(time.Second): t.Error("timeout") }
	}
	if len(hub.Rooms) != 6 {
		t.Error("expected 6 rooms after concurrent access")
	}
}