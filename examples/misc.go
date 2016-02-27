package main

import (
	"github.com/binzume/gob3m/b3m"
	"github.com/tarm/serial"
	"flag"
	"log"
)

func main() {
	var opt_port = flag.String("port", "COM1", "Serial port")
	var opt_id = flag.Int("id", 0, "servo id")
	flag.Parse()

	s, err := serial.OpenPort(&serial.Config{Name: *opt_port, Baud: 1500000})
	if err != nil {
		log.Fatal(err)
	}
	id := byte(*opt_id)
	conn := b3m.New(s)

	servo := conn.GetServo(id)

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

