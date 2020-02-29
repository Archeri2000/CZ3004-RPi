package message

import "bytes"

// Message is a message sent over connection
type Message struct {
	Buf *bytes.Buffer // should this be bytes.Buffer????
}

// Request is an object representing an arbitrary request sent from a connection for a message
type Request struct {
	Kind   Kind         // represents the Connection sending the request
	M      Message      // represents the actual output from a connection (Arduino: sensor values/Algo: next move)
	Result chan Message // channel with the finalized message to send back to the respective channels (Arduino: next move/Algo: sensor values)
}

// Kind represents the kind of channel
type Kind int

const (
	// Algo ...
	Algo Kind = iota
	// Arduino ...
	Arduino
	// Android ...
	Android
)
