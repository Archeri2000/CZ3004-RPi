package connection

import (
	"CZ3004-RPi/src/message"
	"net"
)

// byte offset at which we have to split into algo/arduino
const offset = 10

// AlgoConnection ...
// NewAlgo ...
func NewAlgo(c chan message.Request) *Connection {
	t, e := net.ResolveTCPAddr("tcp4", ":9999")
	if e != nil {
		panic(e)
	}
	conn, e := net.ListenTCP("tcp", t)
	for {
		actual, e := conn.AcceptTCP()
		if e != nil {
			continue
		}
		return &Connection{actual, c, message.Algo}
	}
}
