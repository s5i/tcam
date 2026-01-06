package gamedata

import (
	"context"
	"encoding/binary"
	"io"
	"maps"
	"os"
	"slices"

	"github.com/s5i/tcam/enum"
)

type DATKey struct {
	Category enum.DatCategory
	ID       int
}

type DATAttributes struct {
	Present map[enum.DatAttribute]bool
}

var Attrs map[DATKey]DATAttributes

func ReadFile(ctx context.Context, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var sig uint32
	if err := binary.Read(f, binary.LittleEndian, &sig); err != nil {
		return err
	}

	ret := map[DATKey]DATAttributes{}
	count := map[enum.DatCategory]int{}

	if err := func() error {
		for cat := enum.DatCategory(0); cat < enum.DatCategoryLast; cat++ {
			var c uint16
			if err := binary.Read(f, binary.LittleEndian, &c); err != nil {
				return err
			}
			count[cat] = int(c) + 1
			Logger.Printf("%v => %v entries", cat, count[cat])
		}

		for cat := enum.DatCategory(0); cat < enum.DatCategoryLast; cat++ {
			firstID := 1
			if cat == enum.DatCategoryItem {
				firstID = 100
			}

			for id := firstID; id < count[cat]; id++ {
				key := DATKey{Category: cat, ID: id}
				attrs := DATAttributes{
					Present: map[enum.DatAttribute]bool{},
				}

				for attr := enum.DatAttribute(0); attr < enum.DatAttributeLast; attr++ {
					if err := binary.Read(f, binary.LittleEndian, &attr); err != nil {
						return err
					}

					if attr == enum.DatAttributeLast {
						break
					}

					attrs.Present[attr] = true

					switch attr {
					case enum.DatAttributeDisplacement:
						// x, y displacement (2 bytes each)
						if _, err := f.Seek(4, io.SeekCurrent); err != nil {
							return err
						}
					case enum.DatAttributeLight:
						// intensity (2 bytes), color (2 bytes)
						if _, err := f.Seek(4, io.SeekCurrent); err != nil {
							return err
						}
					case enum.DatAttributeMarket:
						// category (2 bytes), trade as (2 bytes), show as (2 bytes)
						if _, err := f.Seek(6, io.SeekCurrent); err != nil {
							return err
						}

						// name length (2 bytes), name (variable)
						var marketNameLength uint16
						if err := binary.Read(f, binary.LittleEndian, &marketNameLength); err != nil {
							return err
						}
						if _, err := f.Seek(int64(marketNameLength), io.SeekCurrent); err != nil {
							return err
						}

						// restrict vocation (2 bytes), required level (2 bytes)
						if _, err := f.Seek(4, io.SeekCurrent); err != nil {
							return err
						}

					case enum.DatAttributeElevation:
						// elevation (2 bytes)
						if _, err := f.Seek(2, io.SeekCurrent); err != nil {
							return err
						}

					case enum.DatAttributeLensHelp:
						// whatever this is (2 bytes)
						if _, err := f.Seek(2, io.SeekCurrent); err != nil {
							return err
						}

					case enum.DatAttributeBones:
						// x, y (2 bytes each) * [north, south, east, west]
						if _, err := f.Seek(16, io.SeekCurrent); err != nil {
							return err
						}
					}
				}

				var w, h uint8
				if err := binary.Read(f, binary.LittleEndian, &w); err != nil {
					return err
				}
				if err := binary.Read(f, binary.LittleEndian, &h); err != nil {
					return err
				}
				if w > 1 || h > 1 {
					// "real size" (1 byte)
					if _, err := f.Seek(1, io.SeekCurrent); err != nil {
						return err
					}
				}

				var layers, numPatternX, numPatternY, numPatternZ, animationPhases uint8
				if err := binary.Read(f, binary.LittleEndian, &layers); err != nil {
					return err
				}
				if err := binary.Read(f, binary.LittleEndian, &numPatternX); err != nil {
					return err
				}
				if err := binary.Read(f, binary.LittleEndian, &numPatternY); err != nil {
					return err
				}
				if err := binary.Read(f, binary.LittleEndian, &numPatternZ); err != nil {
					return err
				}
				if err := binary.Read(f, binary.LittleEndian, &animationPhases); err != nil {
					return err
				}

				numSprites := w * h * layers * numPatternX * numPatternY * numPatternZ * animationPhases
				if _, err := f.Seek(2*int64(numSprites), io.SeekCurrent); err != nil {
					return err
				}

				ret[key] = attrs
				Logger.Printf("%v => %v", key, slices.Sorted(maps.Keys(attrs.Present)))
			}
		}
		return nil
	}(); err != nil && err != io.EOF {
		return err
	}

	Attrs = ret

	return nil
}
