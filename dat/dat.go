package dat

import (
	"fmt"
	"io"
)

// Properties holds metadata flags parsed from a Tibia .dat item entry.
type Properties struct {
	Ground         bool
	GroundBorder   bool
	OnBottom       bool
	OnTop          bool
	Container      bool
	Stackable      bool
	ForceUse       bool
	MultiUse       bool
	FluidContainer bool
	Fluid          bool
	Unpassable     bool
	Unmoveable     bool
	BlockMissile   bool
	BlockPathfind  bool
	Pickupable     bool
	Hangable       bool
	Usable         bool
}

// File holds item metadata parsed from a Tibia client .dat file.
type File struct {
	Signature    uint32
	ItemCount    uint16
	OutfitCount  uint16
	EffectCount  uint16
	MissileCount uint16

	properties []Properties
}

// Read parses item metadata from a Tibia .dat file stream.
//
// Only the item section is consumed; outfit, effect, and missile sections are
// not loaded. Item IDs start at 100 and run through ItemCount inclusive.
func Read(r io.Reader) (*File, error) {
	br := &binReader{r: r}

	signature, itemCount, outfitCount, effectCount, missileCount, err := readHeader(br)
	if err != nil {
		return nil, fmt.Errorf("dat: read header: %w", err)
	}
	if itemCount < firstItemID {
		return nil, fmt.Errorf("dat: item count %d is below first item id %d", itemCount, firstItemID)
	}

	version := formatForSignature(signature)
	patternZFixed := version == formatV1 || version == formatV2

	f := &File{
		Signature:    signature,
		ItemCount:    itemCount,
		OutfitCount:  outfitCount,
		EffectCount:  effectCount,
		MissileCount: missileCount,
		properties:   make([]Properties, int(itemCount)+1),
	}

	for id := uint16(firstItemID); id <= itemCount; id++ {
		props, err := readProperties(br, version, id)
		if err != nil {
			return nil, fmt.Errorf("dat: item %d properties: %w", id, err)
		}
		if err := readTexturePatterns(br, patternZFixed); err != nil {
			return nil, fmt.Errorf("dat: item %d sprites: %w", id, err)
		}
		f.properties[id] = props
	}

	return f, nil
}

func (f *File) propertiesFor(id int) (Properties, bool) {
	if id < firstItemID || id > int(f.ItemCount) {
		return Properties{}, false
	}
	return f.properties[id], true
}

// Properties returns parsed metadata for an item ID, or false if the ID is out of range.
func (f *File) Properties(id int) (Properties, bool) {
	return f.propertiesFor(id)
}

func (f *File) IsStackable(id int) bool {
	p, ok := f.propertiesFor(id)
	return ok && p.Stackable
}

func (f *File) IsFluidContainer(id int) bool {
	p, ok := f.propertiesFor(id)
	return ok && p.FluidContainer
}

func (f *File) IsFluid(id int) bool {
	p, ok := f.propertiesFor(id)
	return ok && p.Fluid
}

func (f *File) IsContainer(id int) bool {
	p, ok := f.propertiesFor(id)
	return ok && p.Container
}

func (f *File) IsGround(id int) bool {
	p, ok := f.propertiesFor(id)
	return ok && p.Ground
}

func (f *File) IsPickupable(id int) bool {
	p, ok := f.propertiesFor(id)
	return ok && p.Pickupable
}

func (f *File) IsUnmoveable(id int) bool {
	p, ok := f.propertiesFor(id)
	return ok && p.Unmoveable
}
