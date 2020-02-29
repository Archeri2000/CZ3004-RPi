package connection

import (
	"CZ3004-RPi/src/message"
	"net"
)

// byte offset at which we have to split into algo/arduino
const offset = 10

// AlgoConnection ...
// when do we close the channel? or leave open?
type AlgoConnection struct {
	conn  net.TCPConn          // represents the bytestream
	toRPi chan message.Message // messages from rpi to algo
}

// Receive a message from rpi and services that message
func (a *AlgoConnection) Receive(m message.Message) (n int, e error) {
	n, e = a.conn.Read(m.Buf.Bytes())
	if e != nil {
		return n, e
	}
	return n, nil
}

// Send request to rpi from your own internal channel
func (a *AlgoConnection) Send(c chan message.Request) (n int, e error) {
	// check if a.toRPi is nil
	m, ok := <-a.toRPi
	if !ok { // channel already closed
		return 0, nil
	}
	r := message.Request{Kind: message.Algo, M: m} // don't initialise the result channel
	c <- r
	temp := <-r.Result
	a.conn.Write(temp.Buf.Bytes())
	return n, nil
}
