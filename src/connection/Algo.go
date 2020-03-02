package connection

import (
	"CZ3004-RPi/src/message"
	"bytes"
	"net"
)

// byte offset at which we have to split into algo/arduino
const offset = 10

// AlgoConnection ...
// when do we close the channel? or leave open?
type AlgoConnection struct {
	conn  *net.TCPConn         // represents the bytestream
	toRPi chan message.Request // messages from algo to rpi
}

// Receive an outgoing message  from rpi and send to conn without expecting reply
func (a *AlgoConnection) Receive(m message.Message) (n int, e error) {
	n, e = a.conn.Write(m.Buf.Bytes())
	if e != nil {
		return n, e
	}
	return n, nil
}

// Send request to rpi from your own internal channel
func (a *AlgoConnection) Send(b []byte) (n int, e error) {
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

// NewAlgo ...
func NewAlgo(c chan message.Request) *AlgoConnection {
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
		return &AlgoConnection{conn: actual, toRPi: c}
	}
}
