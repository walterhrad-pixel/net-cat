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

func TestMoveClient(t *testing.T) {
	hub := server.NewHub()

	client := &server.Client{
		Send: make(chan []byte, 1),
	}

	hub.MoveClient(client, "room1")

	if client.Room == nil {
		t.Error("client should be assigned to a room")
	}

	if client.Room.Name != "room1" {
		t.Errorf("expected room1, got %s", client.Room.Name)
	}
}

func TestMoveClientSwitchRoom(t *testing.T) {
	hub := server.NewHub()

	client := &server.Client{
		Send: make(chan []byte, 1),
	}

	hub.MoveClient(client, "room1")
	firstRoom := client.Room

	hub.MoveClient(client, "room2")

	if client.Room.Name != "room2" {
		t.Error("client did not move to new room")
	}

	if firstRoom.Name == client.Room.Name {
		t.Error("client should have switched rooms")
	}
}
