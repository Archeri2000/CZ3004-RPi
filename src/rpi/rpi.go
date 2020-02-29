package rpi

import (
	"CZ3004-RPi/src/message"
	"bytes"
)

// RPi represents the rpi multiplexer
// multiplexes over 4 channels so idk - better way???
type RPi struct {
	requests  chan message.Request // incoming requests from all 4 channels
	toAlgo    chan message.Message // a completed op for algo
	toAndroid chan message.Message // a completed op for android
	toArduino chan message.Message // a completed op for arduino
}

const offset = 10 // byte offset between ard/android message

// Get is a abstraction of a client submitting a request to rpi
// this just calls the handler
// can implement a handler interface also
func (rpi *RPi) Get(r message.Request) (m message.Message) {
	switch r.Kind {
	case message.Algo:
		go rpi.AlgoHandler(r)
		return <-r.Result
	case message.Android:
		go rpi.AndroidHandler(r)
		return <-r.Result
	case message.Arduino:
		go rpi.ArduinoHandler(r)
		return <-r.Result
	}
}

// AlgoHandler ...
func (rpi *RPi) AlgoHandler(r message.Request) {
	arduinoBytes := make([]byte, offset)
	arduinoMessage := message.Message{Buf: bytes.NewBuffer(arduinoBytes)}
	androidBytes := r.M.Buf.Bytes()
	androidMessage := message.Message{Buf: bytes.NewBuffer(androidBytes)}
	rpi.toArduino <- arduinoMessage
	rpi.toAndroid <- androidMessage
	m := <-rpi.toAlgo
	r.Result <- m
	close(r.Result)
}
