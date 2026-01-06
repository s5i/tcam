package parser

import (
	"bytes"
	"fmt"
	"time"

	"github.com/s5i/tcam/enum"
	"github.com/s5i/tcam/network"
)

type Talk struct {
	Name       string
	Mode       enum.MessageMode
	Msg        string
	TimeOffset time.Duration
}

func parseTalk(p *network.Packet, checkIntegrity bool) (*Talk, *network.Packet, error) {
	if p.OpCode() != enum.OpCodeTalk {
		return nil, nil, fmt.Errorf("expected op code %s, got %s", enum.OpCodeTalk, p.OpCode())
	}

	r := bytes.NewReader(p.Data)

	if err := opcode(r); err != nil {
		return nil, nil, err
	}

	// Channel statement GUID (4 bytes).
	if err := skip(r, 4); err != nil {
		return nil, nil, err
	}

	name, err := str(r)
	if err != nil {
		return nil, nil, err
	}

	if checkIntegrity {
		for _, c := range []byte(name) {
			if c < 32 || c > 126 {
				return nil, nil, fmt.Errorf("invalid name %q", name)
			}
		}
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

		if err := position(r); err != nil {
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
		enum.MessageModeMessageRVRContinue,
		enum.MessageMode14: // ??

	default:
		if checkIntegrity {
			return nil, nil, fmt.Errorf("unknown message mode %s", mode)
		}
	}

	msg, err := str(r)
	if err != nil {
		return nil, nil, err
	}

	if checkIntegrity {
		for _, c := range []byte(msg) {
			if c < 32 || c > 126 {
				return nil, nil, fmt.Errorf("invalid message %q", msg)
			}
		}
	}

	ret := &Talk{
		Name:       name,
		Mode:       mode,
		Msg:        msg,
		TimeOffset: p.TimeOffset,
	}

	return ret, p.Next(cur(r)), nil
}
