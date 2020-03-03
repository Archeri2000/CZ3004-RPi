package rpi

import (
	"CZ3004-RPi/src/handler"
	"CZ3004-RPi/src/message"
	"bytes"
	"os"
)

// RPi represents the rpi multiplexer
// multiplexes over 4 channels so idk - better way???
type RPi struct {
	Requests          chan message.Request              // incoming requests from all 4 channels
	toAlgo            chan message.Message              // a completed op for algo
	toAndroid         chan message.Message              // a completed op for android
	toArduino         chan message.Message              // a completed op for arduino
	incomingHandlers  map[message.Kind]handler.Handler  // stores incoming handlers
	outgoingReceivers map[message.Kind]handler.Receiver // stores outgoing handlers - wrapper over connections
}

const offset = 1 // byte offset between ard/android message

// Get is a abstraction of a client submitting a request to rpi
// this just calls the handler
func (rpi *RPi) Get(r message.Request) {
	go rpi.incomingHandlers[r.Kind](r)
	return
}

// AlgoHandler ...
func (rpi *RPi) AlgoHandler(r message.Request) {
	arduinoBytes := make([]byte, offset)
	r.M.Buf.Read(arduinoBytes)
	arduinoMessage := message.Message{Buf: bytes.NewBuffer(arduinoBytes)}
	androidBytes := r.M.Buf.Bytes()
	androidMessage := message.Message{Buf: bytes.NewBuffer(androidBytes)}
	rpi.outgoingReceivers[message.Arduino](arduinoMessage)
	rpi.outgoingReceivers[message.Android](androidMessage)
	r.Result <- <-rpi.toAlgo
	close(r.Result)
}

// AndroidHandler handles incoming misc messages from arduino conn
func (rpi *RPi) AndroidHandler(r message.Request) {
	androidBytes := r.M.Buf.Bytes()
	androidMessage := message.Message{Buf: bytes.NewBuffer(androidBytes)}
	os.Stdout.Write(androidBytes)
	rpi.outgoingReceivers[message.Android](androidMessage)
}

// ArduinoHandler handles incoming sensor input from arduino conn
func (rpi *RPi) ArduinoHandler(r message.Request) {
	// format data here
	rpi.toAlgo <- r.M // new message with formatted data not r.m
	close(r.Result)
}

// RegisterHandler registers a given handler to the internal handler hashmap of rpi
func (rpi *RPi) RegisterHandler(h handler.Handler, m message.Kind) {
	rpi.incomingHandlers[m] = h
}

// RegisterReceivers ...
func (rpi *RPi) RegisterReceivers(r handler.Receiver, m message.Kind) {
	rpi.outgoingReceivers[m] = r
}

// NewRPi returns a new RPi
func NewRPi() (rpi *RPi) {
	return &RPi{Requests: make(chan message.Request), toAlgo: make(chan message.Message), toAndroid: make(chan message.Message), toArduino: make(chan message.Message), incomingHandlers: make(map[message.Kind]handler.Handler), outgoingReceivers: make(map[message.Kind]handler.Receiver)}
}
