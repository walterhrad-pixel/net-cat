package utils

import "netcat/server"

func FormatMessage(msg server.Message) []byte {
	return []byte(msg.String() + "\n")
}
