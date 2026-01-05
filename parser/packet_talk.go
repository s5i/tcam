package parser

import (
	"bytes"
	"fmt"

	"github.com/s5i/tcam/enum"
	"github.com/s5i/tcam/network"
)

type Talk struct {
	Name string
	Mode enum.MessageMode
	Msg  string
}

func parseTalk(p *network.Packet) (*Talk, *network.Packet, error) {
	if p.OpCode() != enum.OpCodeTalk {
		return nil, nil, fmt.Errorf("expected op code %s, got %s", enum.OpCodeTalk, p.OpCode())
	}

	r := bytes.NewReader(p.Data)

	// Skip opcode (1 byte), channel statement GUID (4 bytes).
	if err := skip(r, 5); err != nil {
		return nil, nil, err
	}

	name, err := str(r)
	if err != nil {
		return nil, nil, err
	}

	var mode enum.MessageMode
	if err := read(r, &mode); err != nil {
		return nil, nil, err
	}

	switch mode {
	case
		enum.MessageModeMessageSay,
		enum.MessageModeMessageWhisper,
		enum.MessageModeMessageYell,
		enum.MessageModeMessageMonsterSay,
		enum.MessageModeMessageMonsterYell,
		enum.MessageModeMessageBarkLow,
		enum.MessageModeMessageBarkLoud,
		enum.MessageModeMessageNpcFromStartBlock:
		// Position, 5 bytes.
		if err := skip(r, 5); err != nil {
			return nil, nil, err
		}
	case
		enum.MessageModeMessageChannel,
		enum.MessageModeMessageChannelManagement,
		enum.MessageModeMessageChannelHighlight,
		enum.MessageModeMessageGamemasterChannel:
		// Channel ID, 2 bytes.
		if err := skip(r, 2); err != nil {
			return nil, nil, err
		}
	case
		enum.MessageModeMessagePrivateFrom,
		enum.MessageModeMessageGamemasterBroadcast,
		enum.MessageModeMessageGamemasterPrivateFrom,
		enum.MessageModeMessageRVRAnswer,
		enum.MessageModeMessageRVRContinue:
	default:
		return nil, nil, fmt.Errorf("unknown message mode %s", mode)
	}

	msg, err := str(r)
	if err != nil {
		return nil, nil, err
	}

	ret := &Talk{
		Name: name,
		Mode: mode,
		Msg:  msg,
	}

	var next *network.Packet
	if cur := cur(r); cur != len(p.Data) {
		next = &network.Packet{
			Offset: p.Offset,
			Data:   p.Data[cur:],
		}
	}

	return ret, next, nil
}
