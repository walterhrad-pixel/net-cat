package server

func FormatMessage(msg Message) []byte {
	return []byte(msg.String() + "\n")
}
