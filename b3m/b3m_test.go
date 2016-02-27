package b3m

import (
	"testing"
	"bytes"
)

func TestReset(t *testing.T) {
	io := bytes.NewBuffer([]byte{})
	conn := New(io)
	servo := conn.GetServo( 0)
	servo.Reset(0)
	if bytes.Compare([]byte{6,5,0,0,0,11}, io.Bytes()) != 0 {
		t.Fatal("command not matched %v", io.Bytes())
	}
}
