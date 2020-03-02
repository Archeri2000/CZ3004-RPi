package connection

import (
	"CZ3004-RPi/src/message"
	"net"
)

// AndroidConnection ...

// NewAndroid ...
func NewAndroid(c chan message.Request) *Connection {
	t, e := net.ResolveTCPAddr("tcp4", ":9998")
	if e != nil {
		panic(e)
	}
	conn, e := net.ListenTCP("tcp", t)
	for {
		actual, e := conn.AcceptTCP()
		if e != nil {
			continue
		}
		return &Connection{actual, c, message.Android}
	}
}
