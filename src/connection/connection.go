package cz3004

// Message is a concrete message type
type Message struct {
	b []byte
}

// Connection is an abstract interface representing the possible conncetion
type Connection interface {
	Send([]byte) (n int, e error)
	Receive([]byte) (m Message)
}

// AlgoConnection ...
type AlgoConnection struct {
	b []byte
}

func (a *AlgoConnection) Send([]byte) (n int, e error) {

}
