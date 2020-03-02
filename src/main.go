package main

import (
	"CZ3004-RPi/src/connection"
	"CZ3004-RPi/src/handler"
	"CZ3004-RPi/src/message"
	"CZ3004-RPi/src/rpi"
	"bufio"
	"io"
)

// ENDL
const ENDL rune = '\n'

func main() {
	/*
		TODO: how to handle closing channels
		TOOO: initial initialization
		TODO: link conn and rpi

		set up and initialize the rpi module
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
	/*
		AlgoConn := connection.NewAlgo(rpi.Requests)
		AndroidConn := connection.NewAndroid(rpi.Requests)
		ArduinoConn, _ := connection.NewArduino("8080", 8, rpi.Requests)
	*/
	MockAlgo := connection.Connection{&connection.MockConn{"1234\n", true, "algo"}, rpi.Requests, message.Algo}
	MockAndroid := connection.Connection{&connection.MockConn{"2345\n", true, "android"}, rpi.Requests, message.Android}
	MockArduino := connection.Connection{&connection.MockConn{"3456\n", true, "arduino"}, rpi.Requests, message.Arduino}

	rpi.RegisterReceivers(MockAlgo.Receive, message.Algo)
	rpi.RegisterReceivers(MockAndroid.Receive, message.Android)
	rpi.RegisterReceivers(MockArduino.Receive, message.Arduino)

	go initialize(MockAlgo)
	go initialize(MockArduino)
	go initialize(MockAndroid)
	for i := range rpi.Requests {
		rpi.Get(i)
	}
	for {

	}
}

func initialize(c connection.Connection) {
	for {
		r, e := bufio.NewReader(&c).ReadString('\n')
		if e != nil && e != io.EOF {
			continue
		}
		c.Send([]byte(r))
	}
}
