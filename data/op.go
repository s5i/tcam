package data

import "time"

// Operation represents an abstract operation.
// Type-switch into concrete operation types for details.
type Operation interface {
	isOperation()
}

// Login (0x0A).
type Login struct {
	TimeOffset  time.Duration
	PlayerID    uint32
	AccessLevel byte
}

// DisconnectClient (0x14).
type DisconnectClient struct {
	TimeOffset time.Duration
	Message    string
}

// WaitList (0x16).
type WaitList struct {
	TimeOffset time.Duration
	Message    string
	Time       byte
}

// Ping (0x1E).
type Ping struct {
	TimeOffset time.Duration
}

// MapDescription (0x64).
type MapDescription struct {
	TimeOffset time.Duration
	PlayerPos  Location
	Tiles      []Tile
}

// MoveNorth (0x65).
type MoveNorth struct {
	TimeOffset time.Duration
	Tiles      []Tile
}

// MoveEast (0x66).
type MoveEast struct {
	TimeOffset time.Duration
	Tiles      []Tile
}

// MoveSouth (0x67).
type MoveSouth struct {
	TimeOffset time.Duration
	Tiles      []Tile
}

// MoveWest (0x68).
type MoveWest struct {
	TimeOffset time.Duration
	Tiles      []Tile
}

// UpdateTile (0x69).
type UpdateTile struct {
	TimeOffset time.Duration
	Location   Location
	Tile       *Tile // nil if tile was cleared (thingId == 0xFF01)
}

// AddTileItem (0x6A).
type AddTileItem struct {
	TimeOffset time.Duration
	Location   Location
	Thing      Thing
}

// UpdateTileItem (0x6B).
type UpdateTileItem struct {
	TimeOffset time.Duration
	Location   Location
	StackIndex byte
	Thing      Thing
}

// RemoveTileItem (0x6C).
type RemoveTileItem struct {
	TimeOffset time.Duration
	Location   Location
	StackIndex byte
}

// MoveCreature (0x6D).
type MoveCreature struct {
	TimeOffset  time.Duration
	OldLocation Location
	OldStack    byte
	NewLocation Location
}

// Container (0x6E).
type Container struct {
	TimeOffset  time.Duration
	ContainerID byte
	ItemID      uint16
	Name        string
	Volume      byte
	HasParent   byte
	Items       []Thing
}

// CloseContainer (0x6F).
type CloseContainer struct {
	TimeOffset  time.Duration
	ContainerID byte
}

// AddContainerItem (0x70).
type AddContainerItem struct {
	TimeOffset  time.Duration
	ContainerID byte
	Thing       Thing
}

// UpdateContainerItem (0x71).
type UpdateContainerItem struct {
	TimeOffset  time.Duration
	ContainerID byte
	Slot        byte
	Thing       Thing
}

// RemoveContainerItem (0x72).
type RemoveContainerItem struct {
	TimeOffset  time.Duration
	ContainerID byte
	Slot        byte
}

// InventorySetItem (0x78).
type InventorySetItem struct {
	TimeOffset time.Duration
	Slot       byte
	Item       Item
}

// InventoryClearItem (0x79).
type InventoryClearItem struct {
	TimeOffset time.Duration
	Slot       byte
}

// TradeItemRequest (0x7D/0x7E).
type TradeItemRequest struct {
	TimeOffset time.Duration
	Name       string
	Items      []Thing
}

// CloseTrade (0x7F).
type CloseTrade struct {
	TimeOffset time.Duration
}

// WorldLight (0x82).
type WorldLight struct {
	TimeOffset time.Duration
	Level      byte
	Color      byte
}

// MagicEffect (0x83).
type MagicEffect struct {
	TimeOffset time.Duration
	Location   Location
	Effect     byte
}

// AnimatedText (0x84).
type AnimatedText struct {
	TimeOffset time.Duration
	Location   Location
	Color      byte
	Text       string
}

// DistanceShoot (0x85).
type DistanceShoot struct {
	TimeOffset time.Duration
	From       Location
	To         Location
	Effect     byte
}

// CreatureSquare (0x86).
type CreatureSquare struct {
	TimeOffset time.Duration
	CreatureID uint32
	Color      byte
}

// CreatureHealth (0x8C).
type CreatureHealth struct {
	TimeOffset time.Duration
	CreatureID uint32
	Health     byte
}

// CreatureLight (0x8D).
type CreatureLight struct {
	TimeOffset time.Duration
	CreatureID uint32
	Level      byte
	Color      byte
}

// CreatureOutfit (0x8E).
type CreatureOutfit struct {
	TimeOffset time.Duration
	CreatureID uint32
	Outfit     Outfit
}

// ChangeSpeed (0x8F).
type ChangeSpeed struct {
	TimeOffset time.Duration
	CreatureID uint32
	Speed      uint16
}

// CreatureSkull (0x90).
type CreatureSkull struct {
	TimeOffset time.Duration
	CreatureID uint32
	Skull      byte
}

// CreatureShield (0x91).
type CreatureShield struct {
	TimeOffset time.Duration
	CreatureID uint32
	Shield     byte
}

// TextWindow (0x96).
type TextWindow struct {
	TimeOffset time.Duration
	WindowID   uint32
	ItemID     uint16
	MaxLen     uint16
	Text       string
	Author     string
}

// HouseWindow (0x97).
type HouseWindow struct {
	TimeOffset time.Duration
	Unknown    byte
	ID         uint32
	Text       string
}

// PlayerStats (0xA0).
type PlayerStats struct {
	TimeOffset  time.Duration
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
	TimeOffset time.Duration
	Skills     [7]SkillValue
}

// PlayerIcons (0xA2).
type PlayerIcons struct {
	TimeOffset time.Duration
	Icons      byte
}

// CancelTarget (0xA3).
type CancelTarget struct {
	TimeOffset time.Duration
}

// CreatureSpeak (0xAA).
type CreatureSpeak struct {
	TimeOffset  time.Duration
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
	TimeOffset time.Duration
	Channels   []ChannelEntry
}

// Channel (0xAC).
type Channel struct {
	TimeOffset time.Duration
	ID         uint16
	Name       string
}

// OpenPrivateChannel (0xAD).
type OpenPrivateChannel struct {
	TimeOffset time.Duration
	Name       string
}

// RuleViolationsChannel (0xAE).
type RuleViolationsChannel struct {
	TimeOffset time.Duration
	Size       uint16
}

// RemoveReport (0xAF).
type RemoveReport struct {
	TimeOffset time.Duration
	Name       string
}

// RuleViolationCancel (0xB0).
type RuleViolationCancel struct {
	TimeOffset time.Duration
	Name       string
}

// LockRuleViolation (0xB1).
type LockRuleViolation struct {
	TimeOffset time.Duration
}

// CreatePrivateChannel (0xB2).
type CreatePrivateChannel struct {
	TimeOffset time.Duration
	ID         uint16
	Name       string
}

// ClosePrivate (0xB3).
type ClosePrivate struct {
	TimeOffset time.Duration
	ChannelID  uint16
}

// TextMessage (0xB4).
type TextMessage struct {
	TimeOffset time.Duration
	Type       byte
	Text       string
}

// CancelWalk (0xB5).
type CancelWalk struct {
	TimeOffset time.Duration
	Direction  Direction
}

// FloorChangeUp (0xBE).
type FloorChangeUp struct {
	TimeOffset time.Duration
	Tiles      []Tile
}

// FloorChangeDown (0xBF).
type FloorChangeDown struct {
	TimeOffset time.Duration
	Tiles      []Tile
}

// OutfitWindow (0xC8).
type OutfitWindow struct {
	TimeOffset  time.Duration
	Outfit      Outfit
	OutfitStart uint16
	OutfitEnd   uint16
}

// VIP (0xD2).
type VIP struct {
	TimeOffset time.Duration
	ID         uint32
	Name       string
	Online     byte
}

// VIPLogin (0xD3).
type VIPLogin struct {
	TimeOffset time.Duration
	ID         uint32
}

// VIPLogout (0xD4).
type VIPLogout struct {
	TimeOffset time.Duration
	ID         uint32
}

func (Login) isOperation()                 {}
func (DisconnectClient) isOperation()      {}
func (WaitList) isOperation()              {}
func (Ping) isOperation()                  {}
func (MapDescription) isOperation()        {}
func (MoveNorth) isOperation()             {}
func (MoveEast) isOperation()              {}
func (MoveSouth) isOperation()             {}
func (MoveWest) isOperation()              {}
func (UpdateTile) isOperation()            {}
func (AddTileItem) isOperation()           {}
func (UpdateTileItem) isOperation()        {}
func (RemoveTileItem) isOperation()        {}
func (MoveCreature) isOperation()          {}
func (Container) isOperation()             {}
func (CloseContainer) isOperation()        {}
func (AddContainerItem) isOperation()      {}
func (UpdateContainerItem) isOperation()   {}
func (RemoveContainerItem) isOperation()   {}
func (InventorySetItem) isOperation()      {}
func (InventoryClearItem) isOperation()    {}
func (TradeItemRequest) isOperation()      {}
func (CloseTrade) isOperation()            {}
func (WorldLight) isOperation()            {}
func (MagicEffect) isOperation()           {}
func (AnimatedText) isOperation()          {}
func (DistanceShoot) isOperation()         {}
func (CreatureSquare) isOperation()        {}
func (CreatureHealth) isOperation()        {}
func (CreatureLight) isOperation()         {}
func (CreatureOutfit) isOperation()        {}
func (ChangeSpeed) isOperation()           {}
func (CreatureSkull) isOperation()         {}
func (CreatureShield) isOperation()        {}
func (TextWindow) isOperation()            {}
func (HouseWindow) isOperation()           {}
func (PlayerStats) isOperation()           {}
func (SkillValue) isOperation()            {}
func (PlayerSkills) isOperation()          {}
func (PlayerIcons) isOperation()           {}
func (CancelTarget) isOperation()          {}
func (CreatureSpeak) isOperation()         {}
func (ChannelEntry) isOperation()          {}
func (ChannelsDialog) isOperation()        {}
func (Channel) isOperation()               {}
func (OpenPrivateChannel) isOperation()    {}
func (RuleViolationsChannel) isOperation() {}
func (RemoveReport) isOperation()          {}
func (RuleViolationCancel) isOperation()   {}
func (LockRuleViolation) isOperation()     {}
func (CreatePrivateChannel) isOperation()  {}
func (ClosePrivate) isOperation()          {}
func (TextMessage) isOperation()           {}
func (CancelWalk) isOperation()            {}
func (FloorChangeUp) isOperation()         {}
func (FloorChangeDown) isOperation()       {}
func (OutfitWindow) isOperation()          {}
func (VIP) isOperation()                   {}
func (VIPLogin) isOperation()              {}
func (VIPLogout) isOperation()             {}
