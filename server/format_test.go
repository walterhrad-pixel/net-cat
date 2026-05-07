package server_test

import (
	"testing"

	"netcat/server"
)

func TestFormatMessage(t *testing.T) {
	msg := server.Message{
		Type:      server.ChatMessage,
		Username:  "Ryan",
		Content:   "hello",
		Timestamp: "2026-05-07 10:00:00",
	}
	result := server.FormatMessage(msg)

	expected := "[2026-05-07 10:00:00][Ryan]: hello\n"

	if string(result) != expected {
		t.Errorf("expected %q but got %q", expected, string(result))
	}
}
