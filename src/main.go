package main

import (
	"CZ3004-RPi/src/connection"
	"CZ3004-RPi/src/handler"
	"CZ3004-RPi/src/message"
	"CZ3004-RPi/src/rpi"
	"bufio"
	"bytes"
	"fmt"
	"os"
)

// ENDL
const ENDL byte = 'a'

func main() {
	/*
		TODO: how to handle closing channels
		TOOO: initial initialization
		TODO: link conn and rpi

		set up and listenOn the rpi module
		set up individual connections

		go func() {
			persistently listen on connections and store in toRPi channel
		}

		main goroutine (this is the rpi) persistently runs get
	*/
	rpi := rpi.NewRPi()
	AlgoH := handler.Handler(rpi.AlgoHandler)
	rpi.RegisterHandler(AlgoH, message.Algo)
	AndroidH := handler.Handler(rpi.AndroidHandler)
	rpi.RegisterHandler(AndroidH, message.Android)
	ArduinoH := handler.Handler(rpi.ArduinoHandler)
	rpi.RegisterHandler(ArduinoH, message.Arduino)
	And := connection.NewAndroid(rpi.Requests)
	rpi.RegisterReceivers(And.Receive, message.Android)
	Ardu := connection.NewArduino("/dev/ttyACM0", 115200, rpi.Requests)
	rpi.RegisterReceivers(Ardu.Receive, message.Arduino)
	fmt.Printf("Success!")
	MockAlgo := connection.Connection{&connection.MockConn{"stesta", true, "algo"}, rpi.Requests, message.Algo}
	rpi.RegisterReceivers(MockAlgo.Receive, message.Algo)
	go listenOn(And)
	go listenOn(&MockAlgo)
	go listenOn(Ardu)
	for i := range rpi.Requests {
		rpi.Get(i)
	}
	os.Exit(0)
	/*
		AlgoConn := connection.NewAlgo(rpi.Requests)
		AndroidConn := connection.NewAndroid(rpi.Requests)
		ArduinoConn, _ := connection.NewArduino("8080", 8, rpi.Requests)
	*/
	MockAndroid := connection.Connection{&connection.MockConn{"2345\n", true, "android"}, rpi.Requests, message.Android}
	MockArduino := connection.Connection{&connection.MockConn{"3456\n", true, "arduino"}, rpi.Requests, message.Arduino}

	rpi.RegisterReceivers(MockAlgo.Receive, message.Algo)
	rpi.RegisterReceivers(MockAndroid.Receive, message.Android)
	rpi.RegisterReceivers(MockArduino.Receive, message.Arduino)

	go listenOn(&MockAlgo)
	go listenOn(&MockArduino)
	go listenOn(&MockAndroid)
	for i := range rpi.Requests {
		rpi.Get(i)
	}
}

func listenOn(c *connection.Connection) {
	buf := bytes.Buffer{}
	reader := bufio.NewReader(c)
	for {
		r, e := reader.ReadString(ENDL)
		buf.Write([]byte(r))
		fmt.Printf("%d\n", buf.Len())
		if e == nil {
			_, _ = c.Send(buf.Bytes())
			buf = bytes.Buffer{}
		}
	}
}
