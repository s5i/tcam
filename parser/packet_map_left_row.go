package parser

import (
	"bytes"
	"fmt"

	"github.com/s5i/tcam/enum"
	"github.com/s5i/tcam/network"
)

type MapLeftRow struct{}

func parseMapLeftRow(p *network.Packet) (*MapLeftRow, *network.Packet, error) {
	if p.OpCode() != enum.OpCodeMapLeftRow {
		return nil, nil, fmt.Errorf("expected op code %s, got %s", enum.OpCodeMapLeftRow, p.OpCode())
	}

	r := bytes.NewReader(p.Data)

	if err := opcode(r); err != nil {
		return nil, nil, err
	}

	if err := mapDescription(r); err != nil {
		return nil, nil, err
	}

	return &MapLeftRow{}, p.Next(cur(r)), nil
}
