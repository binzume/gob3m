package b3m

import (
	"errors"
	"fmt"
	"io"
	"time"
)

const DefaultTimeout = 100

// commands
type CommandType byte

const CmdLoad CommandType = 1
const CmdSave CommandType = 2
const CmdRead CommandType = 3
const CmdWrite CommandType = 4
const CmdReset CommandType = 5
const CmdPosition CommandType = 6

// error status
const StatusSystemError = 1
const StatusMotorError = 2
const StatusUartError = 4
const StatusCommandError = 8

// servo modes
const RunNormal byte = 0
const RunFree byte = 2
const RunHold byte = 3

// control modes
const ControlPosition = 0
const ControlVelocity = 4
const ControlTorque = 8
const ControlFForword = 12

// trajectory
type TrajectoryType byte

const TrajectoryNormal TrajectoryType = 0
const TrajectoryEven TrajectoryType = 1
const TrajectoryThirdPoly TrajectoryType = 3
const TrajectoryFourthPoly TrajectoryType = 4
const TrajectoryFifthPoly TrajectoryType = 5

type Command struct {
	Cmd    CommandType
	Option byte
	Id     byte
	Data   []byte
}

func Send(s io.Writer, c *Command) (int, error) {
	buf := make([]byte, len(c.Data)+5)
	buf[0] = (byte)(len(c.Data) + 5)
	buf[1] = (byte)(c.Cmd)
	buf[2] = c.Option
	buf[3] = c.Id
	copy(buf[4:], c.Data)
	var sum byte = 0
	for i := 0; i < len(buf)-1; i++ {
		sum += buf[i]
	}
	buf[len(buf)-1] = sum
	return s.Write(buf)
}

func Recv(s io.Reader) (*Command, error) {
	buf := make([]byte, 256)
	n, err := s.Read(buf[0:1])
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, errors.New("timeout")
	}
	sz := (int)(buf[0])
	for i := 1; i < sz; {
		n, err = s.Read(buf[i:])
		if err != nil {
			return nil, err
		}
		i += n
	}
	data := make([]byte, sz-5)
	copy(data, buf[4:])
	cmd := &Command{(CommandType)(buf[1]), buf[2], buf[3], data}
	return cmd, nil
}

func Recv2(s io.Reader, timeout int) (*Command, error) {
	type Result struct {
		value *Command
		err   error
	}
	ch := make(chan Result, 1)
	go func() {
		ret, err := Recv(s)
		ch <- Result{ret, err}
	}()

	select {
	case ret := <-ch:
		return ret.value, ret.err
	case <-time.After(time.Millisecond * time.Duration(timeout)):
		return nil, errors.New("timeout")
	}
}

func ReadMem(s io.ReadWriter, id byte, addr int, size int, timeout int) (*Command, error) {
	cmd := &Command{CmdRead, 0, id, []byte{(byte)(addr), (byte)(size)}}
	_, err := Send(s, cmd)
	if err != nil {
		return nil, err
	}
	return Recv2(s, timeout)
}

func WriteMem(s io.ReadWriter, id byte, addr int, data []byte, timeout int) (*Command, error) {
	buf := make([]byte, len(data)+2)
	copy(buf, data)
	buf[len(buf)-2] = (byte)(addr)
	buf[len(buf)-1] = 1
	cmd := &Command{CmdWrite, 0, id, buf}
	_, err := Send(s, cmd)
	if err != nil {
		return nil, err
	}
	if id == 255 {
		return cmd, nil
	}
	return Recv2(s, timeout)
}

type Servo struct {
	io        io.ReadWriter // serial port
	Id        byte          // device id
	TimeoutMs int           // timeout for replay.
	Status    byte          // last status
}

func GetServo(io io.ReadWriter, id byte) *Servo {
	return &Servo{io, id, DefaultTimeout, 0}
}

func (s *Servo) ReadMem(addr int, size int) ([]byte, error) {
	res, err := ReadMem(s.io, s.Id, addr, size, s.TimeoutMs)
	if err != nil {
		return nil, err
	}
	s.Status = res.Option
	return res.Data, nil
}

func (s *Servo) WriteMem(addr int, data []byte) error {
	res, err := WriteMem(s.io, s.Id, addr, data, s.TimeoutMs)
	if err != nil {
		return err
	}
	s.Status = res.Option
	return nil
}

func (s *Servo) GetVersion() (model string, version string, err error) {
	buf, err := s.ReadMem(0xA2, 12)
	if err != nil {
		return "", "", err
	}
	model = fmt.Sprintf("B3M-%c%c-%v%v%v-%c", buf[7], buf[6], buf[3], buf[2], buf[1], buf[0])
	version = fmt.Sprintf("%v.%v.%v.%v", buf[11], buf[10], buf[9], buf[8])
	return
}

func (s *Servo) GetMode() (byte, error) {
	buf, err := s.ReadMem(0x28, 1)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (s *Servo) SetMode(mode byte) error {
	return s.WriteMem(0x28, []byte{mode})
}

func (s *Servo) Reset(timeAfter byte) error {
	cmd := &Command{CmdReset, 0, s.Id, []byte{timeAfter}}
	_, err := Send(s.io, cmd)
	return err
}

func (s *Servo) Load() error {
	_, err := Send(s.io, &Command{CmdLoad, 0, s.Id, []byte{}})
	if err != nil {
		return err
	}
	if s.Id != 255 {
		_, err = Recv2(s.io, s.TimeoutMs)
	}
	return err
}

func (s *Servo) Save() error {
	_, err := Send(s.io, &Command{CmdSave, 0, s.Id, []byte{}})
	if err != nil {
		return err
	}
	if s.Id != 255 {
		_, err = Recv2(s.io, s.TimeoutMs)
	}
	return err
}

func (s *Servo) SetTrajectoryMode(trajectory TrajectoryType) error {
	return s.WriteMem(0x29, []byte{byte(trajectory)})
}

func (s *Servo) SetPosition(pos int16) error {
	return s.WriteMem(0x2A, []byte{(byte)(pos), (byte)(pos >> 8)})
}

func (s *Servo) GetCurrentPosition() (int16, error) {
	res, err := s.ReadMem(0x2C, 2)
	if err != nil {
		return 0, err
	}
	return (int16)(res[0]) | ((int16)(res[1]) << 8), nil
}

func (s *Servo) SetVelocity(v int16) error {
	return s.WriteMem(0x30, []byte{(byte)(v), (byte)(v >> 8)})
}

func (s *Servo) SetTorque(torque int16) error {
	return s.WriteMem(0x3C, []byte{(byte)(torque), (byte)(torque >> 8)})
}

func (s *Servo) SetPosition2(pos, time int16) error {
	_, err := Send(s.io, &Command{CmdPosition, 0, s.Id, []byte{(byte)(pos), (byte)(pos >> 8), (byte)(time), (byte)(time >> 8)}})
	if err != nil {
		return err
	}
	if s.Id != 255 {
		_, err = Recv2(s.io, s.TimeoutMs)
	}
	return err
}
