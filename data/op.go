package data

import "time"

// Operation represents an abstract operation.
// Type-switch into concrete operation types for details.
type Operation interface {
	isOperation()
}

// LoginPlayerState (0x0A).
type LoginPlayerState struct {
	TimeOffset  time.Duration
	PlayerID    uint32
	AccessLevel byte
}

// LoginError (0x14).
type LoginError struct {
	TimeOffset time.Duration
	Message    string
}

// LoginWaitList (0x16).
type LoginWaitList struct {
	TimeOffset time.Duration
	Message    string
	Time       byte
}

// Ping (0x1E).
type Ping struct {
	TimeOffset time.Duration
}

// Map (0x64).
type Map struct {
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

// TileUpdate (0x69).
type TileUpdate struct {
	TimeOffset time.Duration
	Location   Location
	Tile       Tile
	HasTile    bool // If false, tile was cleared.
}

// TileItemAdd (0x6A).
type TileItemAdd struct {
	TimeOffset time.Duration
	Location   Location
	Thing      Thing
}

// TileItemUpdate (0x6B).
type TileItemUpdate struct {
	TimeOffset time.Duration
	Location   Location
	StackIndex byte
	Thing      Thing
}

// TileItemRemove (0x6C).
type TileItemRemove struct {
	TimeOffset time.Duration
	Location   Location
	StackIndex byte
}

// CreatureMove (0x6D).
type CreatureMove struct {
	TimeOffset  time.Duration
	OldLocation Location
	OldStack    byte
	NewLocation Location
}

// ContainerOpen (0x6E).
type ContainerOpen struct {
	TimeOffset  time.Duration
	ContainerID byte
	ItemID      uint16
	Name        string
	Volume      byte
	HasParent   byte
	Items       []Thing
}

// ContainerClose (0x6F).
type ContainerClose struct {
	TimeOffset  time.Duration
	ContainerID byte
}

// ContainerItemAdd (0x70).
type ContainerItemAdd struct {
	TimeOffset  time.Duration
	ContainerID byte
	Thing       Thing
}

// ContainerItemUpdate (0x71).
type ContainerItemUpdate struct {
	TimeOffset  time.Duration
	ContainerID byte
	Slot        byte
	Thing       Thing
}

// ContainerItemRemove (0x72).
type ContainerItemRemove struct {
	TimeOffset  time.Duration
	ContainerID byte
	Slot        byte
}

// InventoryItemSet (0x78).
type InventoryItemSet struct {
	TimeOffset time.Duration
	Slot       byte
	Item       Item
}

// InventoryItemClear (0x79).
type InventoryItemClear struct {
	TimeOffset time.Duration
	Slot       byte
}

// TradeOwn (0x7D).
type TradeOwn struct {
	TimeOffset time.Duration
	Name       string
	Items      []Thing
}

// TradeCounter (0x7E).
type TradeCounter struct {
	TimeOffset time.Duration
	Name       string
	Items      []Thing
}

// TradeClose (0x7F).
type TradeClose struct {
	TimeOffset time.Duration
}

// EffectLight (0x82).
type EffectLight struct {
	TimeOffset time.Duration
	Level      byte
	Color      byte
}

// EffectGraphical (0x83).
type EffectGraphical struct {
	TimeOffset time.Duration
	Location   Location
	Effect     byte
}

// EffectText (0x84).
type EffectText struct {
	TimeOffset time.Duration
	Location   Location
	Color      byte
	Text       string
}

// EffectMissile (0x85).
type EffectMissile struct {
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

// CreatureSpeed (0x8F).
type CreatureSpeed struct {
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

// CreatureParty (0x91).
type CreatureParty struct {
	TimeOffset time.Duration
	CreatureID uint32
	Shield     byte
}

// PromptTextUpdate (0x96).
type PromptTextUpdate struct {
	TimeOffset time.Duration
	WindowID   uint32
	ItemID     uint16
	MaxLen     uint16
	Text       string
	Author     string
}

// PromptHouseList (0x97).
type PromptHouseList struct {
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

// TargetClear (0xA3).
type TargetClear struct {
	TimeOffset time.Duration
}

// CreatureMessage (0xAA).
type CreatureMessage struct {
	TimeOffset  time.Duration
	StatementID uint32
	Name        string
	Type        byte
	Location    *Location // for types 1,2,3,0x10,0x11
	ChannelID   *uint16   // for types 5,6,0xA,0xC,0xE
	Text        string
}

// ChannelEntry represents a single channel in a ChannelList.
type ChannelEntry struct {
	ID   uint16
	Name string
}

// ChannelList (0xAB).
type ChannelList struct {
	TimeOffset time.Duration
	Channels   []ChannelEntry
}

// ChannelOpen (0xAC).
type ChannelOpen struct {
	TimeOffset time.Duration
	ID         uint16
	Name       string
}

// PrivateChannelOpen (0xAD).
type PrivateChannelOpen struct {
	TimeOffset time.Duration
	Name       string
}

// RuleViolationsChannel (0xAE).
type RuleViolationsChannel struct {
	TimeOffset time.Duration
	Size       uint16
}

// RuleViolationsRemove (0xAF).
type RuleViolationsRemove struct {
	TimeOffset time.Duration
	Name       string
}

// RuleViolationCancel (0xB0).
type RuleViolationCancel struct {
	TimeOffset time.Duration
	Name       string
}

// RuleViolationsLock (0xB1).
type RuleViolationsLock struct {
	TimeOffset time.Duration
}

// PrivateChannelCreate (0xB2).
type PrivateChannelCreate struct {
	TimeOffset time.Duration
	ID         uint16
	Name       string
}

// PrivateChannelClose (0xB3).
type PrivateChannelClose struct {
	TimeOffset time.Duration
	ChannelID  uint16
}

// Message (0xB4).
type Message struct {
	TimeOffset time.Duration
	Type       byte
	Text       string
}

// MoveCancel (0xB5).
type MoveCancel struct {
	TimeOffset time.Duration
	Direction  Direction
}

// MoveFloorUp (0xBE).
type MoveFloorUp struct {
	TimeOffset time.Duration
	Tiles      []Tile
}

// MoveFloorDown (0xBF).
type MoveFloorDown struct {
	TimeOffset time.Duration
	Tiles      []Tile
}

// PromptChooseOutfit (0xC8).
type PromptChooseOutfit struct {
	TimeOffset  time.Duration
	Outfit      Outfit
	OutfitStart uint16
	OutfitEnd   uint16
}

// VIPState (0xD2).
type VIPState struct {
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

func (LoginPlayerState) isOperation()      {}
func (LoginError) isOperation()            {}
func (LoginWaitList) isOperation()         {}
func (Ping) isOperation()                  {}
func (Map) isOperation()                   {}
func (MoveNorth) isOperation()             {}
func (MoveEast) isOperation()              {}
func (MoveSouth) isOperation()             {}
func (MoveWest) isOperation()              {}
func (TileUpdate) isOperation()            {}
func (TileItemAdd) isOperation()           {}
func (TileItemUpdate) isOperation()        {}
func (TileItemRemove) isOperation()        {}
func (CreatureMove) isOperation()          {}
func (ContainerOpen) isOperation()         {}
func (ContainerClose) isOperation()        {}
func (ContainerItemAdd) isOperation()      {}
func (ContainerItemUpdate) isOperation()   {}
func (ContainerItemRemove) isOperation()   {}
func (InventoryItemSet) isOperation()      {}
func (InventoryItemClear) isOperation()    {}
func (TradeOwn) isOperation()              {}
func (TradeCounter) isOperation()          {}
func (TradeClose) isOperation()            {}
func (EffectLight) isOperation()           {}
func (EffectGraphical) isOperation()       {}
func (EffectText) isOperation()            {}
func (EffectMissile) isOperation()         {}
func (CreatureSquare) isOperation()        {}
func (CreatureHealth) isOperation()        {}
func (CreatureLight) isOperation()         {}
func (CreatureOutfit) isOperation()        {}
func (CreatureSpeed) isOperation()         {}
func (CreatureSkull) isOperation()         {}
func (CreatureParty) isOperation()         {}
func (PromptTextUpdate) isOperation()      {}
func (PromptHouseList) isOperation()       {}
func (PlayerStats) isOperation()           {}
func (SkillValue) isOperation()            {}
func (PlayerSkills) isOperation()          {}
func (PlayerIcons) isOperation()           {}
func (TargetClear) isOperation()           {}
func (CreatureMessage) isOperation()       {}
func (ChannelEntry) isOperation()          {}
func (ChannelList) isOperation()           {}
func (ChannelOpen) isOperation()           {}
func (PrivateChannelOpen) isOperation()    {}
func (RuleViolationsChannel) isOperation() {}
func (RuleViolationsRemove) isOperation()  {}
func (RuleViolationCancel) isOperation()   {}
func (RuleViolationsLock) isOperation()    {}
func (PrivateChannelCreate) isOperation()  {}
func (PrivateChannelClose) isOperation()   {}
func (Message) isOperation()               {}
func (MoveCancel) isOperation()            {}
func (MoveFloorUp) isOperation()           {}
func (MoveFloorDown) isOperation()         {}
func (PromptChooseOutfit) isOperation()    {}
func (VIPState) isOperation()              {}
func (VIPLogin) isOperation()              {}
func (VIPLogout) isOperation()             {}
