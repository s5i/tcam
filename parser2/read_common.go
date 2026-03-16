package parser2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/s5i/tcam/gamedata"
	"golang.org/x/text/encoding/charmap"
)

func ReadByte(r *bytes.Reader) (TByte, error) {
	var ret TByte
	if err := binary.Read(r, binary.LittleEndian, &ret); err != nil {
		return 0, err
	}
	return ret, nil
}

func ReadWord(r *bytes.Reader) (TWord, error) {
	var ret TWord
	if err := binary.Read(r, binary.LittleEndian, &ret); err != nil {
		return 0, err
	}
	return ret, nil
}

func ReadQuad(r *bytes.Reader) (TQuad, error) {
	var ret TQuad
	if err := binary.Read(r, binary.LittleEndian, &ret); err != nil {
		return 0, err
	}
	return ret, nil
}

func ReadString(r *bytes.Reader) (TString, error) {
	n, err := ReadWord(r)
	if err != nil {
		return "", err
	}

	raw := make([]byte, n)
	if _, err := io.ReadFull(r, raw); err != nil {
		return "", err
	}

	utf, err := decoder.Bytes(raw)
	if err != nil {
		return "", err
	}

	return TString(utf), nil
}

func ReadOutfit(r *bytes.Reader) (TOutfit, error) {
	var ret TOutfit

	id, err := ReadWord(r)
	if err != nil {
		return TOutfit{}, err
	}
	ret.ID = id

	if id == 0 {
		typ, err := ReadWord(r)
		if err != nil {
			return TOutfit{}, err
		}

		ret.Type = typ
	} else {
		colors, err := ReadQuad(r)
		if err != nil {
			return TOutfit{}, err
		}

		ret.Colors = colors
	}

	return ret, nil
}

// ReadItem reads an item from the stream.
func ReadItem(r *bytes.Reader) (TItem, error) {
	typeID, err := ReadWord(r)
	if err != nil {
		return TItem{}, err
	}

	ret := TItem{TypeID: typeID}

	flags := itemFlagsForID(typeID)
	if flags.IsFluidContainer || flags.IsFluid || flags.IsStackable {
		extra, err := ReadByte(r)
		if err != nil {
			return TItem{}, err
		}
		ret.ExtraByte = extra
		ret.HasExtra = true
	}

	return ret, nil
}

// ReadCreature reads a creature entry from the map stream.
func ReadCreature(r *bytes.Reader) (TCreature, error) {
	knownState, err := ReadWord(r)
	if err != nil {
		return TCreature{}, err
	}

	cr := TCreature{KnownState: knownState}
	needFullState := false

	switch knownState {
	case 99: // KNOWNCREATURE_UPTODATE
		id, err := ReadQuad(r)
		if err != nil {
			return TCreature{}, err
		}
		cr.ID = id

		dir, err := ReadByte(r)
		if err != nil {
			return TCreature{}, err
		}
		cr.Direction = dir

	case 97: // KNOWNCREATURE_FREE (unknown creature)
		needFullState = true

		removeID, err := ReadQuad(r)
		if err != nil {
			return TCreature{}, err
		}
		cr.RemoveID = removeID

		id, err := ReadQuad(r)
		if err != nil {
			return TCreature{}, err
		}
		cr.ID = id

		name, err := ReadString(r)
		if err != nil {
			return TCreature{}, err
		}
		cr.Name = name

	case 98: // KNOWNCREATURE_OUTDATED
		needFullState = true

		id, err := ReadQuad(r)
		if err != nil {
			return TCreature{}, err
		}
		cr.ID = id

	default:
		return TCreature{}, fmt.Errorf("unknown creature known-state: %d", knownState)
	}

	if !needFullState {
		return cr, nil
	}

	health, err := ReadByte(r)
	if err != nil {
		return TCreature{}, err
	}
	cr.Health = health

	dir, err := ReadByte(r)
	if err != nil {
		return TCreature{}, err
	}
	cr.Direction = dir

	outfit, err := ReadOutfit(r)
	if err != nil {
		return TCreature{}, err
	}
	cr.Outfit = outfit

	brightness, err := ReadByte(r)
	if err != nil {
		return TCreature{}, err
	}
	cr.LightBrightness = brightness

	color, err := ReadByte(r)
	if err != nil {
		return TCreature{}, err
	}
	cr.LightColor = color

	speed, err := ReadWord(r)
	if err != nil {
		return TCreature{}, err
	}
	cr.Speed = speed

	skull, err := ReadByte(r)
	if err != nil {
		return TCreature{}, err
	}
	cr.Skull = skull

	party, err := ReadByte(r)
	if err != nil {
		return TCreature{}, err
	}
	cr.Party = party

	return cr, nil
}

// ReadMapObject reads either a creature or an item from the map stream.
func ReadMapObject(r *bytes.Reader) (TMapObject, error) {
	typeID, err := ReadWord(r)
	if err != nil {
		return TMapObject{}, err
	}
	if err := unreadWord(r); err != nil {
		return TMapObject{}, err
	}

	switch {
	case typeID == 97 || typeID == 98 || typeID == 99:
		cr, err := ReadCreature(r)
		if err != nil {
			return TMapObject{}, err
		}
		return TMapObject{Creature: &cr}, nil

	default:
		item, err := ReadItem(r)
		if err != nil {
			return TMapObject{}, err
		}
		return TMapObject{Item: &item}, nil
	}
}

// ReadMapDescription reads a sequence of map tiles encoded with the skip-count
// format used by SendFullScreen, SendRow, SendFloors, and SendFieldData.
// It returns the list of map objects encountered (across all tiles).
func ReadMapDescription(r *bytes.Reader) ([]TMapObject, error) {
	var objects []TMapObject
	skip := -1

	for r.Len() > 0 {
		// Peek at the next word to check for skip markers.
		peekWord, err := ReadWord(r)
		if err != nil {
			return objects, nil // EOF is normal termination
		}

		if peekWord >= 0xFF00 {
			// This is a skip marker. The high byte is 0xFF, the low byte is the count.
			skipCount := int(peekWord & 0x00FF)
			skip += skipCount + 1
			_ = skip
			continue
		}

		// Not a skip marker — unread the 2 bytes and parse a map object.
		if err := unreadWord(r); err != nil {
			return nil, err
		}

		obj, err := ReadMapObject(r)
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}

	return objects, nil
}

// ReadPosition reads a 5-byte map position (x word, y word, z byte).
func ReadPosition(r *bytes.Reader) (TPosition, error) {
	x, err := ReadWord(r)
	if err != nil {
		return TPosition{}, err
	}
	y, err := ReadWord(r)
	if err != nil {
		return TPosition{}, err
	}
	z, err := ReadByte(r)
	if err != nil {
		return TPosition{}, err
	}
	return TPosition{X: x, Y: y, Z: z}, nil
}

// unreadWord seeks back 2 bytes.
func unreadWord(r *bytes.Reader) error {
	_, err := r.Seek(-2, 1) // io.SeekCurrent
	return err
}

type itemFlags struct {
	IsFluidContainer bool
	IsFluid          bool
	IsStackable      bool
}

func itemFlagsForID(id TWord) itemFlags {
	return itemFlags{
		IsFluidContainer: gamedata.IsFluidContainer[int(id)],
		IsFluid:          gamedata.IsFluid[int(id)],
		IsStackable:      gamedata.IsStackable[int(id)],
	}
}

var decoder = charmap.ISO8859_1.NewDecoder()
