// +build linux

package connection

import (
	"fmt"
	. "golang.org/x/sys/unix"
)

type BluetoothSocket struct {
	sockAddr Sockaddr
	nfr      int
}

func NewBluetoothSocket(sockAddr Sockaddr, nfr int) *BluetoothSocket {
	return &BluetoothSocket{sockAddr: sockAddr, nfr: nfr}

}

func (sock *BluetoothSocket) Read(p []byte) (n int, e error) {
	n, e = Read(sock.nfr, p)
	if e != nil {
		fmt.Printf("n: %d, e: %s", n, e)
	}
	return n, e
}

func (sock *BluetoothSocket) Write(p []byte) (n int, e error) {
	return Write(sock.nfr, p)
}

func (sock *BluetoothSocket) Close() error {
	return Close(sock.nfr)
}
