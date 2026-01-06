package parser

import (
	"bytes"
	"fmt"

	"github.com/s5i/tcam/enum"
	"github.com/s5i/tcam/network"
)

type ChangeOnMap struct{}

func parseChangeOnMap(p *network.Packet) (*ChangeOnMap, *network.Packet, error) {
	if p.OpCode() != enum.OpCodeChangeOnMap {
		return nil, nil, fmt.Errorf("expected op code %s, got %s", enum.OpCodeChangeOnMap, p.OpCode())
	}

	r := bytes.NewReader(p.Data)

	if err := opcode(r); err != nil {
		return nil, nil, err
	}

	if err := mappedThing(r); err != nil {
		return nil, nil, err
	}

	if err := thing(r); err != nil {
		return nil, nil, err
	}

	ret := &ChangeOnMap{}
	return ret, p.Next(cur(r)), nil
}
