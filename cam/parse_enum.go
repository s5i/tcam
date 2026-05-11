package cam

// Opcode identifies the operation type.
type Opcode byte

const (
	LoginPlayerState      Opcode = 0x0A
	LoginError            Opcode = 0x14
	LoginWaitList         Opcode = 0x16
	Ping                  Opcode = 0x1E
	Map                   Opcode = 0x64
	MoveNorth             Opcode = 0x65
	MoveEast              Opcode = 0x66
	MoveSouth             Opcode = 0x67
	MoveWest              Opcode = 0x68
	TileUpdate            Opcode = 0x69
	TileItemAdd           Opcode = 0x6A
	TileItemUpdate        Opcode = 0x6B
	TileItemRemove        Opcode = 0x6C
	CreatureMove          Opcode = 0x6D
	ContainerOpen         Opcode = 0x6E
	ContainerClose        Opcode = 0x6F
	ContainerItemAdd      Opcode = 0x70
	ContainerItemUpdate   Opcode = 0x71
	ContainerItemRemove   Opcode = 0x72
	InventoryItemSet      Opcode = 0x78
	InventoryItemClear    Opcode = 0x79
	TradeOwn              Opcode = 0x7D
	TradeCounter          Opcode = 0x7E
	TradeClose            Opcode = 0x7F
	EffectLight           Opcode = 0x82
	EffectGraphical       Opcode = 0x83
	EffectText            Opcode = 0x84
	EffectMissile         Opcode = 0x85
	CreatureSquare        Opcode = 0x86
	CreatureHealth        Opcode = 0x8C
	CreatureLight         Opcode = 0x8D
	CreatureOutfit        Opcode = 0x8E
	CreatureSpeed         Opcode = 0x8F
	CreatureSkull         Opcode = 0x90
	CreatureParty         Opcode = 0x91
	PromptTextUpdate      Opcode = 0x96
	PromptHouseList       Opcode = 0x97
	PlayerStats           Opcode = 0xA0
	PlayerSkills          Opcode = 0xA1
	PlayerIcons           Opcode = 0xA2
	TargetClear           Opcode = 0xA3
	CreatureMessage       Opcode = 0xAA
	ChannelList           Opcode = 0xAB
	ChannelOpen           Opcode = 0xAC
	PrivateChannelOpen    Opcode = 0xAD
	RuleViolationsChannel Opcode = 0xAE
	RuleViolationsRemove  Opcode = 0xAF
	RuleViolationCancel   Opcode = 0xB0
	RuleViolationsLock    Opcode = 0xB1
	PrivateChannelCreate  Opcode = 0xB2
	PrivateChannelClose   Opcode = 0xB3
	Message               Opcode = 0xB4
	MoveCancel            Opcode = 0xB5
	MoveFloorUp           Opcode = 0xBE
	MoveFloorDown         Opcode = 0xBF
	PromptChooseOutfit    Opcode = 0xC8
	VIPState              Opcode = 0xD2
	VIPLogin              Opcode = 0xD3
	VIPLogout             Opcode = 0xD4
)
