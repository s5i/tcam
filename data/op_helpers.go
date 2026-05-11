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

var (
	IsFluidContainer = map[int]bool{}
	IsFluid          = map[int]bool{}
	IsStackable      = map[int]bool{}

	sFluidContainers = []int{
		2524, 2873, 2874, 2875, 2876, 2877, 2879, 2880, 2881, 2882,
		2885, 2893, 2901, 2902, 2903, 2904, 3465, 3477, 3478, 3479,
		3480,
	}
	sFluids = []int{
		2886, 2887, 2888, 2889, 2890, 2891, 2895, 2896, 2897, 2898,
		2899, 2900,
	}
	sStackable = []int{
		1781, 2784, 2992, 3026, 3027, 3028, 3029, 3030, 3031, 3032,
		3033, 3034, 3035, 3040, 3042, 3043, 3114, 3145, 3146, 3207,
		3215, 3250, 3277, 3287, 3298, 3445, 3446, 3447, 3448, 3449,
		3450, 3492, 3533, 3534, 3547, 3548, 3560, 3577, 3578, 3579,
		3580, 3581, 3582, 3583, 3584, 3585, 3586, 3587, 3588, 3589,
		3590, 3591, 3595, 3596, 3597, 3598, 3599, 3600, 3601, 3602,
		3603, 3604, 3605, 3606, 3721, 3722, 3723, 3724, 3725, 3726,
		3727, 3728, 3729, 3730, 3731, 3732, 3734, 3735, 3736, 3737,
		3738, 3739, 3740, 3741, 5021,
	}
)

func init() {
	for _, x := range sFluidContainers {
		IsFluidContainer[x] = true
	}
	for _, x := range sFluids {
		IsFluid[x] = true
	}
	for _, x := range sStackable {
		IsStackable[x] = true
	}
}
