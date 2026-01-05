package enum

import "fmt"

type Item uint16

func (c Item) String() string {
	if str, ok := itemMap[c]; ok {
		return str
	}
	return fmt.Sprintf("Unknown-%d", c)
}

const (
	ItemInvalid          Item = 0
	ItemUnknownCreature  Item = 97
	ItemOutdatedCreature Item = 98
	ItemCreature         Item = 99
)

var itemMap = map[Item]string{
	ItemInvalid:          "Invalid",
	ItemUnknownCreature:  "UnknownCreature",
	ItemOutdatedCreature: "OutdatedCreature",
	ItemCreature:         "Creature",
}
