package connection

import (
	"CZ3004-RPi/src/message"
	"bytes"
	"io"
	"net"
	"os"
)

// AndroidConnection ...
type AndroidConnection struct {
	conn  io.ReadWriteCloser   // change to bluetooth conn
	toRPi chan message.Request // messages from algo to rpi
}

// Receive an outgoing message and send without expecting reply
func (a *AndroidConnection) Receive(m message.Message) (n int, e error) {
	n, e = a.conn.Read(m.Buf.Bytes())
	if e != nil {
		return n, e
	}
	return n, nil
}

// Send request to rpi from your own internal channel
func (a *AndroidConnection) Send(b []byte) (n int, e error) {
	// check if a.toRPi is nil
	// wrap data
	m := message.Message{Buf: bytes.NewBuffer(b)}
	r := message.Request{Kind: message.Algo, M: m} // don't initialise the result channel
	a.toRPi <- r
	temp := <-r.Result
	if r.Result != nil {
		a.conn.Write(temp.Buf.Bytes())
	}
	return n, nil
}

// NewAndroid ...
func NewAndroid(c chan message.Request) *AndroidConnection {
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
		return &AndroidConnection{conn: actual, toRPi: c}
	}
}

// Mock returns a mock andriod conn that writes to stdout
func Mock(c chan message.Request) *AndroidConnection {
	return &AndroidConnection{conn: os.Stdout, toRPi: c}
}
