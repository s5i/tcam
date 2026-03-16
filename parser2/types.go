package parser2

type TByte uint8
type TWord uint16
type TQuad uint32
type TString string
type TOutfit struct {
	ID     TWord
	Type   TWord
	Colors TQuad
}

// TPosition represents a map coordinate.
type TPosition struct {
	X TWord
	Y TWord
	Z TByte
}

// TItem represents an item on the map or in a container.
type TItem struct {
	TypeID    TWord
	ExtraByte TByte // fluid type or stack count; valid when HasExtra is true
	HasExtra  bool
}

// TCreature represents a creature as sent over the wire.
// Depending on the known-state (Unknown=97, Outdated=98, Known=99),
// only some fields are populated.
type TCreature struct {
	KnownState TWord // 97, 98, or 99

	// KnownState == 97 (unknown / new creature)
	RemoveID TQuad
	ID       TQuad
	Name     TString

	// KnownState == 97 or 98 (full update)
	Health          TByte
	Direction       TByte
	Outfit          TOutfit
	LightBrightness TByte
	LightColor      TByte
	Speed           TWord
	Skull           TByte
	Party           TByte

	// KnownState == 99 (already known)
	// ID and Direction are set; other fields are zero.
}

// TMapObject is either a creature or an item on the map.
type TMapObject struct {
	Creature *TCreature
	Item     *TItem
}

// TPlayerData contains the player status bar information.
type TPlayerData struct {
	HitPoints         TWord
	MaxHitPoints      TWord
	Capacity          TWord
	Experience        TQuad
	Level             TWord
	LevelPercent      TByte
	ManaPoints        TWord
	MaxManaPoints     TWord
	MagicLevel        TByte
	MagicLevelPercent TByte
	SoulPoints        TByte
}

// TSkillEntry is a single combat skill.
type TSkillEntry struct {
	Level   TByte
	Percent TByte
}

// TPlayerSkills contains all 7 combat skills.
type TPlayerSkills struct {
	Fist      TSkillEntry
	Club      TSkillEntry
	Sword     TSkillEntry
	Axe       TSkillEntry
	Distance  TSkillEntry
	Shielding TSkillEntry
	Fishing   TSkillEntry
}

// TContainer describes an open container and its contents.
type TContainer struct {
	ContainerNr TByte
	TypeID      TWord
	Name        TString
	Capacity    TByte
	HasParent   bool
	Items       []TItem
}

// TTradeOffer describes a trade offer (own or partner).
type TTradeOffer struct {
	Name  TString
	Items []TItem
}

// TChannel identifies a chat channel.
type TChannel struct {
	ID   TWord
	Name TString
}

// TTalk represents an incoming chat message.
// Depending on Mode, either Position or ChannelID is meaningful.
type TTalk struct {
	StatementID TQuad
	Sender      TString
	Mode        TByte

	// Modes that carry a map position (say, whisper, yell, etc.)
	Position TPosition

	// Modes that carry a channel ID.
	ChannelID TWord

	// Mode-specific extra data (e.g. gamemaster request data).
	ExtraData TQuad

	Text TString
}

// TEditText describes the edit-text dialog data.
type TEditText struct {
	ObjectID  TQuad
	TypeID    TWord
	MaxLength TWord
	Text      TString
	Editor    TString
}

// TEditList describes the edit-list dialog data.
type TEditList struct {
	Type TByte
	ID   TQuad
	Text TString
}

// TBuddyData contains buddy list information.
type TBuddyData struct {
	CharacterID TQuad
	Name        TString
	Online      bool
}

// TOutfitWindow describes the outfit selection dialog.
type TOutfitWindow struct {
	CurrentOutfit TOutfit
	FirstOutfit   TWord
	LastOutfit    TWord
}

// TInitGame is the initial game data sent after login.
type TInitGame struct {
	CreatureID    TQuad
	Beat          TWord
	CanReportBugs bool
}

// TGraphicalEffect represents a graphical effect at a position.
type TGraphicalEffect struct {
	Position TPosition
	Type     TByte
}

// TTextualEffect represents floating text at a position.
type TTextualEffect struct {
	Position TPosition
	Color    TByte
	Text     TString
}

// TMissileEffect represents a missile projectile.
type TMissileEffect struct {
	Origin      TPosition
	Destination TPosition
	Type        TByte
}

// TMessage represents a system/status message.
type TMessage struct {
	Mode TByte
	Text TString
}

// TFieldChange represents an add/change/delete on a map field.
type TFieldChange struct {
	Position   TPosition
	StackIndex TByte       // used by change/delete
	Object     *TMapObject // used by add/change
}

// TMoveCreature represents a creature moving on the map.
type TMoveCreature struct {
	OrigPosition TPosition
	OrigIndex    TByte
	DestPosition TPosition
}

// TAmbiente represents ambient light settings.
type TAmbiente struct {
	Brightness TByte
	Color      TByte
}

// TCreatureUpdate represents a single-field creature update.
type TCreatureUpdate struct {
	CreatureID TQuad
}

// TCreatureHealthUpdate is a creature health change.
type TCreatureHealthUpdate struct {
	CreatureID TQuad
	Health     TByte
}

// TCreatureLightUpdate is a creature light change.
type TCreatureLightUpdate struct {
	CreatureID TQuad
	Brightness TByte
	Color      TByte
}

// TCreatureOutfitUpdate is a creature outfit change.
type TCreatureOutfitUpdate struct {
	CreatureID TQuad
	Outfit     TOutfit
}

// TCreatureSpeedUpdate is a creature speed change.
type TCreatureSpeedUpdate struct {
	CreatureID TQuad
	Speed      TWord
}

// TMarkCreature is a creature mark/color indicator.
type TMarkCreature struct {
	CreatureID TQuad
	Color      TByte
}
