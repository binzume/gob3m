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
	var opt_pos = flag.Int("pos", 0, "position")
	flag.Parse()

	s, err := serial.OpenPort(&serial.Config{Name: *opt_port, Baud: 1500000})
	if err != nil {
		log.Fatal(err)
	}
	id := byte(*opt_id)
	pos := int16(*opt_pos)

	conn := b3m.New(s)

	servo := conn.GetServo(id)

	err = servo.SetMode(b3m.ControlPosition | b3m.RunNormal)
	if err != nil {
		log.Fatal(err)
	}

	err = servo.SetPosition(pos)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ok")
}
