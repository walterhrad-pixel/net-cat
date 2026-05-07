package server

import (
	"net"
	"strings"
	"testing"
	"time"
)

// --- helpers ---

// newTestClient creates a Client backed by an in-memory net.Pipe connection.
// It returns both the client-side and server-side halves of the pipe so tests
// can read/write the "remote" end.
func newTestClient(t *testing.T, username string, hub *Hub) (*Client, net.Conn) {
	t.Helper()
	serverSide, clientSide := net.Pipe()
	c := NewClient(serverSide, username, hub)
	return c, clientSide
}

// --- NewClient ---

func TestNewClient_Fields(t *testing.T) {
	hub := &Hub{}
	c, remote := newTestClient(t, "  alice  ", hub)
	defer remote.Close()

	if c.Username != "alice" {
		t.Errorf("expected username trimmed to 'alice', got %q", c.Username)
	}
	if c.Hub != hub {
		t.Error("Hub not set correctly")
	}
	if c.Send == nil {
		t.Error("Send channel is nil")
	}
	if c.Quit == nil {
		t.Error("Quit channel is nil")
	}
}

func TestNewClient_ChannelBuffers(t *testing.T) {
	hub := &Hub{}
	c, remote := newTestClient(t, "bob", hub)
	defer remote.Close()

	// Send channel should accept 256 messages without blocking
	for i := 0; i < 256; i++ {
		select {
		case c.Send <- []byte("x"):
		default:
			t.Fatalf("Send channel blocked at iteration %d (expected buffer of 256)", i)
		}
	}
}

// --- Write ---

func TestWrite_ForwardsMessages(t *testing.T) {
	hub := &Hub{}
	c, remote := newTestClient(t, "alice", hub)
	defer remote.Close()

	go c.Write()

	want := "hello from server\n"
	c.Send <- []byte(want)

	remote.SetReadDeadline(time.Now().Add(time.Second))
	buf := make([]byte, len(want))
	if _, err := remote.Read(buf); err != nil {
		t.Fatalf("Read from remote: %v", err)
	}
	if string(buf) != want {
		t.Errorf("got %q, want %q", string(buf), want)
	}
}

func TestWrite_ExitsOnQuit(t *testing.T) {
	hub := &Hub{}
	c, remote := newTestClient(t, "alice", hub)
	defer remote.Close()

	done := make(chan struct{})
	go func() {
		c.Write()
		close(done)
	}()

	// Sending on Quit must unblock Write
	c.Quit <- struct{}{}

	select {
	case <-done:
		// success
	case <-time.After(time.Second):
		t.Error("Write() did not exit after Quit signal")
	}
}

// --- SendMessage ---

func TestSendMessage_FormatsCorrectly(t *testing.T) {
	hub := &Hub{}
	c, remote := newTestClient(t, "alice", hub)
	defer remote.Close()

	go c.Write()

	msg := NewMessage(ChatMessage, "alice", "hi there", "general")
	msg.Timestamp = "2024-01-01 12:00:00"
	c.SendMessage(msg)

	remote.SetReadDeadline(time.Now().Add(time.Second))
	buf := make([]byte, 512)
	n, err := remote.Read(buf)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	got := string(buf[:n])
	if !strings.Contains(got, "alice") || !strings.Contains(got, "hi there") {
		t.Errorf("formatted message missing expected content, got: %q", got)
	}
	if !strings.HasSuffix(got, "\n") {
		t.Errorf("SendMessage should append newline, got: %q", got)
	}
}

// --- handleCommand ---

func TestHandleCommand_Name(t *testing.T) {
	hub := &Hub{}

	serverSide, clientSide := net.Pipe()
	defer clientSide.Close()

	room := &Room{Name: "general"}
	room.Broadcast = make(chan Message, 10)

	c := NewClient(serverSide, "alice", hub)
	c.Room = room

	// Intercept BroadcastMessage — use a simple channel-based room stub
	// by replacing room.Broadcast; the real BroadcastMessage sends to it.
	// (If BroadcastMessage is not channel-based in the real implementation,
	//  just verify Username changed.)

	c.handleCommand("/name bob")

	if c.Username != "bob" {
		t.Errorf("expected username 'bob' after /name, got %q", c.Username)
	}
}

func TestHandleCommand_Name_NoArg(t *testing.T) {
	hub := &Hub{}
	serverSide, clientSide := net.Pipe()
	defer clientSide.Close()

	room := &Room{Name: "general"}
	room.Broadcast = make(chan Message, 10)

	c := NewClient(serverSide, "alice", hub)
	c.Room = room

	c.handleCommand("/name")

	// No argument: username must remain unchanged
	if c.Username != "alice" {
		t.Errorf("username should be unchanged, got %q", c.Username)
	}
}

func TestHandleCommand_UnknownCommand(t *testing.T) {
	hub := &Hub{}
	serverSide, clientSide := net.Pipe()
	defer clientSide.Close()

	room := &Room{Name: "general"}
	room.Broadcast = make(chan Message, 10)

	c := NewClient(serverSide, "alice", hub)
	c.Room = room

	// Should not panic
	c.handleCommand("/unknown something")
}

func TestHandleCommand_Empty(t *testing.T) {
	hub := &Hub{}
	serverSide, clientSide := net.Pipe()
	defer clientSide.Close()

	room := &Room{Name: "general"}
	room.Broadcast = make(chan Message, 10)

	c := NewClient(serverSide, "alice", hub)
	c.Room = room

	// Should not panic on empty command
	c.handleCommand("/")
}