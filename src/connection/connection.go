package connection

import (
	"CZ3004-RPi/src/message"
	"net"
)

// Connection is an abstract interface representing the possible connection force clients to implement a connection interface; send/rx is for multiplexing with other goroutines
type Connection interface {
	// Send a message from your conn to RPi
	Send(b []byte) (n int, e error)
	// Receive a message from RPi to your own conn
	Receive(r message.Request) (n int, e error)
	net.Conn
}
