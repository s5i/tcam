package parser

import (
	"bytes"
	"fmt"

	"github.com/s5i/tcam/enum"
	"github.com/s5i/tcam/network"
)

type CreatureHealth struct{}

func parseCreatureHealth(p *network.Packet) (*CreatureHealth, *network.Packet, error) {
	if p.OpCode() != enum.OpCodeCreatureHealth {
		return nil, nil, fmt.Errorf("expected op code %s, got %s", enum.OpCodeCreatureHealth, p.OpCode())
	}

	r := bytes.NewReader(p.Data)

	if err := opcode(r); err != nil {
		return nil, nil, err
	}

	// ID (4 bytes), health pct (1 byte)
	if err := skip(r, 5); err != nil {
		return nil, nil, err
	}

	return &CreatureHealth{}, p.Next(cur(r)), nil
}
