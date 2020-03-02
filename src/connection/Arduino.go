package connection

import (
	"CZ3004-RPi/src/message"
	"github.com/jacobsa/go-serial/serial"
	"io"
)

// AlgoConnection ...
// when do we close the channel? or leave open?
type ArduinoConnection struct {
	conn  io.ReadWriteCloser   // represents the bytestream
	toRPi chan message.Message // messages from rpi to algo
}

func NewArduinoConnection(port string, baud uint, toRPi chan message.Message) (*ArduinoConnection, error) {
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
	n, e = a.conn.Read(m.Buf.Bytes())
	if e != nil {
		return n, e
	}
	return n, nil
}

// Send request to rpi from your own internal channel
func (a *ArduinoConnection) Send(c chan message.Request) (n int, e error) {
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
