package enum

import "fmt"

type MessageMode uint8

func (c MessageMode) String() string {
	if str, ok := messageModeMap[c]; ok {
		return str
	}
	return fmt.Sprintf("Unknown-%d", c)
}

const (
	MessageModeMessageNone                  MessageMode = 0
	MessageModeMessageSay                   MessageMode = 1
	MessageModeMessageWhisper               MessageMode = 2
	MessageModeMessageYell                  MessageMode = 3
	MessageModeMessagePrivateFrom           MessageMode = 4
	MessageModeMessageChannel               MessageMode = 5
	MessageModeMessageChannelManagement     MessageMode = 6
	MessageModeMessageRVRAnswer             MessageMode = 7
	MessageModeMessageRVRContinue           MessageMode = 8
	MessageModeMessageGamemasterBroadcast   MessageMode = 9
	MessageModeMessageGamemasterChannel     MessageMode = 10
	MessageModeMessageGamemasterPrivateFrom MessageMode = 11
	MessageModeMessageChannelHighlight      MessageMode = 12
	MessageMode14                           MessageMode = 14 // To be checked.
	MessageModeMessageMonsterSay            MessageMode = 16
	MessageModeMessageMonsterYell           MessageMode = 17
	MessageModeMessageWarning               MessageMode = 18
	MessageModeMessageGame                  MessageMode = 19
	MessageModeMessageLogin                 MessageMode = 20
	MessageModeMessageStatus                MessageMode = 21
	MessageModeMessageLook                  MessageMode = 22
	MessageModeMessageFailure               MessageMode = 23
	MessageModeMessageBlue                  MessageMode = 24
	MessageModeMessageRed                   MessageMode = 25
	MessageModeMessageHealOthers            MessageMode = 26
	MessageModeMessageExpOthers             MessageMode = 27
	MessageModeMessageLoot                  MessageMode = 29
	MessageModeMessageTradeNpc              MessageMode = 30
	MessageModeMessageGuild                 MessageMode = 31
	MessageModeMessagePartyManagement       MessageMode = 32
	MessageModeMessageParty                 MessageMode = 33
	MessageModeMessageBarkLow               MessageMode = 34
	MessageModeMessageBarkLoud              MessageMode = 35
	MessageModeMessageReport                MessageMode = 36
	MessageModeMessageHotkeyUse             MessageMode = 37
	MessageModeMessageTutorialHint          MessageMode = 38
	MessageModeMessageThankyou              MessageMode = 39
	MessageModeMessageMarket                MessageMode = 40
	MessageModeMessageMana                  MessageMode = 41
	MessageModeMessageBeyondLast            MessageMode = 42
	MessageModeMessageGameHighlight         MessageMode = 50
	MessageModeMessageNpcFromStartBlock     MessageMode = 51
	MessageModeLastMessage                  MessageMode = 52
	MessageModeMessageInvalid               MessageMode = 255
)

var messageModeMap = map[MessageMode]string{
	MessageModeMessageNone:                  "MessageNone",
	MessageModeMessageSay:                   "MessageSay",
	MessageModeMessageWhisper:               "MessageWhisper",
	MessageModeMessageYell:                  "MessageYell",
	MessageModeMessagePrivateFrom:           "MessagePrivateFrom",
	MessageModeMessageChannel:               "MessageChannel",
	MessageModeMessageChannelManagement:     "MessageChannelManagement",
	MessageModeMessageRVRAnswer:             "MessageRVRAnswer",
	MessageModeMessageRVRContinue:           "MessageRVRContinue",
	MessageModeMessageGamemasterBroadcast:   "MessageGamemasterBroadcast",
	MessageModeMessageGamemasterChannel:     "MessageGamemasterChannel",
	MessageModeMessageGamemasterPrivateFrom: "MessageGamemasterPrivateFrom",
	MessageModeMessageChannelHighlight:      "MessageChannelHighlight",
	MessageMode14:                           "Weird-14",
	MessageModeMessageMonsterSay:            "MessageMonsterSay",
	MessageModeMessageMonsterYell:           "MessageMonsterYell",
	MessageModeMessageWarning:               "MessageWarning",
	MessageModeMessageGame:                  "MessageGame",
	MessageModeMessageLogin:                 "MessageLogin",
	MessageModeMessageStatus:                "MessageStatus",
	MessageModeMessageLook:                  "MessageLook",
	MessageModeMessageFailure:               "MessageFailure",
	MessageModeMessageBlue:                  "MessageBlue",
	MessageModeMessageRed:                   "MessageRed",
	MessageModeMessageHealOthers:            "MessageHealOthers",
	MessageModeMessageExpOthers:             "MessageExpOthers",
	MessageModeMessageLoot:                  "MessageLoot",
	MessageModeMessageTradeNpc:              "MessageTradeNpc",
	MessageModeMessageGuild:                 "MessageGuild",
	MessageModeMessagePartyManagement:       "MessagePartyManagement",
	MessageModeMessageParty:                 "MessageParty",
	MessageModeMessageBarkLow:               "MessageBarkLow",
	MessageModeMessageBarkLoud:              "MessageBarkLoud",
	MessageModeMessageReport:                "MessageReport",
	MessageModeMessageHotkeyUse:             "MessageHotkeyUse",
	MessageModeMessageTutorialHint:          "MessageTutorialHint",
	MessageModeMessageThankyou:              "MessageThankyou",
	MessageModeMessageMarket:                "MessageMarket",
	MessageModeMessageMana:                  "MessageMana",
	MessageModeMessageBeyondLast:            "MessageBeyondLast",
	MessageModeMessageGameHighlight:         "MessageGameHighlight",
	MessageModeMessageNpcFromStartBlock:     "MessageNpcFromStartBlock",
	MessageModeLastMessage:                  "LastMessage",
	MessageModeMessageInvalid:               "MessageInvalid",
}
