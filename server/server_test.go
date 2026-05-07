package server_test

import (
	"testing"

	"netcat/server"
)

func TestGetOrCreateRoom(t *testing.T) {
	hub := server.NewHub()

	room1 := hub.GetOrCreateRoom("main")
	room2 := hub.GetOrCreateRoom("main")

	if room1 != room2 {
		t.Error("expected same room instance, got different ones")
	}
}

func TestRoomCount(t *testing.T) {
	hub := server.NewHub()

	hub.GetOrCreateRoom("roomA")
	hub.GetOrCreateRoom("roomB")
	hub.GetOrCreateRoom("roomC")

	if hub.RoomCount() != 3 {
		t.Errorf("expected 3 rooms, got %d", hub.RoomCount())
	}
}
