package server

import (
	"testing"
)

func TestNewRoom(t *testing.T) {
	r := NewRoom("testroom")
	if r.Name != "testroom" {
		t.Fatalf("expected room name 'testroom', got '%s'", r.Name)
	}

	if r.ClientCount() != 0 {
		t.Fatalf("expected 0 clients, got %d", r.ClientCount())
	}

	if r.Broadcast == nil {
		t.Fatal("expected Broadcast channel to be initialized")
	}
}
