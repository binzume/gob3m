package b3m

import (
	"fmt"
	"time"
	"errors"
	"github.com/tarm/serial"
)

const DefaultTimeout = 100

// commands
type CommandType byte
const CmdLoad CommandType = 1;
const CmdSave CommandType = 2;
const Read CommandType = 3;
const Write CommandType = 4;
const CmdReset CommandType = 5;
const CmdPosition CommandType = 6;

// servo modes
const ModeNormal byte = 0;
const ModeFree byte = 2;
const ModeHold byte = 3;

// control modes
const ModeP = 0;
const ModeV = 4;
const ModeT = 8;
const ModeFF = 12;

type Command struct {
	Cmd CommandType
	Option byte
	Id byte
	Data []byte
}

func Send(s *serial.Port, c *Command) (int, error)  {
	buf := make([]byte, len(c.Data) + 5)
	buf[0] = (byte)(len(c.Data) + 5);
	buf[1] = (byte)(c.Cmd)
	buf[2] = c.Option
	buf[3] = c.Id
	copy(buf[4:], c.Data)
	var sum byte = 0
	for i:= 0; i<len(buf)-1; i++ {
		sum += buf[i]
	}
	buf[len(buf)-1] = sum
	return s.Write(buf)
}

func Recv(s *serial.Port) (*Command, error) {
	buf := make([]byte, 256)
	n, err := s.Read(buf[0:1])
	if err != nil || n == 0{
		return nil, err
	}
	sz := (int)(buf[0])
	for i := 1; i < sz ; {
		n, err = s.Read(buf[i:])
		if err != nil {
			return nil, err
		}
		i += n
	}
	data := make([]byte, sz - 5)
	copy(data, buf[4:])
	cmd := &Command{(CommandType)(buf[1]), buf[2], buf[3],data}
	return cmd, nil
}

func ReadMem(s *serial.Port, id byte, addr int, size int, timeout int) (*Command, error) {
	cmd := &Command{Read, 0, id, []byte{(byte)(addr),(byte)(size)}}
	_, err := Send(s, cmd)
	if err != nil {
		return nil, err
	}

	type Result struct { value *Command; err error}
	ch := make(chan Result, 1)
	go func() {
		ret,err := Recv(s)
		ch <- Result{ret, err}
	}()

	select {
	case ret := <- ch:
		return ret.value, ret.err
	case <-time.After(time.Millisecond * time.Duration(timeout)):
		return nil, errors.New("timeout")
	}
}

func WriteMem(s *serial.Port, id byte, addr int, data []byte, timeout int) (*Command, error) {
	buf := make([]byte, len(data) + 2)
	copy(buf, data)
	buf[len(buf)-2] = (byte)(addr)
	buf[len(buf)-1] = 1
	cmd := &Command{Write, 0, id, buf}
	_, err := Send(s, cmd)
	if err != nil {
		return nil, err
	}

	type Result struct { value *Command; err error}
	ch := make(chan Result, 1)
	go func() {
		ret, err := Recv(s)
		ch <- Result{ret, err}
	}()

	select {
	case ret := <- ch:
		return ret.value, ret.err
	case <-time.After(time.Millisecond * time.Duration(timeout)):
		return nil, errors.New("timeout")
	}
}


func GetVersion(s *serial.Port, id byte) (model string, version string, err error) {
	res, err := ReadMem(s, id, 0xA2, 12, DefaultTimeout)
	if err != nil {
		return "", "", err
	}
	model = fmt.Sprintf("B3M-%c%c-%v%v%v-%c", res.Data[7], res.Data[6], res.Data[3], res.Data[2], res.Data[1], res.Data[0])
	version = fmt.Sprintf("%v.%v.%v.%v", res.Data[11], res.Data[10], res.Data[9], res.Data[8])
	return
}


func GetMode(s *serial.Port, id byte) (byte, error) {
	res, err := ReadMem(s, id, 0x28, 1, DefaultTimeout)
	if err != nil {
		return 0, err
	}
	return res.Data[0], nil
}


func SetMode(s *serial.Port, id byte, mode byte) (*Command, error) {
	res, err := WriteMem(s, id, 0x28, []byte{mode}, DefaultTimeout)
	return res, err
}

func Reset(s *serial.Port, id byte, time byte) error {
	cmd := &Command{CmdReset, 0, id, []byte{time}}
	_, err := Send(s, cmd)
	return err
}

func SetPosition(s *serial.Port, id byte, pos int16) (error) {
	_, err := WriteMem(s, id, 0x2A, []byte{(byte)(pos), (byte)(pos>>8)}, DefaultTimeout)
	return err
}

func GetCurrentPosition(s *serial.Port, id byte) (int16, error) {
	res, err := ReadMem(s, id, 0x2C, 2, DefaultTimeout)
	if err != nil {
		return 0, err
	}
	return (int16)(res.Data[0]) | ((int16)(res.Data[1]) << 8), nil
}

