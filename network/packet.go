package network

import (
	"fmt"
	"time"

	"github.com/s5i/tcam/enum"
)

// Packet represents a game packet with a time offset.
type Packet struct {
	GlobalOffset int
	LocalOffset  int
	TimeOffset   time.Duration
	Data         []byte
}

func (p Packet) OpCode() enum.OpCode {
	if len(p.Data) == 0 {
		return enum.OpCode(0)
	}
	return enum.OpCode(p.Data[0])
}

func (p *Packet) Next(offset int) *Packet {
	if offset == len(p.Data) {
		return nil
	}

	return &Packet{
		GlobalOffset: p.GlobalOffset,
		LocalOffset:  p.LocalOffset + offset,
		TimeOffset:   p.TimeOffset,
		Data:         p.Data[offset:],
	}
}

func (p Packet) String() string {
	pos := p.GlobalOffset + p.LocalOffset
	return fmt.Sprintf("[%8x+%x] - %s - %s - %v", pos/16, pos%16, p.TimeOffset.Truncate(time.Second), p.OpCode(), p.Data)
}
