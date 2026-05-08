package cam

import (
	"github.com/s5i/tcam/data"
)

func getMapDescription(m *message, x, y, z, width, height int) ([]data.Tile, error) {
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
		t, err := parseFloorDescription(m, x, y, nz, width, height, z-nz, &skip)
		if err != nil {
			return nil, err
		}
		tiles = append(tiles, t...)
	}
	return tiles, nil
}

func parseFloorDescription(m *message, x, y, z, width, height, offset int, skip *int) ([]data.Tile, error) {
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
					tile, err := parseTileDescription(m, loc)
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

func parseTileDescription(m *message, loc data.Location) (data.Tile, error) {
	tile := data.Tile{Location: loc}
	for {
		peek, err := m.peekU16()
		if err != nil {
			return data.Tile{}, err
		}
		if peek >= 0xFF00 {
			break
		}
		thing, err := getThing(m)
		if err != nil {
			return data.Tile{}, err
		}
		tile.Things = append(tile.Things, thing)
	}
	return tile, nil
}

func getThing(m *message) (data.Thing, error) {
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
			c.Name, err = m.getString()
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
