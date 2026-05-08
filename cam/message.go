package cam

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/s5i/tcam/data"
)

type message struct {
	r   io.ReadSeeker
	len int64
}

func newMessage(buf []byte) *message {
	return &message{
		r:   bytes.NewReader(buf),
		len: int64(len(buf)),
	}
}

func (m *message) getByte() (byte, error) {
	var b byte
	err := binary.Read(m.r, binary.LittleEndian, &b)
	return b, err
}

func (m *message) getU16() (uint16, error) {
	var v uint16
	err := binary.Read(m.r, binary.LittleEndian, &v)
	return v, err
}

func (m *message) getU32() (uint32, error) {
	var v uint32
	err := binary.Read(m.r, binary.LittleEndian, &v)
	return v, err
}

func (m *message) getString() (string, error) {
	length, err := m.getU16()
	if err != nil {
		return "", err
	}
	buf := make([]byte, length)
	if _, err := io.ReadFull(m.r, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func (m *message) getLocation() (data.Location, error) {
	x, err := m.getU16()
	if err != nil {
		return data.Location{}, err
	}
	y, err := m.getU16()
	if err != nil {
		return data.Location{}, err
	}
	z, err := m.getByte()
	if err != nil {
		return data.Location{}, err
	}
	return data.Location{X: int(x), Y: int(y), Z: int(z)}, nil
}

func (m *message) getOutfit() (data.Outfit, error) {
	lookType, err := m.getU16()
	if err != nil {
		return data.Outfit{}, err
	}
	if lookType != 0 {
		head, err := m.getByte()
		if err != nil {
			return data.Outfit{}, err
		}
		body, err := m.getByte()
		if err != nil {
			return data.Outfit{}, err
		}
		legs, err := m.getByte()
		if err != nil {
			return data.Outfit{}, err
		}
		feet, err := m.getByte()
		if err != nil {
			return data.Outfit{}, err
		}
		return data.Outfit{LookType: lookType, Head: head, Body: body, Legs: legs, Feet: feet}, nil
	}
	lookItem, err := m.getU16()
	if err != nil {
		return data.Outfit{}, err
	}
	return data.Outfit{LookType: lookType, LookItem: lookItem}, nil
}

func (m *message) peekU16() (uint16, error) {
	v, err := m.getU16()
	if err != nil {
		return 0, err
	}
	if _, err := m.r.Seek(-2, io.SeekCurrent); err != nil {
		return 0, fmt.Errorf("seek back after peek: %w", err)
	}
	return v, nil
}

func (m *message) remaining() int {
	cur, err := m.r.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0
	}
	return int(m.len - cur)
}
