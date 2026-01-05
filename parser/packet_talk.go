package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

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
	if _, err := r.Seek(5, io.SeekCurrent); err != nil {
		return nil, nil, err
	}

	var nameLen uint16
	if err := binary.Read(r, binary.LittleEndian, &nameLen); err != nil {
		return nil, nil, err
	}

	name := make([]byte, nameLen)
	if _, err := io.ReadFull(r, name); err != nil {
		Logger.Printf("Failed to read name: %v", err)
		return nil, nil, err
	}

	var mode enum.MessageMode
	if err := binary.Read(r, binary.LittleEndian, &mode); err != nil {
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
		if _, err := r.Seek(5, io.SeekCurrent); err != nil {
			return nil, nil, err
		}
	case
		enum.MessageModeMessageChannel,
		enum.MessageModeMessageChannelManagement,
		enum.MessageModeMessageChannelHighlight,
		enum.MessageModeMessageGamemasterChannel:
		// Channel ID, 2 bytes.
		if _, err := r.Seek(2, io.SeekCurrent); err != nil {
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

	var msgLen uint16
	if err := binary.Read(r, binary.LittleEndian, &msgLen); err != nil {
		return nil, nil, err
	}

	msg := make([]byte, msgLen)
	if _, err := io.ReadFull(r, msg); err != nil {
		Logger.Printf("Failed to read message: %v", err)
		return nil, nil, err
	}

	ret := &Talk{
		Name: string(name),
		Mode: mode,
		Msg:  string(msg),
	}

	cur, _ := r.Seek(0, io.SeekCurrent)

	var next *network.Packet
	if int(cur) != len(p.Data) {
		next = &network.Packet{
			Offset: p.Offset,
			Data:   p.Data[cur:],
		}
	}

	return ret, next, nil
}
