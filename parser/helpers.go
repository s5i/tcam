package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/s5i/tcam/enum"
	"github.com/s5i/tcam/gamedata"
)

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

func opcode(r *bytes.Reader) error {
	return skip(r, 1)
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

func position(r *bytes.Reader) error {
	// Position, 5 bytes.
	if err := skip(r, 5); err != nil {
		return err
	}
	return nil
}

func mappedThing(r *bytes.Reader) error {
	var x uint16
	if err := read(r, &x); err != nil {
		return err
	}

	if x == 0xFFFF {
		// Creature ID, 4 bytes.
		if err := skip(r, 4); err != nil {
			return err
		}

		return nil
	}

	// Y (2 bytes), Z (1 byte), stack position (1 byte).
	if err := skip(r, 4); err != nil {
		return err
	}

	return nil
}

func thing(r *bytes.Reader) error {
	var id enum.Item
	if err := read(r, &id); err != nil {
		return err
	}

	switch id {
	case enum.ItemInvalid:
		return fmt.Errorf("invalid item id: %s", enum.ItemInvalid)
	case enum.ItemUnknownCreature, enum.ItemOutdatedCreature, enum.ItemCreature:
		return creature(r, id)
	default:
		return item(r, id)
	}
}

func creature(r *bytes.Reader, t enum.Item) error {
	switch t {
	case enum.ItemInvalid:
		// Type (2 bytes).
		if err := read(r, &t); err != nil {
			return err
		}
		fallthrough

	case enum.ItemOutdatedCreature:
		// ID (4 bytes).
		if err := skip(r, 4); err != nil {
			return err
		}

		// Health percent (1 byte), direction (1 byte).
		if err := skip(r, 2); err != nil {
			return err
		}

		if err := outfit(r); err != nil {
			return err
		}

		// Light (2 bytes), speed (2 bytes), skull (1 byte), shield (1 byte).
		if err := skip(r, 6); err != nil {
			return err
		}

	case enum.ItemUnknownCreature:
		// RemoveID (4 bytes), ID (4 bytes).
		if err := skip(r, 8); err != nil {
			return err
		}

		// Name.
		if _, err := str(r); err != nil {
			return err
		}

		// Health percent (1 byte), direction (1 byte).
		if err := skip(r, 2); err != nil {
			return err
		}

		if err := outfit(r); err != nil {
			return err
		}

		// Light (2 bytes), speed (2 bytes), skull (1 byte), shield (1 byte).
		if err := skip(r, 6); err != nil {
			return err
		}

	case enum.ItemCreature:
		// ID (4 bytes), direction (1 byte).
		if err := skip(r, 5); err != nil {
			return err
		}
	}

	return nil
}

func outfit(r *bytes.Reader) error {
	var t uint16
	if err := read(r, &t); err != nil {
		return err
	}

	// LookTypeEx (2 bytes).
	if t == 0 {
		return skip(r, 2)
	}

	// Head (1 byte), body (1 byte), legs (1 byte), feet (1 byte).
	return skip(r, 4)
}

func item(r *bytes.Reader, t enum.Item) error {
	if t == enum.ItemInvalid {
		// Type (2 bytes).
		if err := read(r, &t); err != nil {
			return err
		}
	}

	attr := gamedata.Attrs[gamedata.DATKey{Category: enum.DatCategoryItem, ID: int(t)}]

	if forcedSkip, ok := itemSkipOverrides[t]; ok {
		if forcedSkip {
			return skip(r, 1)
		}
		return nil
	}

	switch {
	case
		attr.Present[enum.DatAttributeStackable],
		attr.Present[enum.DatAttributeChargeable],
		attr.Present[enum.DatAttributeFluidContainer],
		attr.Present[enum.DatAttributeSplash]:

		if err := skip(r, 1); err != nil {
			return err
		}
	}
	return nil
}

var itemSkipOverrides = map[enum.Item]bool{
	1644: false,
	2887: true,
	2888: true,
	3031: true,
	3277: true,
	3577: true,
	3582: true,
	3606: true,
	3725: true,
}
