package connection

import (
	"CZ3004-RPi/src/message"
	"bufio"
	"os/exec"
)

func NewImage(c chan message.Request) *Connection {
	// TODO: add command path
	cmd := exec.Command("python", "something.py")
	cmd.Start()
	pw, _ := cmd.StdinPipe()
	pr, _ := cmd.StdoutPipe()
	cmdline := cmdline{w: bufio.NewWriter(pw), r: bufio.NewReader(pr)}
	return &Connection{&cmdline, c, message.Image}
}

func (conn *Connection) ImgReceive(m message.Message) (n int, e error) {
	// get o, x, y from img, submit to pyscript for approval
	predictionInfo := m.Buf.Bytes()
	conn.getImageRec(int(predictionInfo[0]), int(predictionInfo[1]), int(predictionInfo[2]))
	return 1, nil // TODO: check this
}

// gets orientation + co-ordinates, sends to pyscript and receives result
func (conn *Connection) getImageRec(o, x, y int) {
	// translate orientation to actual ints
	conn.Write([]byte{byte(rune(o)), byte(rune(x)), byte(rune(y))})
}

// unexported - will be wrapped by newimage; this implements read/write
type cmdline struct {
	w *bufio.Writer
	r *bufio.Reader
}

// Read from stdout of the cmd
func (c *cmdline) Read(p []byte) (n int, err error) {
	line, _, err := c.r.ReadLine()
	if err != nil {
		return 0, err
	}
	n = copy(p, line)
	return n, nil
}

// write to stdin
func (c *cmdline) Write(p []byte) (n int, err error) {
	n, err = c.w.WriteString(string(p))
	if err != nil {
		return n, err
	}
	return n, err
}
