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

	if err := opcode(r); err != nil {
		return nil, nil, err
	}

	if err := mappedThing(r); err != nil {
		return nil, nil, err
	}

	if err := position(r); err != nil {
		return nil, nil, err
	}

	ret := &MoveCreature{}
	return ret, p.Next(cur(r)), nil
}
