package main

import (
	"./b3m"
	"github.com/tarm/serial"
	"log"
)

func main() {
	s, err := serial.OpenPort(&serial.Config{Name: "COM3", Baud: 1500000})
	if err != nil {
		log.Fatal(err)
	}
	var id byte = 0

	servo := b3m.GetServo(s, id)

	model, version, err := servo.GetVersion()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Model:%v Version:%v", model, version)

	mode, err := servo.GetMode()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("mode %v", mode)

	pos, err := servo.GetCurrentPosition()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("pos: %v", pos)

	err = servo.SetMode(b3m.ControlPosition | b3m.RunNormal)
	if err != nil {
		log.Fatal(err)
	}

	err = servo.SetPosition(pos + 200)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ok")

	// servo.Reset(0)

}
