package connection

import (
	"CZ3004-RPi/src/message"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

// Connection is a struct representing the possible connection force clients to implement a connection interface; send/rx is for multiplexing with other goroutines
type Connection struct {
	io.ReadWriteCloser                      // represents the bytestream
	ToRPi              chan message.Request // messages from algo to rpi
	Kind               message.Kind
}

// Receive an outgoing message from rpi and send to conn without expecting reply
func (conn *Connection) Receive(m message.Message) (n int, e error) {
	n, e = conn.Write(m.Buf.Bytes())
	if e != nil {
		return n, e
	}
	return n, nil
}

// Send request to rpi from your own internal channel
func (conn *Connection) Send(b []byte) (n int, e error) {
	// check if a.toRPi is nil
	// wrap data
	m := message.Message{Buf: bytes.NewBuffer(b[1:])}
	head, _ := strconv.Atoi(string(string(b)[0]))
	fmt.Printf("Header byte: %d\n", head)
	r := message.Request{Kind: conn.Kind, M: m, Result: make(chan message.Message), Header: message.Header(head)} // don't initialise the result channel
	conn.ToRPi <- r
	temp, ok := <-r.Result
	if ok {
		conn.Write(temp.Buf.Bytes())
	}
	return n, nil
}
