package cam

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/s5i/tcam/dat"
	"github.com/s5i/tcam/data"
	"golang.org/x/text/encoding/charmap"
)

type message struct {
	r   io.ReadSeeker
	len int64
	dat *dat.File
}

func newMessage(buf []byte, dat *dat.File) *message {
	return &message{
		r:   bytes.NewReader(buf),
		len: int64(len(buf)),
		dat: dat,
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
	utf, err := decoder.Bytes(buf)
	if err != nil {
		return err
	}

	*ret = string(utf)
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
	stats        *ParseStats
	playerPos    data.Location
	playerID     uint32
	playerName   string
	serverName   string
	lastVisit    time.Time
	seenMessage  bool
}

const lastVisitPrefix = "Your last visit in "

func parseLastVisitMessage(text string) (serverName string, lastVisit time.Time, ok bool) {
	if !strings.HasPrefix(text, lastVisitPrefix) {
		return "", time.Time{}, false
	}
	rest := text[len(lastVisitPrefix):]
	colon := strings.LastIndex(rest, ": ")
	if colon < 0 {
		return "", time.Time{}, false
	}
	ts := strings.TrimSuffix(rest[colon+2:], ".")
	lastVisit, err := time.Parse("02. Jan 2006 15:04:05 MST", ts)
	if err != nil {
		return "", time.Time{}, false
	}
	return rest[:colon], lastVisit, true
}

func (s *parseState) resolvePlayerName(tiles []data.Tile) {
	if s.playerName != "" || s.playerID == 0 {
		return
	}

	for _, tile := range tiles {
		for _, t := range tile.Things {
			if t.HasCreature && t.Creature.ID == s.playerID {
				s.playerName = t.Creature.Name
				return
			}
		}
	}
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
		if !ignore {
			tiles = append(tiles, t...)
		}
	}
	return tiles, nil
}

func (m *message) parseFloorDescription(ignore bool, x, y, z, width, height, offset int, skip *int) ([]data.Tile, error) {
	var tiles []data.Tile
	for nx := range width {
		for ny := range height {
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
					if !ignore {
						tiles = append(tiles, tile)
					}

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
		if !ignore {
			tile.Things = append(tile.Things, thing)
		}
	}
	return tile, nil
}

func (m *message) getThing(ignore bool) (data.Thing, error) {
	thingID, err := m.getU16()
	if err != nil {
		return data.Thing{}, err
	}

	if thingID == 0x0061 || thingID == 0x0062 {
		c := data.Creature{}
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
		return data.Thing{HasCreature: true, Creature: c}, nil
	}

	if thingID == 0x0063 {
		// Creature turn.
		c := data.Creature{}
		c.ID, err = m.getU32()
		if err != nil {
			return data.Thing{}, err
		}
		dir, err := m.getByte()
		if err != nil {
			return data.Thing{}, err
		}
		c.Direction = data.Direction(dir)
		return data.Thing{HasCreature: true, Creature: c}, nil
	}

	item, err := getItem(m, thingID)
	if err != nil {
		return data.Thing{}, err
	}
	return data.Thing{HasItem: true, Item: item}, nil
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
	if m.dat.IsStackable(int(itemID)) {
		item.Count, err = m.getByte()
		if err != nil {
			return data.Item{}, err
		}
	} else if m.dat.IsFluid(int(itemID)) || m.dat.IsFluidContainer(int(itemID)) {
		item.SubType, err = m.getByte()
		if err != nil {
			return data.Item{}, err
		}
	}
	return item, nil
}

var decoder = charmap.Windows1252.NewDecoder()
