package data

// Location represents a position in the game world.
type Location struct {
	X, Y, Z int
}

// IsCreature returns true if this location encodes a creature reference (X == 65535).
func (l Location) IsCreature() bool {
	return l.X == 65535
}

// CreatureID extracts a creature ID from a creature-reference location.
func (l Location) CreatureID(stack byte) uint32 {
	return uint32(l.Y) | uint32(l.Z)<<16 | uint32(stack)<<24
}

// Direction represents a facing direction.
type Direction byte

const (
	North Direction = 0
	East  Direction = 1
	South Direction = 2
	West  Direction = 3
	NE    Direction = 4
	NW    Direction = 5
	SE    Direction = 6
	SW    Direction = 7
)

// Outfit represents a creature's visual appearance.
type Outfit struct {
	LookType               uint16
	LookItem               uint16 // only when LookType == 0
	Head, Body, Legs, Feet byte   // only when LookType != 0
}

// Item represents an in-game item.
type Item struct {
	ID      uint16
	Count   byte // for stackable items
	SubType byte // for splash/fluid container items
}

// Creature represents a creature (player, NPC, or monster).
type Creature struct {
	ID         uint32
	RemovedID  uint32 // for "unknown creature" (0x0061), the old creature ID being replaced
	Name       string
	Health     byte
	Direction  Direction
	Outfit     Outfit
	LightLevel byte
	LightColor byte
	Speed      uint16
	Skull      byte
	Shield     byte
}

// Thing represents either a Creature or an Item on a tile.
type Thing struct {
	Creature    Creature
	HasCreature bool

	Item    Item
	HasItem bool
}

// Tile represents a map tile at a specific location containing things.
type Tile struct {
	Location Location
	Things   []Thing
}
