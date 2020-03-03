// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package connection

import (
	"CZ3004-RPi/src/message"
	"fmt"
	"golang.org/x/sys/unix"
)

// AndroidConnection ...

// NewAndroid ...
func NewAndroid(port string, baud uint, toRPi chan message.Request) *Connection {
	fd, _ := Socket(AF_BLUETOOTH, SOCK_STREAM, BTPROTO_RFCOMM)
	_ = unix.Bind(fd, &unix.SockaddrRFCOMM{
		Channel: 1,
		Addr:    [6]uint8{0, 0, 0, 0, 0, 0}, // BDADDR_ANY or 00:00:00:00:00:00
	})
	_ = Listen(fd, 1)
	nfd, sa, _ := Accept(fd)
	fmt.Printf("conn addr=%v fd=%d", sa.(*unix.SockaddrRFCOMM).Addr, nfd)
	Read(nfd, buf)
	for {
		if err != nil {
			continue
		}
		return &Connection{conn, toRPi, message.Android}
	}
}
