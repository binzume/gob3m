package main

import (
	"github.com/binzume/gob3m/b3m"
	"github.com/tarm/serial"
	"encoding/binary"
	"os"
	"os/signal"
	"syscall"
	"flag"
	"log"
	"time"
)

type MotorStatus struct {
	Temperature float32
	Current float32
	Voltage float32
	DutyRatio float32
}

func GetMotorStatus(servo *b3m.Servo) (*MotorStatus, error) {
	data, err := servo.ReadMem(0x46, 10)
	if err != nil {
		return nil, err
	}
	t := float32(binary.LittleEndian.Uint16(data[0:2])) / 100
	c := float32(binary.LittleEndian.Uint16(data[2:4])) / 1000
	v := float32(binary.LittleEndian.Uint16(data[4:6])) / 1000
	d := float32(binary.LittleEndian.Uint16(data[6:8])) / float32(binary.LittleEndian.Uint16(data[8:10]))
	return &MotorStatus{t, c, v, d}, nil
}

func main() {
	var opt_port = flag.String("port", "COM1", "Serial port")
	var opt_id = flag.Int("id", 0, "servo id")
	flag.Parse()

	s, err := serial.OpenPort(&serial.Config{Name: *opt_port, Baud: 1500000})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		s.Close()
	}()
	id := byte(*opt_id)
	conn := b3m.New(s)

	servo := conn.GetServo(id)

	err = servo.SetMode(b3m.ControlPosition | b3m.RunNormal)
	if err != nil {
		log.Fatal(err)
	}

	pos, err := servo.GetCurrentPosition()

	err = servo.SetPosition(pos)
	if err != nil {
		log.Fatal(err)
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)

	for {
		ms, err := GetMotorStatus(servo)
		if err != nil {
			log.Printf("error %v", err)
		} else {
			pos, _ = servo.GetCurrentPosition()
			log.Printf("Motor: T:%v\tI:%v(A)\tV:%v(V)\tD:%v", ms.Temperature, ms.Current, ms.Voltage, ms.DutyRatio, pos)
		}
		time.Sleep(20 * time.Millisecond)
		select {
		case s := <- sigch:
			_ = servo.SetMode(b3m.RunFree)
			log.Printf("exit %v", s)
			return
		default:
		}
	}
}
