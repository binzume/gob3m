package main

import (
	"github.com/binzume/gob3m/b3m"
	"github.com/tarm/serial"
	"log"
	"time"
	"flag"
)

func main() {
	var opt_port = flag.String("port", "COM1", "Serial port")
	flag.Parse()

	s, err := serial.OpenPort(&serial.Config{Name: *opt_port, Baud: 1500000, ReadTimeout: 100 * time.Millisecond})
	if err != nil {
		log.Fatal(err)
	}

	// scan all servo
	found := 0
	for id := 0; id < 256; id ++{
		servo := b3m.GetServo(s, (byte)(id))
		model, version, err := servo.GetVersion()
		if err != nil {
			log.Printf("id:%v %v", id, err)
		} else {
			log.Printf("id:%v Model:%v Version:%v", id, model, version)
			found ++
		}
	}
	log.Printf("ok found: %v", found)
	s.Close()
}
