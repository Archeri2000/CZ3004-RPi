package connection

import (
	"CZ3004-RPi/src/message"
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
)

func NewImage(c chan message.Request) *Connection {
	actual := os.Stdout
	return &Connection{actual, c, message.Image}
}

func (conn *Connection) ImgReceive(m message.Message) (n int, e error) {
	// get o, x, y from img, submit to pyscript for approval
	predictionInfo := m.Buf.Bytes()
	result := getImageRec(int(predictionInfo[0]), int(predictionInfo[1]), int(predictionInfo[2]))
	return result["result"], nil
}

// gets orientation + co-ordinates, sends to pyscript and receives result
func getImageRec(o, x, y int) map[string]int {
	// translate orientation to actual ints
	cmd := exec.Command("python", "cmd goes here", strconv.Itoa(o), strconv.Itoa(x), strconv.Itoa(y))
	output, err := cmd.Output()
	if err != nil {
		println(err) // bad handling but who cares
	}
	jsonMap := make(map[string]int)
	err = json.Unmarshal([]byte(output), &jsonMap)
	if err != nil {
		panic(err)
	}
	return jsonMap
}
