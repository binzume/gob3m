package main

import (
	"log"
	"github.com/tarm/serial"
	"./b3m"
)

func main() {
	c := &serial.Config{Name: "COM3", Baud: 1500000}
	s, err := serial.OpenPort(c)
	if err != nil {
	        log.Fatal(err)
	}
	var id byte = 0


	model, version, err := b3m.GetVersion(s, id)
	if err != nil {
	        log.Fatal(err)
	} else {
		log.Printf("Model:%v Version:%v", model, version)
	}

	// get id test
	res, err := b3m.ReadMem(s, id, 0, 1, 100)
	if err != nil {
	        log.Fatal(err)
	} else {
		log.Printf("%v", res.Data)
	}


	mode, err := b3m.GetMode(s, id)
	if err != nil {
	        log.Fatal(err)
	} else {
		log.Printf("mode %v", mode)
	}

	pos, err := b3m.GetCurrentPosition(s, id)
	if err != nil {
	        log.Fatal(err)
	} else {
		log.Printf("pos: %v", pos)
	}

	b3m.SetPosition(s, id, pos + 100)

	_, err = b3m.SetMode(s, id, b3m.ModeP | b3m.ModeNormal)
	if err != nil {
        log.Fatal(err)
	}

	// b3m.Reset(s, id, 1)

}
