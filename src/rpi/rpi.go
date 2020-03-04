package rpi

import (
	"CZ3004-RPi/src/handler"
	"CZ3004-RPi/src/message"
	"bytes"
	"strconv"
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
const discard = 'x'

// Get is a abstraction of a client submitting a request to rpi
// this just calls the handler
func (rpi *RPi) Get(r message.Request) {
	go rpi.incomingHandlers[r.Kind](r)
	return
}

// AlgoHandler handles incoming messages from Algo conn
func (rpi *RPi) AlgoHandler(r message.Request) {
	switch r.Header {
	case message.Move:
		// Split for ardu
		arduinoBytes := make([]byte, offset)
		r.M.Buf.Read(arduinoBytes)
		arduinoBytes = append([]byte(strconv.Itoa(int(message.Move))), arduinoBytes...)
		arduinoMessage := message.Message{Buf: bytes.NewBuffer(arduinoBytes)}
		rpi.outgoingReceivers[message.Arduino](arduinoMessage)
		// Split for android
		// assumption - algo adds the pipe separator
		androidBytes := r.M.Buf.Bytes()
		androidMessage := message.Message{Buf: bytes.NewBuffer(androidBytes)}
		rpi.outgoingReceivers[message.Android](androidMessage)
		r.Result <- <-rpi.toAlgo
	case message.FastestPath:
		fastestPath := r.M.Buf.Bytes()                                                       // grab byte array representing moves
		fastestPath = append([]byte(strconv.Itoa(int(message.FastestPath))), fastestPath...) // assumption - moves can be broken into bytes
		arduinoMessage := message.Message{Buf: bytes.NewBuffer(fastestPath)}
		rpi.outgoingReceivers[message.Arduino](arduinoMessage)
	case message.Calibration:
		// request from algo for calibration - route to arduino
		arduinoBytes := r.M.Buf.Bytes()
		arduinoBytes = append([]byte(strconv.Itoa(int(message.Calibration))), arduinoBytes...)
		arduinoMessage := message.Message{bytes.NewBuffer(arduinoBytes)}
		rpi.outgoingReceivers[message.Arduino](arduinoMessage)
	}
	close(r.Result)
}

// AndroidHandler handles incoming misc messages from android conn
func (rpi *RPi) AndroidHandler(r message.Request) {
	// append \n to exploration/setwaypoint
	switch r.Header {
	// implicit assumption to do calibration
	case message.FastestPathStart:
		algoBytes := []byte{'\n'}
		algoBytes = append([]byte(strconv.Itoa(int(message.FastestPathStart))), algoBytes...)
		algoMessage := message.Message{Buf: bytes.NewBuffer(algoBytes)}
		rpi.outgoingReceivers[message.Arduino](algoMessage) // only fp start routes to ardu
	case message.ExplorationStart:
		arduinoBytes := []byte{'\n'}
		arduinoBytes = append([]byte(strconv.Itoa(int(message.ExplorationStart))), arduinoBytes...)
		arduinoMessage := message.Message{Buf: bytes.NewBuffer(arduinoBytes)}
		rpi.outgoingReceivers[message.Algo](arduinoMessage) // exploration start + waypoint start routes to algo
	case message.SetWaypoint:
		algoBytes := r.M.Buf.Bytes()
		algoBytes = append([]byte(strconv.Itoa(int(message.SetWaypoint))), algoBytes...)
		algoMessage := message.Message{Buf: bytes.NewBuffer(algoBytes)}
		rpi.outgoingReceivers[message.Algo](algoMessage) // exploration start + waypoint start routes to algo
	}
	close(r.Result)
}

// ArduinoHandler handles incoming sensor input from arduino conn
func (rpi *RPi) ArduinoHandler(r message.Request) {
	// format data here
	var leftShort byte
	if leftShort, _ = r.M.Buf.ReadByte(); leftShort == discard {
		leftShort, _ = r.M.Buf.ReadByte()
	} else {
		_, _ = r.M.Buf.ReadByte()
	}
	algoBytes := append([]byte(strconv.Itoa(int(message.Sensor))), leftShort)
	algoBytes = append(algoBytes, r.M.Buf.Bytes()...)
	algoMessage := message.Message{bytes.NewBuffer(algoBytes)}
	rpi.toAlgo <- algoMessage // new message with formatted data not r.m
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
