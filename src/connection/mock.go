package connection

import (
	"fmt"
	"io"
	"os"
)

// MockConn ...
type MockConn struct {
	TestValue string
	CanRead   bool
	Name      string
}

func (m *MockConn) Read(b []byte) (n int, e error) {
	if m.CanRead || true {
		copy(b, m.TestValue)
		m.CanRead = false
		return len(m.TestValue), nil
	}
	return 0, io.EOF
}

func (m *MockConn) Write(b []byte) (n int, e error) {
	os.Stdout.Write([]byte(fmt.Sprintf("%s:%c\n", m.Name, b)))
	return len(m.TestValue), nil
}

func (m *MockConn) Close() error {
	return nil
}
