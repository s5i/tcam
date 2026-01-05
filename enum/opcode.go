package enum

import "fmt"

type OpCode uint8

func (c OpCode) String() string {
	if str, ok := opCodeMap[c]; ok {
		return str
	}
	return fmt.Sprintf("Unknown-%d", c)
}

const (
	OpCodeLoginOrPendingState  OpCode = 10
	OpCodeGMActions            OpCode = 11
	OpCodeUpdateNeeded         OpCode = 12
	OpCodeLoginError           OpCode = 13
	OpCodeLoginAdvice          OpCode = 14
	OpCodeLoginWait            OpCode = 15
	OpCodeLoginToken           OpCode = 16
	OpCodePing                 OpCode = 29
	OpCodePingBack             OpCode = 30
	OpCodeChallenge            OpCode = 31
	OpCodeNewPing              OpCode = 32
	OpCodeDeath                OpCode = 40
	OpCodeFullMap              OpCode = 100
	OpCodeMapTopRow            OpCode = 101
	OpCodeMapRightRow          OpCode = 102
	OpCodeMapBottomRow         OpCode = 103
	OpCodeMapLeftRow           OpCode = 104
	OpCodeUpdateTile           OpCode = 105
	OpCodeCreateOnMap          OpCode = 106
	OpCodeChangeOnMap          OpCode = 107
	OpCodeDeleteOnMap          OpCode = 108
	OpCodeMoveCreature         OpCode = 109
	OpCodeOpenContainer        OpCode = 110
	OpCodeCloseContainer       OpCode = 111
	OpCodeCreateContainer      OpCode = 112
	OpCodeChangeInContainer    OpCode = 113
	OpCodeDeleteInContainer    OpCode = 114
	OpCodeSetInventory         OpCode = 120
	OpCodeDeleteInventory      OpCode = 121
	OpCodeOpenNpcTrade         OpCode = 122
	OpCodePlayerGoods          OpCode = 123
	OpCodeCloseNpcTrade        OpCode = 124
	OpCodeOwnTrade             OpCode = 125
	OpCodeCounterTrade         OpCode = 126
	OpCodeCloseTrade           OpCode = 127
	OpCodeAmbient              OpCode = 130
	OpCodeGraphicalEffect      OpCode = 131
	OpCodeTextEffect           OpCode = 132
	OpCodeMissleEffect         OpCode = 133
	OpCodeMarkCreature         OpCode = 134
	OpCodeTrappers             OpCode = 135
	OpCodeCreatureHealth       OpCode = 140
	OpCodeCreatureLight        OpCode = 141
	OpCodeCreatureOutfit       OpCode = 142
	OpCodeCreatureSpeed        OpCode = 143
	OpCodeCreatureSkull        OpCode = 144
	OpCodeCreatureParty        OpCode = 145
	OpCodeCreatureUnpass       OpCode = 146
	OpCodeCreatureMarks        OpCode = 147
	OpCodePlayerHelpers        OpCode = 148
	OpCodeCreatureType         OpCode = 149
	OpCodeEditText             OpCode = 150
	OpCodeEditList             OpCode = 151
	OpCodeBlessings            OpCode = 156
	OpCodePreset               OpCode = 157
	OpCodePremiumTrigger       OpCode = 158
	OpCodePlayerDataBasic      OpCode = 159
	OpCodePlayerData           OpCode = 160
	OpCodePlayerSkills         OpCode = 161
	OpCodePlayerState          OpCode = 162
	OpCodeClearTarget          OpCode = 163
	OpCodePlayerModes          OpCode = 167
	OpCodeSpellDelay           OpCode = 164
	OpCodeSpellGroupDelay      OpCode = 165
	OpCodeMultiUseDelay        OpCode = 166
	OpCodeSetStoreDeepLink     OpCode = 168
	OpCodeTalk                 OpCode = 170
	OpCodeChannels             OpCode = 171
	OpCodeOpenChannel          OpCode = 172
	OpCodeOpenPrivateChannel   OpCode = 173
	OpCodeRuleViolationChannel OpCode = 174
	OpCodeRuleViolationRemove  OpCode = 175
	OpCodeRuleViolationCancel  OpCode = 176
	OpCodeRuleViolationLock    OpCode = 177
	OpCodeOpenOwnChannel       OpCode = 178
	OpCodeCloseChannel         OpCode = 179
	OpCodeTextMessage          OpCode = 180
	OpCodeCancelWalk           OpCode = 181
	OpCodeWalkWait             OpCode = 182
	OpCodeFloorChangeUp        OpCode = 190
	OpCodeFloorChangeDown      OpCode = 191
	OpCodeChooseOutfit         OpCode = 200
	OpCodeVipAdd               OpCode = 210
	OpCodeVipState             OpCode = 211
	OpCodeVipLogout            OpCode = 212
	OpCodeTutorialHint         OpCode = 220
	OpCodeAutomapFlag          OpCode = 221
	OpCodeQuestLog             OpCode = 240
	OpCodeQuestLine            OpCode = 241
)

var opCodeMap = map[OpCode]string{
	OpCodeLoginOrPendingState:  "LoginOrPendingState",
	OpCodeGMActions:            "GMActions",
	OpCodeUpdateNeeded:         "UpdateNeeded",
	OpCodeLoginError:           "LoginError",
	OpCodeLoginAdvice:          "LoginAdvice",
	OpCodeLoginWait:            "LoginWait",
	OpCodeLoginToken:           "LoginToken",
	OpCodePing:                 "Ping",
	OpCodePingBack:             "PingBack",
	OpCodeChallenge:            "Challenge",
	OpCodeNewPing:              "NewPing",
	OpCodeDeath:                "Death",
	OpCodeFullMap:              "FullMap",
	OpCodeMapTopRow:            "MapTopRow",
	OpCodeMapRightRow:          "MapRightRow",
	OpCodeMapBottomRow:         "MapBottomRow",
	OpCodeMapLeftRow:           "MapLeftRow",
	OpCodeUpdateTile:           "UpdateTile",
	OpCodeCreateOnMap:          "CreateOnMap",
	OpCodeChangeOnMap:          "ChangeOnMap",
	OpCodeDeleteOnMap:          "DeleteOnMap",
	OpCodeMoveCreature:         "MoveCreature",
	OpCodeOpenContainer:        "OpenContainer",
	OpCodeCloseContainer:       "CloseContainer",
	OpCodeCreateContainer:      "CreateContainer",
	OpCodeChangeInContainer:    "ChangeInContainer",
	OpCodeDeleteInContainer:    "DeleteInContainer",
	OpCodeSetInventory:         "SetInventory",
	OpCodeDeleteInventory:      "DeleteInventory",
	OpCodeOpenNpcTrade:         "OpenNpcTrade",
	OpCodePlayerGoods:          "PlayerGoods",
	OpCodeCloseNpcTrade:        "CloseNpcTrade",
	OpCodeOwnTrade:             "OwnTrade",
	OpCodeCounterTrade:         "CounterTrade",
	OpCodeCloseTrade:           "CloseTrade",
	OpCodeAmbient:              "Ambient",
	OpCodeGraphicalEffect:      "GraphicalEffect",
	OpCodeTextEffect:           "TextEffect",
	OpCodeMissleEffect:         "MissleEffect",
	OpCodeMarkCreature:         "MarkCreature",
	OpCodeTrappers:             "Trappers",
	OpCodeCreatureHealth:       "CreatureHealth",
	OpCodeCreatureLight:        "CreatureLight",
	OpCodeCreatureOutfit:       "CreatureOutfit",
	OpCodeCreatureSpeed:        "CreatureSpeed",
	OpCodeCreatureSkull:        "CreatureSkull",
	OpCodeCreatureParty:        "CreatureParty",
	OpCodeCreatureUnpass:       "CreatureUnpass",
	OpCodeCreatureMarks:        "CreatureMarks",
	OpCodePlayerHelpers:        "PlayerHelpers",
	OpCodeCreatureType:         "CreatureType",
	OpCodeEditText:             "EditText",
	OpCodeEditList:             "EditList",
	OpCodeBlessings:            "Blessings",
	OpCodePreset:               "Preset",
	OpCodePremiumTrigger:       "PremiumTrigger",
	OpCodePlayerDataBasic:      "PlayerDataBasic",
	OpCodePlayerData:           "PlayerData",
	OpCodePlayerSkills:         "PlayerSkills",
	OpCodePlayerState:          "PlayerState",
	OpCodeClearTarget:          "ClearTarget",
	OpCodePlayerModes:          "PlayerModes",
	OpCodeSpellDelay:           "SpellDelay",
	OpCodeSpellGroupDelay:      "SpellGroupDelay",
	OpCodeMultiUseDelay:        "MultiUseDelay",
	OpCodeSetStoreDeepLink:     "SetStoreDeepLink",
	OpCodeTalk:                 "Talk",
	OpCodeChannels:             "Channels",
	OpCodeOpenChannel:          "OpenChannel",
	OpCodeOpenPrivateChannel:   "OpenPrivateChannel",
	OpCodeRuleViolationChannel: "RuleViolationChannel",
	OpCodeRuleViolationRemove:  "RuleViolationRemove",
	OpCodeRuleViolationCancel:  "RuleViolationCancel",
	OpCodeRuleViolationLock:    "RuleViolationLock",
	OpCodeOpenOwnChannel:       "OpenOwnChannel",
	OpCodeCloseChannel:         "CloseChannel",
	OpCodeTextMessage:          "TextMessage",
	OpCodeCancelWalk:           "CancelWalk",
	OpCodeWalkWait:             "WalkWait",
	OpCodeFloorChangeUp:        "FloorChangeUp",
	OpCodeFloorChangeDown:      "FloorChangeDown",
	OpCodeChooseOutfit:         "ChooseOutfit",
	OpCodeVipAdd:               "VipAdd",
	OpCodeVipState:             "VipState",
	OpCodeVipLogout:            "VipLogout",
	OpCodeTutorialHint:         "TutorialHint",
	OpCodeAutomapFlag:          "AutomapFlag",
	OpCodeQuestLog:             "QuestLog",
	OpCodeQuestLine:            "QuestLine",
}
