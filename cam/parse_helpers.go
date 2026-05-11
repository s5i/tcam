package cam

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/s5i/tcam/data"
)

const emptyString = ""

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

func (m *message) getString(ret *string, ignore bool) error {
	length, err := m.getU16()
	if err != nil {
		return err
	}
	if ignore {
		if _, err := m.r.Seek(int64(length), io.SeekCurrent); err != nil {
			return err
		}
		return nil
	}
	buf := make([]byte, length)
	if _, err := io.ReadFull(m.r, buf); err != nil {
		return err
	}
	*ret = string(buf)
	return nil
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

type parseState struct {
	playerPos data.Location
	tiles     map[tileKey][]data.Thing
}

func (m *message) getMapDescription(ignore bool, x, y, z, width, height int) ([]data.Tile, error) {
	startz, endz, zstep := 0, 0, 0
	if z > 7 {
		startz = z - 2
		endz = min(15, z+2)
		zstep = 1
	} else {
		startz = 7
		endz = 0
		zstep = -1
	}

	var tiles []data.Tile
	skip := 0
	for nz := startz; nz != endz+zstep; nz += zstep {
		t, err := m.parseFloorDescription(ignore, x, y, nz, width, height, z-nz, &skip)
		if err != nil {
			return nil, err
		}
		tiles = append(tiles, t...)
	}
	return tiles, nil
}

func (m *message) parseFloorDescription(ignore bool, x, y, z, width, height, offset int, skip *int) ([]data.Tile, error) {
	var tiles []data.Tile
	for nx := 0; nx < width; nx++ {
		for ny := 0; ny < height; ny++ {
			if *skip == 0 {
				tileOpt, err := m.peekU16()
				if err != nil {
					return nil, err
				}
				if tileOpt >= 0xFF00 {
					v, err := m.getU16()
					if err != nil {
						return nil, err
					}
					*skip = int(v & 0xFF)
				} else {
					loc := data.Location{X: x + nx + offset, Y: y + ny + offset, Z: z}
					tile, err := m.parseTileDescription(ignore, loc)
					if err != nil {
						return nil, err
					}
					tiles = append(tiles, tile)
					v, err := m.getU16()
					if err != nil {
						return nil, err
					}
					*skip = int(v & 0xFF)
				}
			} else {
				*skip--
			}
		}
	}
	return tiles, nil
}

func (m *message) parseTileDescription(ignore bool, loc data.Location) (data.Tile, error) {
	tile := data.Tile{Location: loc}
	for {
		peek, err := m.peekU16()
		if err != nil {
			return data.Tile{}, err
		}
		if peek >= 0xFF00 {
			break
		}
		thing, err := m.getThing(ignore)
		if err != nil {
			return data.Tile{}, err
		}
		tile.Things = append(tile.Things, thing)
	}
	return tile, nil
}

func (m *message) getThing(ignore bool) (data.Thing, error) {
	thingID, err := m.getU16()
	if err != nil {
		return data.Thing{}, err
	}

	if thingID == 0x0061 || thingID == 0x0062 {
		c := &data.Creature{}
		if thingID == 0x0062 {
			// Known creature.
			c.ID, err = m.getU32()
			if err != nil {
				return data.Thing{}, err
			}
			c.Health, err = m.getByte()
			if err != nil {
				return data.Thing{}, err
			}
		} else {
			// Unknown creature (0x0061).
			c.RemovedID, err = m.getU32()
			if err != nil {
				return data.Thing{}, err
			}
			c.ID, err = m.getU32()
			if err != nil {
				return data.Thing{}, err
			}
			err = m.getString(&c.Name, ignore)
			if err != nil {
				return data.Thing{}, err
			}
			c.Health, err = m.getByte()
			if err != nil {
				return data.Thing{}, err
			}
		}

		dir, err := m.getByte()
		if err != nil {
			return data.Thing{}, err
		}
		c.Direction = data.Direction(dir)

		c.Outfit, err = m.getOutfit()
		if err != nil {
			return data.Thing{}, err
		}
		c.LightLevel, err = m.getByte()
		if err != nil {
			return data.Thing{}, err
		}
		c.LightColor, err = m.getByte()
		if err != nil {
			return data.Thing{}, err
		}
		c.Speed, err = m.getU16()
		if err != nil {
			return data.Thing{}, err
		}
		c.Skull, err = m.getByte()
		if err != nil {
			return data.Thing{}, err
		}
		c.Shield, err = m.getByte()
		if err != nil {
			return data.Thing{}, err
		}
		return data.Thing{Creature: c}, nil
	}

	if thingID == 0x0063 {
		// Creature turn.
		c := &data.Creature{}
		c.ID, err = m.getU32()
		if err != nil {
			return data.Thing{}, err
		}
		dir, err := m.getByte()
		if err != nil {
			return data.Thing{}, err
		}
		c.Direction = data.Direction(dir)
		return data.Thing{Creature: c}, nil
	}

	item, err := getItem(m, thingID)
	if err != nil {
		return data.Thing{}, err
	}
	return data.Thing{Item: &item}, nil
}

func getItem(m *message, itemID uint16) (data.Item, error) {
	var err error
	if itemID == 0xFFFF {
		itemID, err = m.getU16()
		if err != nil {
			return data.Item{}, err
		}
	}

	item := data.Item{ID: itemID}
	if data.IsStackable[int(itemID)] {
		item.Count, err = m.getByte()
		if err != nil {
			return data.Item{}, err
		}
	} else if data.IsFluid[int(itemID)] || data.IsFluidContainer[int(itemID)] {
		item.SubType, err = m.getByte()
		if err != nil {
			return data.Item{}, err
		}
	}
	return item, nil
}

type tileKey struct {
	x, y, z int
}

func (s *parseState) updateTiles(tiles []data.Tile) {
	for _, t := range tiles {
		k := tileKey{t.Location.X, t.Location.Y, t.Location.Z}
		s.tiles[k] = append([]data.Thing(nil), t.Things...)
	}
}

func (s *parseState) getThing(loc data.Location, stack int) *data.Thing {
	k := tileKey{loc.X, loc.Y, loc.Z}
	things := s.tiles[k]
	if stack < 0 || stack >= len(things) {
		return nil
	}
	return &things[stack]
}

func (s *parseState) addThing(loc data.Location, thing data.Thing) {
	k := tileKey{loc.X, loc.Y, loc.Z}
	s.tiles[k] = append(s.tiles[k], thing)
}

func (s *parseState) removeThing(loc data.Location, stack int) {
	k := tileKey{loc.X, loc.Y, loc.Z}
	things := s.tiles[k]
	if stack < 0 || stack >= len(things) {
		return
	}
	s.tiles[k] = append(things[:stack], things[stack+1:]...)
}

func (s *parseState) replaceThing(loc data.Location, stack int, thing data.Thing) {
	k := tileKey{loc.X, loc.Y, loc.Z}
	things := s.tiles[k]
	if stack < 0 || stack >= len(things) {
		return
	}
	things[stack] = thing
}
