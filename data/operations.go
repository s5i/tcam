package data

// Login (0x0A).
type Login struct {
	PlayerID    uint32
	AccessLevel byte
}

// DisconnectClient (0x14).
type DisconnectClient struct {
	Message string
}

// WaitList (0x16).
type WaitList struct {
	Message string
	Time    byte
}

// Ping (0x1E).
type Ping struct{}

// MapDescription (0x64).
type MapDescription struct {
	PlayerPos Location
	Tiles     []Tile
}

// MoveNorth (0x65).
type MoveNorth struct {
	Tiles []Tile
}

// MoveEast (0x66).
type MoveEast struct {
	Tiles []Tile
}

// MoveSouth (0x67).
type MoveSouth struct {
	Tiles []Tile
}

// MoveWest (0x68).
type MoveWest struct {
	Tiles []Tile
}

// UpdateTile (0x69).
type UpdateTile struct {
	Location Location
	Tile     *Tile // nil if tile was cleared (thingId == 0xFF01)
}

// AddTileItem (0x6A).
type AddTileItem struct {
	Location Location
	Thing    Thing
}

// UpdateTileItem (0x6B).
type UpdateTileItem struct {
	Location   Location
	StackIndex byte
	Thing      Thing
}

// RemoveTileItem (0x6C).
type RemoveTileItem struct {
	Location   Location
	StackIndex byte
}

// MoveCreature (0x6D).
type MoveCreature struct {
	OldLocation Location
	OldStack    byte
	NewLocation Location
}

// Container (0x6E).
type Container struct {
	ContainerID byte
	ItemID      uint16
	Name        string
	Volume      byte
	HasParent   byte
	Items       []Thing
}

// CloseContainer (0x6F).
type CloseContainer struct {
	ContainerID byte
}

// AddContainerItem (0x70).
type AddContainerItem struct {
	ContainerID byte
	Thing       Thing
}

// UpdateContainerItem (0x71).
type UpdateContainerItem struct {
	ContainerID byte
	Slot        byte
	Thing       Thing
}

// RemoveContainerItem (0x72).
type RemoveContainerItem struct {
	ContainerID byte
	Slot        byte
}

// InventorySetItem (0x78).
type InventorySetItem struct {
	Slot byte
	Item Item
}

// InventoryClearItem (0x79).
type InventoryClearItem struct {
	Slot byte
}

// TradeItemRequest (0x7D/0x7E).
type TradeItemRequest struct {
	Name  string
	Items []Thing
}

// CloseTrade (0x7F).
type CloseTrade struct{}

// WorldLight (0x82).
type WorldLight struct {
	Level byte
	Color byte
}

// MagicEffect (0x83).
type MagicEffect struct {
	Location Location
	Effect   byte
}

// AnimatedText (0x84).
type AnimatedText struct {
	Location Location
	Color    byte
	Text     string
}

// DistanceShoot (0x85).
type DistanceShoot struct {
	From   Location
	To     Location
	Effect byte
}

// CreatureSquare (0x86).
type CreatureSquare struct {
	CreatureID uint32
	Color      byte
}

// CreatureHealth (0x8C).
type CreatureHealth struct {
	CreatureID uint32
	Health     byte
}

// CreatureLight (0x8D).
type CreatureLight struct {
	CreatureID uint32
	Level      byte
	Color      byte
}

// CreatureOutfit (0x8E).
type CreatureOutfit struct {
	CreatureID uint32
	Outfit     Outfit
}

// ChangeSpeed (0x8F).
type ChangeSpeed struct {
	CreatureID uint32
	Speed      uint16
}

// CreatureSkull (0x90).
type CreatureSkull struct {
	CreatureID uint32
	Skull      byte
}

// CreatureShield (0x91).
type CreatureShield struct {
	CreatureID uint32
	Shield     byte
}

// TextWindow (0x96).
type TextWindow struct {
	WindowID uint32
	ItemID   uint16
	MaxLen   uint16
	Text     string
	Author   string
}

// HouseWindow (0x97).
type HouseWindow struct {
	Unknown byte
	ID      uint32
	Text    string
}

// PlayerStats (0xA0).
type PlayerStats struct {
	HP          uint16
	MaxHP       uint16
	Capacity    uint16
	Exp         uint32
	Level       byte
	LevelPct    byte
	Mana        uint16
	MaxMana     uint16
	MagicLvl    byte
	MagicLvlPct byte
	Soul        uint16
}

// SkillValue represents a single skill's level and percent.
type SkillValue struct {
	Level   byte
	Percent byte
}

// PlayerSkills (0xA1).
type PlayerSkills struct {
	Skills [7]SkillValue
}

// PlayerIcons (0xA2).
type PlayerIcons struct {
	Icons byte
}

// CancelTarget (0xA3).
type CancelTarget struct{}

// CreatureSpeak (0xAA).
type CreatureSpeak struct {
	StatementID uint32
	Name        string
	Type        byte
	Location    *Location // for types 1,2,3,0x10,0x11
	ChannelID   *uint16   // for types 5,6,0xA,0xC,0xE
	Text        string
}

// ChannelEntry represents a single channel in a ChannelsDialog.
type ChannelEntry struct {
	ID   uint16
	Name string
}

// ChannelsDialog (0xAB).
type ChannelsDialog struct {
	Channels []ChannelEntry
}

// Channel (0xAC).
type Channel struct {
	ID   uint16
	Name string
}

// OpenPrivateChannel (0xAD).
type OpenPrivateChannel struct {
	Name string
}

// RuleViolationsChannel (0xAE).
type RuleViolationsChannel struct {
	Size uint16
}

// RemoveReport (0xAF).
type RemoveReport struct {
	Name string
}

// RuleViolationCancel (0xB0).
type RuleViolationCancel struct {
	Name string
}

// LockRuleViolation (0xB1).
type LockRuleViolation struct{}

// CreatePrivateChannel (0xB2).
type CreatePrivateChannel struct {
	ID   uint16
	Name string
}

// ClosePrivate (0xB3).
type ClosePrivate struct {
	ChannelID uint16
}

// TextMessage (0xB4).
type TextMessage struct {
	Type byte
	Text string
}

// CancelWalk (0xB5).
type CancelWalk struct {
	Direction Direction
}

// FloorChangeUp (0xBE).
type FloorChangeUp struct {
	Tiles []Tile
}

// FloorChangeDown (0xBF).
type FloorChangeDown struct {
	Tiles []Tile
}

// OutfitWindow (0xC8).
type OutfitWindow struct {
	Outfit     Outfit
	OutfitStart uint16
	OutfitEnd   uint16
}

// VIP (0xD2).
type VIP struct {
	ID     uint32
	Name   string
	Online byte
}

// VIPLogin (0xD3).
type VIPLogin struct {
	ID uint32
}

// VIPLogout (0xD4).
type VIPLogout struct {
	ID uint32
}
