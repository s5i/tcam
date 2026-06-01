package data

// OpType identifies the operation type.
type OpType byte

const (
	TLoginPlayerState      OpType = 0x0A
	TLoginError            OpType = 0x14
	TLoginWaitList         OpType = 0x16
	TPing                  OpType = 0x1E
	TMap                   OpType = 0x64
	TMoveNorth             OpType = 0x65
	TMoveEast              OpType = 0x66
	TMoveSouth             OpType = 0x67
	TMoveWest              OpType = 0x68
	TTileUpdate            OpType = 0x69
	TTileItemAdd           OpType = 0x6A
	TTileItemUpdate        OpType = 0x6B
	TTileItemRemove        OpType = 0x6C
	TCreatureMove          OpType = 0x6D
	TContainerOpen         OpType = 0x6E
	TContainerClose        OpType = 0x6F
	TContainerItemAdd      OpType = 0x70
	TContainerItemUpdate   OpType = 0x71
	TContainerItemRemove   OpType = 0x72
	TInventoryItemSet      OpType = 0x78
	TInventoryItemClear    OpType = 0x79
	TTradeOwn              OpType = 0x7D
	TTradeCounter          OpType = 0x7E
	TTradeClose            OpType = 0x7F
	TEffectLight           OpType = 0x82
	TEffectGraphical       OpType = 0x83
	TEffectText            OpType = 0x84
	TEffectMissile         OpType = 0x85
	TCreatureSquare        OpType = 0x86
	TCreatureHealth        OpType = 0x8C
	TCreatureLight         OpType = 0x8D
	TCreatureOutfit        OpType = 0x8E
	TCreatureSpeed         OpType = 0x8F
	TCreatureSkull         OpType = 0x90
	TCreatureParty         OpType = 0x91
	TPromptTextUpdate      OpType = 0x96
	TPromptHouseList       OpType = 0x97
	TPlayerStats           OpType = 0xA0
	TPlayerSkills          OpType = 0xA1
	TPlayerIcons           OpType = 0xA2
	TTargetClear           OpType = 0xA3
	TCreatureMessage       OpType = 0xAA
	TChannelList           OpType = 0xAB
	TChannelOpen           OpType = 0xAC
	TPrivateChannelOpen    OpType = 0xAD
	TRuleViolationsChannel OpType = 0xAE
	TRuleViolationsRemove  OpType = 0xAF
	TRuleViolationCancel   OpType = 0xB0
	TRuleViolationsLock    OpType = 0xB1
	TPrivateChannelCreate  OpType = 0xB2
	TPrivateChannelClose   OpType = 0xB3
	TMessage               OpType = 0xB4
	TMoveCancel            OpType = 0xB5
	TMoveFloorUp           OpType = 0xBE
	TMoveFloorDown         OpType = 0xBF
	TPromptChooseOutfit    OpType = 0xC8
	TVIPState              OpType = 0xD2
	TVIPLogin              OpType = 0xD3
	TVIPLogout             OpType = 0xD4
	TCamMetadata           OpType = 0xFF
)

var OpName = map[OpType]string{
	TLoginPlayerState:      "LoginPlayerState",
	TLoginError:            "LoginError",
	TLoginWaitList:         "LoginWaitList",
	TPing:                  "Ping",
	TMap:                   "Map",
	TMoveNorth:             "MoveNorth",
	TMoveEast:              "MoveEast",
	TMoveSouth:             "MoveSouth",
	TMoveWest:              "MoveWest",
	TTileUpdate:            "TileUpdate",
	TTileItemAdd:           "TileItemAdd",
	TTileItemUpdate:        "TileItemUpdate",
	TTileItemRemove:        "TileItemRemove",
	TCreatureMove:          "CreatureMove",
	TContainerOpen:         "ContainerOpen",
	TContainerClose:        "ContainerClose",
	TContainerItemAdd:      "ContainerItemAdd",
	TContainerItemUpdate:   "ContainerItemUpdate",
	TContainerItemRemove:   "ContainerItemRemove",
	TInventoryItemSet:      "InventoryItemSet",
	TInventoryItemClear:    "InventoryItemClear",
	TTradeOwn:              "TradeOwn",
	TTradeCounter:          "TradeCounter",
	TTradeClose:            "TradeClose",
	TEffectLight:           "EffectLight",
	TEffectGraphical:       "EffectGraphical",
	TEffectText:            "EffectText",
	TEffectMissile:         "EffectMissile",
	TCreatureSquare:        "CreatureSquare",
	TCreatureHealth:        "CreatureHealth",
	TCreatureLight:         "CreatureLight",
	TCreatureOutfit:        "CreatureOutfit",
	TCreatureSpeed:         "CreatureSpeed",
	TCreatureSkull:         "CreatureSkull",
	TCreatureParty:         "CreatureParty",
	TPromptTextUpdate:      "PromptTextUpdate",
	TPromptHouseList:       "PromptHouseList",
	TPlayerStats:           "PlayerStats",
	TPlayerSkills:          "PlayerSkills",
	TPlayerIcons:           "PlayerIcons",
	TTargetClear:           "TargetClear",
	TCreatureMessage:       "CreatureMessage",
	TChannelList:           "ChannelList",
	TChannelOpen:           "ChannelOpen",
	TPrivateChannelOpen:    "PrivateChannelOpen",
	TRuleViolationsChannel: "RuleViolationsChannel",
	TRuleViolationsRemove:  "RuleViolationsRemove",
	TRuleViolationCancel:   "RuleViolationCancel",
	TRuleViolationsLock:    "RuleViolationsLock",
	TPrivateChannelCreate:  "PrivateChannelCreate",
	TPrivateChannelClose:   "PrivateChannelClose",
	TMessage:               "Message",
	TMoveCancel:            "MoveCancel",
	TMoveFloorUp:           "MoveFloorUp",
	TMoveFloorDown:         "MoveFloorDown",
	TPromptChooseOutfit:    "PromptChooseOutfit",
	TVIPState:              "VIPState",
	TVIPLogin:              "VIPLogin",
	TVIPLogout:             "VIPLogout",
	TCamMetadata:           "CamMetadata",
}
