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
	"strconv"
)

// ENDL ...
const ENDL byte = '\n'

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
	ImageH := handler.Handler(rpi.ImageHandler)
	rpi.RegisterHandler(ImageH, message.Image)

	//TODO: Revert to real android
	//Andr := connection.NewAndroid(rpi.Requests)
	Andr := connection.Connection{&connection.MockConn{"2345\n", false, "android"}, rpi.Requests, message.Android}
	fmt.Printf("Android Connected!\n")
	Ardu := connection.NewArduino("/dev/ttyACM0", 115200, rpi.Requests)
	fmt.Printf("Arduino Connected!\n")
	Algo := connection.NewAlgo(rpi.Requests)
	fmt.Printf("Algo Connected!\n")
	Image := connection.NewImage(rpi.Requests)

	rpi.RegisterReceivers(Andr.Receive, message.Android)
	rpi.RegisterReceivers(Ardu.Receive, message.Arduino)
	rpi.RegisterReceivers(Algo.Receive, message.Algo)
	rpi.RegisterReceivers(handler.Receiver(Image.ImgReceive), message.Image)
	fmt.Printf("Success!\n")
	//TODO: Remove as android is supposed to provide this signal
	algoBytes := []byte{'\n'}
	algoBytes = append([]byte(strconv.Itoa(int(message.ExplorationStart))), algoBytes...)
	algoMessage := message.Message{Buf: bytes.NewBuffer(algoBytes)}
	_, _ = Algo.Receive(algoMessage) // exploration start + waypoint start routes to algo
	algoBytes2 := []byte{'\n'}
	algoBytes2 = append([]byte(strconv.Itoa(int(message.FastestPathStart))), algoBytes2...)
	algoMessage2 := message.Message{Buf: bytes.NewBuffer(algoBytes2)}
	_, _ = Algo.Receive(algoMessage2) // exploration start + waypoint start routes to algo

	go listenOn(&Andr)
	go listenOn(Algo)
	go listenOn(Ardu)
	for i := range rpi.Requests {
		rpi.Get(i)
	}
	os.Exit(0)
	//
	//MockAlgo := connection.Connection{&connection.MockConn{"stest\n", true, "algo"}, rpi.Requests, message.Algo}
	//MockAndroid := connection.Connection{&connection.MockConn{"2345\n", true, "android"}, rpi.Requests, message.Android}
	//MockArduino := connection.Connection{&connection.MockConn{"3456\n", true, "arduino"}, rpi.Requests, message.Arduino}
	//
	//rpi.RegisterReceivers(MockAlgo.Receive, message.Algo)
	//rpi.RegisterReceivers(MockAndroid.Receive, message.Android)
	//rpi.RegisterReceivers(MockArduino.Receive, message.Arduino)
	//
	//go listenOn(&MockAlgo)
	//go listenOn(&MockArduino)
	//go listenOn(&MockAndroid)
	//for i := range rpi.Requests {
	//	rpi.Get(i)
	//}
}

func listenOn(c *connection.Connection) {
	buf := bytes.Buffer{}
	reader := bufio.NewReader(c)
	for {
		r, e := reader.ReadString(ENDL)
		fmt.Printf("Channel %s: (%s)\n", strconv.Itoa(int(c.Kind)), r)
		buf.Write([]byte(r))
		//fmt.Printf("%d\n", buf.Len())
		//fmt.Printf("%d\n", buf.Len())
		if e == nil {
			_, _ = c.Send(buf.Bytes())
			buf = bytes.Buffer{}
		}
	}
}
