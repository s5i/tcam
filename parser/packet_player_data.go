package parser

import (
	"bytes"
	"fmt"

	"github.com/s5i/tcam/enum"
	"github.com/s5i/tcam/network"
)

type PlayerData struct{}

func parsePlayerData(p *network.Packet) (*PlayerData, *network.Packet, error) {
	if p.OpCode() != enum.OpCodePlayerData {
		return nil, nil, fmt.Errorf("expected op code %s, got %s", enum.OpCodePlayerData, p.OpCode())
	}

	r := bytes.NewReader(p.Data)

	if err := opcode(r); err != nil {
		return nil, nil, err
	}

	// Health (2 bytes), max health (2 bytes), free cap (2 bytes), exp (4 bytes), level (2 bytes), level percent (1 byte), mana (2 bytes), max mana (2 bytes).
	// Mlvl (1 byte), mlvl pct (1 byte), soul (1 byte),

	if err := skip(r, 20); err != nil {
		return nil, nil, err
	}

	return &PlayerData{}, p.Next(cur(r)), nil
}
