package main

import (
	"github.com/binzume/gob3m/b3m"
	"github.com/tarm/serial"
	"flag"
	"time"
	"log"
)

func main() {
	var opt_port = flag.String("port", "COM1", "Serial port")
	var opt_id = flag.Int("id", 0, "servo id")
	var opt_newid = flag.Int("newid", 1, "servo NEW id")
	flag.Parse()

	id := byte(*opt_id)
	newid := byte(*opt_newid)

	s, err := serial.OpenPort(&serial.Config{Name: *opt_port, Baud: 1500000, ReadTimeout: 100 * time.Millisecond})
	if err != nil {
		log.Fatal(err)
	}

	servo := b3m.GetServo(s, id)

	err = servo.WriteMem(0, []byte{newid})
	if err != nil {
		log.Fatal(err)
	}
	servo.Id = newid

	err = servo.Save()
	if err != nil {
		log.Fatal(err)
	}

	err = servo.Reset(0)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ok")
}
