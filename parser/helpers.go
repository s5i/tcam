package parser

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Position struct{}
type MappedThing struct{}

func read(r io.Reader, data any) error {
	return binary.Read(r, binary.LittleEndian, data)
}

func cur(r *bytes.Reader) int {
	cur, _ := r.Seek(0, io.SeekCurrent)
	return int(cur)
}

func skip(r *bytes.Reader, n int) error {
	if _, err := r.Seek(int64(n), io.SeekCurrent); err != nil {
		return err
	}
	return nil
}

func str(r *bytes.Reader) (string, error) {
	var length uint16
	if err := read(r, &length); err != nil {
		return "", err
	}

	content := make([]byte, length)
	if _, err := io.ReadFull(r, content); err != nil {
		return "", err
	}

	return string(content), nil
}

func position(r *bytes.Reader) (*Position, error) {
	// Position, 5 bytes.
	if err := skip(r, 5); err != nil {
		return nil, err
	}
	return &Position{}, nil
}

func mappedThing(r *bytes.Reader) (*MappedThing, error) {
	var x uint16
	if err := read(r, &x); err != nil {
		return nil, err
	}

	if x == 0xFFFF {
		// Creature ID, 4 bytes.
		if err := skip(r, 4); err != nil {
			return nil, err
		}

		return &MappedThing{}, nil
	}

	// Y (2 bytes), Z (1 byte), stack position (1 byte).
	if err := skip(r, 4); err != nil {
		return nil, err
	}

	return &MappedThing{}, nil
}
