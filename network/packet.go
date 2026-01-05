package network

import (
	"time"

	"github.com/s5i/tcam/enum"
)

// Packet represents a game packet with a time offset.
type Packet struct {
	Offset time.Duration
	Data   []byte
}

func (p Packet) OpCode() enum.OpCode {
	if len(p.Data) == 0 {
		return enum.OpCode(0)
	}
	return enum.OpCode(p.Data[0])
}
