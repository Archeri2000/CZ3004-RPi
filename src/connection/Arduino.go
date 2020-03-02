package connection

import (
	"CZ3004-RPi/src/message"
	"bytes"
	"io"

	"github.com/jacobsa/go-serial/serial"
)

// ArduinoConnection ...
// when do we close the channel? or leave open?
type ArduinoConnection struct {
	conn  io.ReadWriteCloser   // represents the bytestream
	toRPi chan message.Request // messages from rpi to algo
}

// NewArduino ...
func NewArduino(port string, baud uint, toRPi chan message.Request) (*ArduinoConnection, error) {
	ardu := new(ArduinoConnection)
	options := serial.OpenOptions{
		PortName:        port,
		BaudRate:        baud,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}
	conn, err := serial.Open(options)
	if err != nil {
		return ardu, err
	}
	ardu.conn = conn
	ardu.toRPi = toRPi
	return ardu, nil
}

// Receive a message from rpi and services that message
func (a *ArduinoConnection) Receive(m message.Message) (n int, e error) {
	n, e = a.conn.Write(m.Buf.Bytes())
	if e != nil {
		return n, e
	}
	return n, nil
}

// Send request to rpi from your own internal channel
func (a *ArduinoConnection) Send(b []byte) (n int, e error) {
	// formatting
	m := message.Message{Buf: bytes.NewBuffer(b)}
	r := message.Request{Kind: message.Algo, M: m} // don't initialise the result channel
	a.toRPi <- r
	temp := <-r.Result
	if r.Result != nil {
		a.conn.Write(temp.Buf.Bytes())
	}
	return n, nil
}
