package server

import (
	"encoding/json"
	"fmt"
)

//MessageType defines the type of message 
type MessageType int

const (
	ChatMessage MessageType = iota
	SystemMessage
	CommandMessage
	JoinMessage
	LeaveMessage
)

//Message represents a chat message 
type Message struct {
	Type      MessageType  `json:"type"`
	Username  string       `json:"username"`
	Content   string       `json:"content"`
	Room      string       `json:"room"`
	Timestamp string       `json:"timestamp"`
	Sender    *Client      `json:"-"`
}

//NewMessage creates a new message
func NewMessage(msgType MessageType, username, content, room string) Message {
	return Message{
		Type:     msgType,
		Username: username,
		Content:  content,
		Room:     room,
	}
}

//String returns formatted message string
func (m Message) String() string {
	switch m.Type { 
	case ChatMessage:
		return fmt.Sprintf("[%s][%s]: %s", m.Timestamp, m.Username, m.Content)
	case SystemMessage:
		return fmt.Sprintf("[%s][SYSTEM]: %s", m.Timestamp, m.Content)
	case JoinMessage:
		return fmt.Sprintf("[%s][SYSTEM]: %s has joined the room", m.Timestamp, m.Username)
	case LeaveMessage:
		return fmt.Sprintf("[%s][SYSTEM]: %s has left the room", m.Timestamp, m.Username)
	default:
		return fmt.Sprintf("[%s][%s]: %s", m.Timestamp, m.Username, m.Content)				
	}
}

//ToJSON converts message to JSON
func (m Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

//FromJSON converts mesage from JSON
func FromJSON(data []byte) (Message, error) {
	var msg Message 
	err := json.Unmarshal(data, &msg)
	return msg, err
}
