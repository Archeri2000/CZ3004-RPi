package connection

import (
	"CZ3004-RPi/src/message"

	"github.com/jacobsa/go-serial/serial"
)

// ArduinoConnection ...

// NewArduino ...
func NewArduino(port string, baud uint, toRPi chan message.Request) *Connection {
	options := serial.OpenOptions{
		PortName:        port,
		BaudRate:        baud,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}
	for {
		conn, err := serial.Open(options)
		if err != nil {
			continue
		}
		return &Connection{conn, toRPi, message.Arduino}
	}
}
