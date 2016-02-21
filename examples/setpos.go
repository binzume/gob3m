package main

import (
	"github.com/binzume/gob3m/b3m"
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

	err = servo.SetMode(b3m.ControlPosition | b3m.RunNormal)
	if err != nil {
		log.Fatal(err)
	}

	err = servo.SetPosition(500)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ok")
}
