package parser

import (
	"bytes"
	"fmt"

	"github.com/s5i/tcam/enum"
	"github.com/s5i/tcam/network"
)

type MarkCreature struct{}

func parseMarkCreature(p *network.Packet) (*MarkCreature, *network.Packet, error) {
	if p.OpCode() != enum.OpCodeMarkCreature {
		return nil, nil, fmt.Errorf("expected op code %s, got %s", enum.OpCodeMarkCreature, p.OpCode())
	}

	r := bytes.NewReader(p.Data)

	if err := opcode(r); err != nil {
		return nil, nil, err
	}

	// ID (4 bytes), color (1 byte).
	if err := skip(r, 5); err != nil {
		return nil, nil, err
	}

	return &MarkCreature{}, p.Next(cur(r)), nil
}
