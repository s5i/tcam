package parser

import (
	"bytes"
	"fmt"

	"github.com/s5i/tcam/enum"
	"github.com/s5i/tcam/network"
)

type MoveCreature struct{}

func parseMoveCreature(p *network.Packet) (*MoveCreature, *network.Packet, error) {
	if p.OpCode() != enum.OpCodeMoveCreature {
		return nil, nil, fmt.Errorf("expected op code %s, got %s", enum.OpCodeMoveCreature, p.OpCode())
	}

	r := bytes.NewReader(p.Data)

	// Skip opcode (1 byte).
	if err := skip(r, 1); err != nil {
		return nil, nil, err
	}

	if _, err := mappedThing(r); err != nil {
		return nil, nil, err
	}

	if _, err := position(r); err != nil {
		return nil, nil, err
	}

	ret := &MoveCreature{}

	var next *network.Packet
	if cur := cur(r); cur != len(p.Data) {
		next = &network.Packet{
			Offset: p.Offset,
			Data:   p.Data[cur:],
		}
	}

	return ret, next, nil
}
