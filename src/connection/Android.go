// +build linux

package connection

import (
	"CZ3004-RPi/src/message"
	"fmt"
	. "golang.org/x/sys/unix"
)

// AndroidConnection ...

// NewAndroid ...
func NewAndroid(toRPi chan message.Request) *Connection {
	fd, _ := Socket(AF_BLUETOOTH, SOCK_STREAM, BTPROTO_RFCOMM)
	_ = Bind(fd, &SockaddrRFCOMM{
		Channel: 2,
		Addr:    [6]uint8{0, 0, 0, 0, 0, 0}, // BDADDR_ANY or 00:00:00:00:00:00
	})
	_ = Listen(fd, 1)
	//fmt.Printf("listening %s\n", er)
	//fmt.Printf("file descriptor: %d\n", fd)
	for {
		nfd, sa, err := Accept(fd)
		fmt.Printf("conn addr=%v fd=%d", sa.(*SockaddrRFCOMM).Addr, nfd)
		if err != nil {
			continue
		}
		return &Connection{NewBluetoothSocket(sa, nfd), toRPi, message.Android}
	}
}
