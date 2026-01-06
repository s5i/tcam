package parser

import (
	"bytes"
	"fmt"

	"github.com/s5i/tcam/enum"
	"github.com/s5i/tcam/network"
)

type LoginOrPendingState struct{}

func parseLoginOrPendingState(p *network.Packet) (*LoginOrPendingState, *network.Packet, error) {
	if p.OpCode() != enum.OpCodeLoginOrPendingState {
		return nil, nil, fmt.Errorf("expected op code %s, got %s", enum.OpCodeLoginOrPendingState, p.OpCode())
	}

	r := bytes.NewReader(p.Data)

	if err := opcode(r); err != nil {
		return nil, nil, err
	}

	// Player ID (4 bytes), Server Beat (2 bytes), Can Report Bugs (1 byte)
	if err := skip(r, 7); err != nil {
		return nil, nil, err
	}

	ret := &LoginOrPendingState{}
	return ret, p.Next(cur(r)), nil
}
