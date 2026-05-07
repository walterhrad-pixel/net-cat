package server

import (
	"fmt"
	"testing"
	"time"
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

func TestAddRemoveAndBroadcast(t *testing.T) {
	r := NewRoom("r1")
	go r.Run()
	defer close(r.Broadcast)

	// create clients
	sender := &Client{Username: "alice", Send: make(chan []byte, 10)}
	receiver := &Client{Username: "bob", Send: make(chan []byte, 10)}

	if !r.AddClient(sender) {
		t.Fatal("expected to add sender")
	}
	if !r.AddClient(receiver) {
		t.Fatal("expected to add receiver")
	}

	drain := func(c *Client) {
		for {
			select {
			case <-c.Send:
			default:
				return
			}
		}
	}
	drain(sender)
	drain(receiver)

	// send a chat message from sender
	msg := NewMessage(ChatMessage, sender.Username, "hello world", r.Name)
	msg.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	msg.Sender = sender
	r.BroadcastMessage(msg)

	// receiver should get the message
	select {
	case b := <-receiver.Send:
		s := string(b)
		if s == "" || !contains(s, "hello world") {
			t.Fatalf("unexpected message delivered to receiver: %q", s)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for receiver to get broadcast message")
	}

	// sender should NOT receive its own message
	select {
	case b := <-sender.Send:
		t.Fatalf("sender should not receive own message, got %q", string(b))
	default:
	}

	// remove receiver and expect a leave notice delivered to others
	r.RemoveClient(receiver)

	select {
	case b := <-sender.Send:
		if !contains(string(b), "has left the room") {
			t.Fatalf("expected leave notice, got %q", string(b))
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for leave notice")
	}
}

func TestRoomCapacity(t *testing.T) {
	r := NewRoom("cap")

	// fill room up to capacity
	for i := 0; i < maxClientsPerRoom; i++ {
		c := &Client{Username: fmt.Sprintf("u%d", i), Send: make(chan []byte, 1)}
		if !r.AddClient(c) {
			t.Fatalf("expected add to succeed for index %d", i)
		}
	}

	// next add should fail
	extra := &Client{Username: "extra", Send: make(chan []byte, 1)}
	if r.AddClient(extra) {
		t.Fatal("expected add to fail when room is at capacity")
	}
}
