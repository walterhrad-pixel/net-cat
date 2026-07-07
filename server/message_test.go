package server

import (
	"encoding/json"
	"testing"
)

func TestNewMessage(t *testing.T) {
	msg := NewMessage(ChatMessage, "alice", "hello world", "general")

	if msg.Type != ChatMessage {
		t.Errorf("expected ChatMessage, got %v", msg.Type)
	}
	if msg.Username != "alice" {
		t.Errorf("expected username 'alice', got %q", msg.Username)
	}
	if msg.Content != "hello world" {
		t.Errorf("expected content 'hello world', got %q", msg.Content)
	}
	if msg.Room != "general" {
		t.Errorf("expected room 'general', got %q", msg.Room)
	}
}

func TestMessageString_Chat(t *testing.T) {
	msg := NewMessage(ChatMessage, "alice", "hello", "general")
	msg.Timestamp = "2024-01-01 12:00:00"

	got := msg.String()
	want := "[2024-01-01 12:00:00][alice]: hello"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestMessageString_System(t *testing.T) {
	msg := NewMessage(SystemMessage, "SYSTEM", "server restarting", "general")
	msg.Timestamp = "2024-01-01 12:00:00"

	got := msg.String()
	want := "[2024-01-01 12:00:00][SYSTEM]: server restarting"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestMessageString_Join(t *testing.T) {
	msg := NewMessage(JoinMessage, "bob", "", "general")
	msg.Timestamp = "2024-01-01 12:00:00"

	got := msg.String()
	want := "[2024-01-01 12:00:00][SYSTEM]: bob has joined the room"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestMessageString_Leave(t *testing.T) {
	msg := NewMessage(LeaveMessage, "bob", "", "general")
	msg.Timestamp = "2024-01-01 12:00:00"

	got := msg.String()
	want := "[2024-01-01 12:00:00][SYSTEM]: bob has left the room"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestMessageString_DefaultFallback(t *testing.T) {
	msg := Message{
		Type:      CommandMessage,
		Username:  "alice",
		Content:   "cmd output",
		Timestamp: "2024-01-01 12:00:00",
	}

	got := msg.String()
	want := "[2024-01-01 12:00:00][alice]: cmd output"
	if got != want {
		t.Errorf("String() default fallback = %q, want %q", got, want)
	}
}

func TestMessageToJSON(t *testing.T) {
	msg := NewMessage(ChatMessage, "alice", "hello", "general")
	msg.Timestamp = "2024-01-01 12:00:00"

	data, err := msg.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error: %v", err)
	}

	var decoded Message
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if decoded.Username != msg.Username {
		t.Errorf("username mismatch: got %q, want %q", decoded.Username, msg.Username)
	}
	if decoded.Content != msg.Content {
		t.Errorf("content mismatch: got %q, want %q", decoded.Content, msg.Content)
	}
	if decoded.Room != msg.Room {
		t.Errorf("room mismatch: got %q, want %q", decoded.Room, msg.Room)
	}
}

func TestFromJSON_Valid(t *testing.T) {
	original := NewMessage(ChatMessage, "alice", "hello", "general")
	original.Timestamp = "2024-01-01 12:00:00"

	data, err := original.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error: %v", err)
	}

	decoded, err := FromJSON(data)
	if err != nil {
		t.Fatalf("FromJSON() error: %v", err)
	}

	if decoded.Type != original.Type {
		t.Errorf("Type mismatch: got %v, want %v", decoded.Type, original.Type)
	}
	if decoded.Username != original.Username {
		t.Errorf("Username mismatch: got %q, want %q", decoded.Username, original.Username)
	}
	if decoded.Content != original.Content {
		t.Errorf("Content mismatch: got %q, want %q", decoded.Content, original.Content)
	}
}

func TestFromJSON_Invalid(t *testing.T) {
	_, err := FromJSON([]byte("not valid json"))
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestSenderFieldNotExportedToJSON(t *testing.T) {
	msg := NewMessage(ChatMessage, "alice", "hello", "general")
	msg.Sender = &Client{Username: "alice"} // should be excluded from JSON

	data, err := msg.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error: %v", err)
	}

	// Sender field tagged json:"-" so it must not appear in output
	if bytes := string(data); len(bytes) == 0 {
		t.Fatal("empty JSON output")
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if _, found := raw["sender"]; found {
		t.Error("Sender field should not be present in JSON output")
	}
}
