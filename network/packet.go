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
	return fmt.Sprintf("%s - %s - %s - %03d", p.Offset(), p.Time(), p.OpCode(), p.Data)
}

func (p Packet) Offset() string {
	pos := p.GlobalOffset + p.LocalOffset
	return fmt.Sprintf("[%8x +%2d]", pos-pos%16, pos%16)
}

func (p Packet) Time() string {
	return fmt.Sprintf("%v", p.TimeOffset.Truncate(time.Second))
}

func (p Packet) OpCode() enum.OpCode {
	if len(p.Data) == 0 {
		return enum.OpCode(0)
	}
	return enum.OpCode(p.Data[0])
}
