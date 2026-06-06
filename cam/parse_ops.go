package cam

import (
	"fmt"
	"time"

	"github.com/s5i/tcam/data"
)

func parsePacket(state *parseState, buf []byte, timeOffset time.Duration, opts *ParseOpts) ([]data.Operation, error) {
	m := newMessage(buf)
	var ops []data.Operation

	for m.remaining() > 0 {
		head, err := m.getByte()
		if err != nil {
			return ops, fmt.Errorf("reading packet head: %w", err)
		}

		t := time.Now()
		opcode := data.OpType(head)
		ignore := opts.TFilter != nil && !opts.TFilter[opcode]
		f, ok := parseFunc[opcode]

		if state.stats != nil {
			state.stats.Count[opcode]++
			state.stats.Duration[opcode] += time.Since(t)
		}

		if !ok {
			return ops, fmt.Errorf("unknown packet head: 0x%02X", head)
		}

		op, err := f(m, state, ignore, timeOffset)
		if err != nil {
			return ops, fmt.Errorf("parsing 0x%02X: %w", head, err)
		}

		if !ignore {
			ops = append(ops, op)
		}
	}

	return ops, nil

}

var parseFunc = map[data.OpType]func(*message, *parseState, bool, time.Duration) (data.Operation, error){
	data.TLoginPlayerState:      parseLoginPlayerState,
	data.TLoginError:            parseLoginError,
	data.TLoginWaitList:         parseLoginWaitList,
	data.TPing:                  parsePing,
	data.TMap:                   parseMap,
	data.TMoveNorth:             parseMoveNorth,
	data.TMoveEast:              parseMoveEast,
	data.TMoveSouth:             parseMoveSouth,
	data.TMoveWest:              parseMoveWest,
	data.TTileUpdate:            parseTileUpdate,
	data.TTileItemAdd:           parseTileItemAdd,
	data.TTileItemUpdate:        parseTileItemUpdate,
	data.TTileItemRemove:        parseTileItemRemove,
	data.TCreatureMove:          parseCreatureMove,
	data.TContainerOpen:         parseContainerOpen,
	data.TContainerClose:        parseContainerClose,
	data.TContainerItemAdd:      parseContainerItemAdd,
	data.TContainerItemUpdate:   parseContainerItemUpdate,
	data.TContainerItemRemove:   parseContainerItemRemove,
	data.TInventoryItemSet:      parseInventoryItemSet,
	data.TInventoryItemClear:    parseInventoryItemClear,
	data.TTradeOwn:              parseTradeOwn,
	data.TTradeCounter:          parseTradeCounter,
	data.TTradeClose:            parseTradeClose,
	data.TEffectLight:           parseEffectLight,
	data.TEffectGraphical:       parseEffectGraphical,
	data.TEffectText:            parseEffectText,
	data.TEffectMissile:         parseEffectMissile,
	data.TCreatureSquare:        parseCreatureSquare,
	data.TCreatureHealth:        parseCreatureHealth,
	data.TCreatureLight:         parseCreatureLight,
	data.TCreatureOutfit:        parseCreatureOutfit,
	data.TCreatureSpeed:         parseCreatureSpeed,
	data.TCreatureSkull:         parseCreatureSkull,
	data.TCreatureParty:         parseCreatureParty,
	data.TPromptTextUpdate:      parsePromptTextUpdate,
	data.TPromptHouseList:       parsePromptHouseList,
	data.TPlayerStats:           parsePlayerStats,
	data.TPlayerSkills:          parsePlayerSkills,
	data.TPlayerIcons:           parsePlayerIcons,
	data.TTargetClear:           parseTargetClear,
	data.TCreatureMessage:       parseCreatureMessage,
	data.TChannelList:           parseChannelList,
	data.TChannelOpen:           parseChannelOpen,
	data.TPrivateChannelOpen:    parsePrivateChannelOpen,
	data.TRuleViolationsChannel: parseRuleViolationsChannel,
	data.TRuleViolationsRemove:  parseRuleViolationsRemove,
	data.TRuleViolationCancel:   parseRuleViolationCancel,
	data.TRuleViolationsLock:    parseRuleViolationsLock,
	data.TPrivateChannelCreate:  parsePrivateChannelCreate,
	data.TPrivateChannelClose:   parsePrivateChannelClose,
	data.TMessage:               parseMessage,
	data.TMoveCancel:            parseMoveCancel,
	data.TMoveFloorUp:           parseMoveFloorUp,
	data.TMoveFloorDown:         parseMoveFloorDown,
	data.TPromptChooseOutfit:    parsePromptChooseOutfit,
	data.TVIPState:              parseVIPState,
	data.TVIPLogin:              parseVIPLogin,
	data.TVIPLogout:             parseVIPLogout,
}

func parseLoginPlayerState(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	if s.playerName == "" {
		ignore = false
	}

	op := data.LoginPlayerState{TimeOffset: offset, PlayerPos: s.playerPos}
	var err error
	op.PlayerID, err = m.getU32()
	if err != nil {
		return nil, err
	}
	if _, err := m.getByte(); err != nil { // unknown
		return nil, err
	}
	if _, err := m.getByte(); err != nil { // unknown
		return nil, err
	}
	op.AccessLevel, err = m.getByte()
	if err != nil {
		return nil, err
	}
	if op.AccessLevel == 1 {
		if _, err := m.getByte(); err != nil { // "loop" byte (unused)
			return nil, err
		}
		for range 32 {
			if _, err := m.getByte(); err != nil {
				return nil, err
			}
		}
	}
	s.playerID = op.PlayerID
	return op, nil
}

func parseLoginError(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	var msg string
	if err := m.getString(&msg, ignore); err != nil {
		return nil, err
	}
	return data.LoginError{TimeOffset: offset, PlayerPos: s.playerPos, Message: msg}, nil
}

func parseLoginWaitList(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	var msg string
	if err := m.getString(&msg, ignore); err != nil {
		return nil, err
	}
	t, err := m.getByte()
	if err != nil {
		return nil, err
	}
	return data.LoginWaitList{TimeOffset: offset, PlayerPos: s.playerPos, Message: msg, Time: t}, nil
}

func parsePing(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	return data.Ping{TimeOffset: offset, PlayerPos: s.playerPos}, nil
}

func parseMap(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	if s.playerName == "" {
		ignore = false
	}

	loc, err := m.getLocation()
	if err != nil {
		return nil, err
	}
	s.playerPos = loc
	tiles, err := m.getMapDescription(ignore, loc.X-8, loc.Y-6, loc.Z, 18, 14)
	if err != nil {
		return nil, err
	}
	s.resolvePlayerName(tiles)
	return data.Map{TimeOffset: offset, PlayerPos: loc, Tiles: tiles}, nil
}

func parseMoveNorth(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	s.playerPos.Y--
	loc := s.playerPos
	tiles, err := m.getMapDescription(ignore, loc.X-8, loc.Y-6, loc.Z, 18, 1)
	if err != nil {
		return nil, err
	}
	return data.MoveNorth{TimeOffset: offset, PlayerPos: s.playerPos, Tiles: tiles}, nil
}

func parseMoveEast(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	s.playerPos.X++
	loc := s.playerPos
	tiles, err := m.getMapDescription(ignore, loc.X+9, loc.Y-6, loc.Z, 1, 14)
	if err != nil {
		return nil, err
	}
	return data.MoveEast{TimeOffset: offset, PlayerPos: s.playerPos, Tiles: tiles}, nil
}

func parseMoveSouth(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	s.playerPos.Y++
	loc := s.playerPos
	tiles, err := m.getMapDescription(ignore, loc.X-8, loc.Y+7, loc.Z, 18, 1)
	if err != nil {
		return nil, err
	}
	return data.MoveSouth{TimeOffset: offset, PlayerPos: s.playerPos, Tiles: tiles}, nil
}

func parseMoveWest(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	s.playerPos.X--
	loc := s.playerPos
	tiles, err := m.getMapDescription(ignore, loc.X-8, loc.Y-6, loc.Z, 1, 14)
	if err != nil {
		return nil, err
	}
	return data.MoveWest{TimeOffset: offset, PlayerPos: s.playerPos, Tiles: tiles}, nil
}

func parseTileUpdate(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	loc, err := m.getLocation()
	if err != nil {
		return nil, err
	}
	thingID, err := m.peekU16()
	if err != nil {
		return nil, err
	}
	if thingID == 0xFF01 {
		if _, err := m.getU16(); err != nil {
			return nil, err
		}
		return data.TileUpdate{TimeOffset: offset, PlayerPos: s.playerPos, Location: loc}, nil
	}
	tile, err := m.parseTileDescription(ignore, loc)
	if err != nil {
		return nil, err
	}
	if _, err := m.getU16(); err != nil {
		return nil, err
	}
	return data.TileUpdate{TimeOffset: offset, PlayerPos: s.playerPos, Location: loc, Tile: tile, HasTile: true}, nil
}

func parseTileItemAdd(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	loc, err := m.getLocation()
	if err != nil {
		return nil, err
	}
	thing, err := m.getThing(ignore)
	if err != nil {
		return nil, err
	}
	return data.TileItemAdd{TimeOffset: offset, PlayerPos: s.playerPos, Location: loc, Thing: thing}, nil
}

func parseTileItemUpdate(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	loc, err := m.getLocation()
	if err != nil {
		return nil, err
	}
	stack, err := m.getByte()
	if err != nil {
		return nil, err
	}
	thing, err := m.getThing(ignore)
	if err != nil {
		return nil, err
	}
	return data.TileItemUpdate{TimeOffset: offset, PlayerPos: s.playerPos, Location: loc, StackIndex: stack, Thing: thing}, nil
}

func parseTileItemRemove(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	loc, err := m.getLocation()
	if err != nil {
		return nil, err
	}
	stack, err := m.getByte()
	if err != nil {
		return nil, err
	}
	return data.TileItemRemove{TimeOffset: offset, PlayerPos: s.playerPos, Location: loc, StackIndex: stack}, nil
}

func parseCreatureMove(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	oldLoc, err := m.getLocation()
	if err != nil {
		return nil, err
	}
	oldStack, err := m.getByte()
	if err != nil {
		return nil, err
	}
	newLoc, err := m.getLocation()
	if err != nil {
		return nil, err
	}

	return data.CreatureMove{TimeOffset: offset, PlayerPos: s.playerPos, OldLocation: oldLoc, OldStack: oldStack, NewLocation: newLoc}, nil
}

func parseContainerOpen(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	op := data.ContainerOpen{TimeOffset: offset, PlayerPos: s.playerPos}
	var err error
	op.ContainerID, err = m.getByte()
	if err != nil {
		return nil, err
	}
	op.ItemID, err = m.getU16()
	if err != nil {
		return nil, err
	}
	err = m.getString(&op.Name, ignore)
	if err != nil {
		return nil, err
	}
	op.Volume, err = m.getByte()
	if err != nil {
		return nil, err
	}
	op.HasParent, err = m.getByte()
	if err != nil {
		return nil, err
	}
	size, err := m.getByte()
	if err != nil {
		return nil, err
	}
	for i := 0; i < int(size); i++ {
		thing, err := m.getThing(ignore)
		if err != nil {
			return nil, err
		}
		op.Items = append(op.Items, thing)
	}
	return op, nil
}

func parseContainerClose(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	id, err := m.getByte()
	if err != nil {
		return nil, err
	}
	return data.ContainerClose{TimeOffset: offset, PlayerPos: s.playerPos, ContainerID: id}, nil
}

func parseContainerItemAdd(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	id, err := m.getByte()
	if err != nil {
		return nil, err
	}
	thing, err := m.getThing(ignore)
	if err != nil {
		return nil, err
	}
	return data.ContainerItemAdd{TimeOffset: offset, PlayerPos: s.playerPos, ContainerID: id, Thing: thing}, nil
}

func parseContainerItemUpdate(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	id, err := m.getByte()
	if err != nil {
		return nil, err
	}
	slot, err := m.getByte()
	if err != nil {
		return nil, err
	}
	thing, err := m.getThing(ignore)
	if err != nil {
		return nil, err
	}
	return data.ContainerItemUpdate{TimeOffset: offset, PlayerPos: s.playerPos, ContainerID: id, Slot: slot, Thing: thing}, nil
}

func parseContainerItemRemove(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	id, err := m.getByte()
	if err != nil {
		return nil, err
	}
	slot, err := m.getByte()
	if err != nil {
		return nil, err
	}
	return data.ContainerItemRemove{TimeOffset: offset, PlayerPos: s.playerPos, ContainerID: id, Slot: slot}, nil
}

func parseInventoryItemSet(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	slot, err := m.getByte()
	if err != nil {
		return nil, err
	}
	item, err := getItem(m, 0xFFFF)
	if err != nil {
		return nil, err
	}
	return data.InventoryItemSet{TimeOffset: offset, PlayerPos: s.playerPos, Slot: slot, Item: item}, nil

}

func parseInventoryItemClear(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	slot, err := m.getByte()
	if err != nil {
		return nil, err
	}
	return data.InventoryItemClear{TimeOffset: offset, PlayerPos: s.playerPos, Slot: slot}, nil
}

func parseTradeOwn(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	var name string
	if err := m.getString(&name, ignore); err != nil {
		return nil, err
	}
	size, err := m.getByte()
	if err != nil {
		return nil, err
	}
	op := data.TradeOwn{TimeOffset: offset, PlayerPos: s.playerPos, Name: name}
	for i := 0; i < int(size); i++ {
		thing, err := m.getThing(ignore)
		if err != nil {
			return nil, err
		}
		op.Items = append(op.Items, thing)
	}
	return op, nil
}

func parseTradeCounter(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	var name string
	if err := m.getString(&name, ignore); err != nil {
		return nil, err
	}
	size, err := m.getByte()
	if err != nil {
		return nil, err
	}
	op := data.TradeCounter{TimeOffset: offset, PlayerPos: s.playerPos, Name: name}
	for i := 0; i < int(size); i++ {
		thing, err := m.getThing(ignore)
		if err != nil {
			return nil, err
		}
		op.Items = append(op.Items, thing)
	}
	return op, nil
}

func parseTradeClose(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	return data.TradeClose{TimeOffset: offset, PlayerPos: s.playerPos}, nil
}

func parseEffectLight(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	level, err := m.getByte()
	if err != nil {
		return nil, err
	}
	color, err := m.getByte()
	if err != nil {
		return nil, err
	}
	return data.EffectLight{TimeOffset: offset, PlayerPos: s.playerPos, Level: level, Color: color}, nil
}

func parseEffectGraphical(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	loc, err := m.getLocation()
	if err != nil {
		return nil, err
	}
	effect, err := m.getByte()
	if err != nil {
		return nil, err
	}
	return data.EffectGraphical{TimeOffset: offset, PlayerPos: s.playerPos, Location: loc, Effect: effect}, nil
}

func parseEffectText(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	loc, err := m.getLocation()
	if err != nil {
		return nil, err
	}
	color, err := m.getByte()
	if err != nil {
		return nil, err
	}
	var text string
	err = m.getString(&text, ignore)
	if err != nil {
		return nil, err
	}
	return data.EffectText{TimeOffset: offset, PlayerPos: s.playerPos, Location: loc, Color: color, Text: text}, nil
}

func parseEffectMissile(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	from, err := m.getLocation()
	if err != nil {
		return nil, err
	}
	to, err := m.getLocation()
	if err != nil {
		return nil, err
	}
	effect, err := m.getByte()
	if err != nil {
		return nil, err
	}
	return data.EffectMissile{TimeOffset: offset, PlayerPos: s.playerPos, From: from, To: to, Effect: effect}, nil
}

func parseCreatureSquare(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	id, err := m.getU32()
	if err != nil {
		return nil, err
	}
	color, err := m.getByte()
	if err != nil {
		return nil, err
	}
	return data.CreatureSquare{TimeOffset: offset, PlayerPos: s.playerPos, CreatureID: id, Color: color}, nil
}

func parseCreatureHealth(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	id, err := m.getU32()
	if err != nil {
		return nil, err
	}
	health, err := m.getByte()
	if err != nil {
		return nil, err
	}
	return data.CreatureHealth{TimeOffset: offset, PlayerPos: s.playerPos, CreatureID: id, Health: health}, nil
}

func parseCreatureLight(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	id, err := m.getU32()
	if err != nil {
		return nil, err
	}
	level, err := m.getByte()
	if err != nil {
		return nil, err
	}
	color, err := m.getByte()
	if err != nil {
		return nil, err
	}
	return data.CreatureLight{TimeOffset: offset, PlayerPos: s.playerPos, CreatureID: id, Level: level, Color: color}, nil
}

func parseCreatureOutfit(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	id, err := m.getU32()
	if err != nil {
		return nil, err
	}
	outfit, err := m.getOutfit()
	if err != nil {
		return nil, err
	}
	return data.CreatureOutfit{TimeOffset: offset, PlayerPos: s.playerPos, CreatureID: id, Outfit: outfit}, nil
}

func parseCreatureSpeed(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	id, err := m.getU32()
	if err != nil {
		return nil, err
	}
	speed, err := m.getU16()
	if err != nil {
		return nil, err
	}
	return data.CreatureSpeed{TimeOffset: offset, PlayerPos: s.playerPos, CreatureID: id, Speed: speed}, nil
}

func parseCreatureSkull(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	id, err := m.getU32()
	if err != nil {
		return nil, err
	}
	skull, err := m.getByte()
	if err != nil {
		return nil, err
	}
	return data.CreatureSkull{TimeOffset: offset, PlayerPos: s.playerPos, CreatureID: id, Skull: skull}, nil
}

func parseCreatureParty(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	id, err := m.getU32()
	if err != nil {
		return nil, err
	}
	shield, err := m.getByte()
	if err != nil {
		return nil, err
	}
	return data.CreatureParty{TimeOffset: offset, PlayerPos: s.playerPos, CreatureID: id, Shield: shield}, nil
}

func parsePromptTextUpdate(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	op := data.PromptTextUpdate{TimeOffset: offset, PlayerPos: s.playerPos}
	var err error
	op.WindowID, err = m.getU32()
	if err != nil {
		return nil, err
	}
	op.ItemID, err = m.getU16()
	if err != nil {
		return nil, err
	}
	op.MaxLen, err = m.getU16()
	if err != nil {
		return nil, err
	}
	err = m.getString(&op.Text, ignore)
	if err != nil {
		return nil, err
	}
	err = m.getString(&op.Author, ignore)
	if err != nil {
		return nil, err
	}
	return op, nil
}

func parsePromptHouseList(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	op := data.PromptHouseList{TimeOffset: offset, PlayerPos: s.playerPos}
	var err error
	op.Unknown, err = m.getByte()
	if err != nil {
		return nil, err
	}
	op.ID, err = m.getU32()
	if err != nil {
		return nil, err
	}
	err = m.getString(&op.Text, ignore)
	if err != nil {
		return nil, err
	}
	return op, nil
}

func parsePlayerStats(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	op := data.PlayerStats{TimeOffset: offset, PlayerPos: s.playerPos}
	var err error
	op.HP, err = m.getU16()
	if err != nil {
		return nil, err
	}
	op.MaxHP, err = m.getU16()
	if err != nil {
		return nil, err
	}
	op.Capacity, err = m.getU16()
	if err != nil {
		return nil, err
	}
	op.Exp, err = m.getU32()
	if err != nil {
		return nil, err
	}
	op.Level, err = m.getByte()
	if err != nil {
		return nil, err
	}
	op.LevelPct, err = m.getByte()
	if err != nil {
		return nil, err
	}
	op.Mana, err = m.getU16()
	if err != nil {
		return nil, err
	}
	op.MaxMana, err = m.getU16()
	if err != nil {
		return nil, err
	}
	op.MagicLvl, err = m.getByte()
	if err != nil {
		return nil, err
	}
	op.MagicLvlPct, err = m.getByte()
	if err != nil {
		return nil, err
	}
	op.Soul, err = m.getU16()
	if err != nil {
		return nil, err
	}
	return op, nil
}

func parsePlayerSkills(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	op := data.PlayerSkills{TimeOffset: offset, PlayerPos: s.playerPos}
	for i := 0; i < 7; i++ {
		level, err := m.getByte()
		if err != nil {
			return nil, err
		}
		pct, err := m.getByte()
		if err != nil {
			return nil, err
		}
		op.Skills[i] = data.SkillValue{Level: level, Percent: pct}
	}
	return op, nil
}

func parsePlayerIcons(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	icons, err := m.getByte()
	if err != nil {
		return nil, err
	}
	return data.PlayerIcons{TimeOffset: offset, PlayerPos: s.playerPos, Icons: icons}, nil
}

func parseTargetClear(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	return data.TargetClear{TimeOffset: offset, PlayerPos: s.playerPos}, nil
}

func parseCreatureMessage(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	op := data.CreatureMessage{TimeOffset: offset, PlayerPos: s.playerPos}
	var err error
	op.StatementID, err = m.getU32()
	if err != nil {
		return nil, err
	}
	err = m.getString(&op.Name, ignore)
	if err != nil {
		return nil, err
	}
	op.Type, err = m.getByte()
	if err != nil {
		return nil, err
	}
	switch op.Type {
	case 1, 2, 3, 0x10, 0x11:
		loc, err := m.getLocation()
		if err != nil {
			return nil, err
		}
		op.Location = &loc
	case 5, 6, 0xA, 0xC, 0xE:
		ch, err := m.getU16()
		if err != nil {
			return nil, err
		}
		op.ChannelID = &ch
	}
	err = m.getString(&op.Text, ignore)
	if err != nil {
		return nil, err
	}
	return op, nil
}

func parseChannelList(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	size, err := m.getByte()
	if err != nil {
		return nil, err
	}
	op := data.ChannelList{
		TimeOffset: offset,
		PlayerPos:  s.playerPos,
	}

	if !ignore {
		op.Channels = make([]data.ChannelEntry, 0, size)
	}

	for range size {
		id, err := m.getU16()
		if err != nil {
			return nil, err
		}
		var name string
		err = m.getString(&name, ignore)
		if err != nil {
			return nil, err
		}
		if !ignore {
			op.Channels = append(op.Channels, data.ChannelEntry{ID: id, Name: name})
		}
	}
	return op, nil
}

func parseChannelOpen(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	id, err := m.getU16()
	if err != nil {
		return nil, err
	}
	var name string
	err = m.getString(&name, ignore)
	if err != nil {
		return nil, err
	}
	return data.ChannelOpen{TimeOffset: offset, PlayerPos: s.playerPos, ID: id, Name: name}, nil
}

func parsePrivateChannelOpen(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	var name string
	if err := m.getString(&name, ignore); err != nil {
		return nil, err
	}
	return data.PrivateChannelOpen{TimeOffset: offset, PlayerPos: s.playerPos, Name: name}, nil
}

func parseRuleViolationsChannel(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	size, err := m.getU16()
	if err != nil {
		return nil, err
	}
	return data.RuleViolationsChannel{TimeOffset: offset, PlayerPos: s.playerPos, Size: size}, nil
}

func parseRuleViolationsRemove(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	var name string
	if err := m.getString(&name, ignore); err != nil {
		return nil, err
	}
	return data.RuleViolationsRemove{TimeOffset: offset, PlayerPos: s.playerPos, Name: name}, nil
}

func parseRuleViolationCancel(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	var name string
	if err := m.getString(&name, ignore); err != nil {
		return nil, err
	}
	return data.RuleViolationCancel{TimeOffset: offset, PlayerPos: s.playerPos, Name: name}, nil
}

func parseRuleViolationsLock(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	return data.RuleViolationsLock{TimeOffset: offset, PlayerPos: s.playerPos}, nil
}

func parsePrivateChannelCreate(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	id, err := m.getU16()
	if err != nil {
		return nil, err
	}
	var name string
	if err := m.getString(&name, ignore); err != nil {
		return nil, err
	}
	return data.PrivateChannelCreate{TimeOffset: offset, PlayerPos: s.playerPos, ID: id, Name: name}, nil
}

func parsePrivateChannelClose(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	id, err := m.getU16()
	if err != nil {
		return nil, err
	}
	return data.PrivateChannelClose{TimeOffset: offset, PlayerPos: s.playerPos, ChannelID: id}, nil
}

func parseMessage(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	t, err := m.getByte()
	if err != nil {
		return nil, err
	}
	var text string
	if err := m.getString(&text, ignore); err != nil {
		return nil, err
	}
	return data.Message{TimeOffset: offset, PlayerPos: s.playerPos, Type: t, Text: text}, nil
}

func parseMoveCancel(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	dir, err := m.getByte()
	if err != nil {
		return nil, err
	}
	return data.MoveCancel{TimeOffset: offset, PlayerPos: s.playerPos, Direction: data.Direction(dir)}, nil
}

func parseMoveFloorUp(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	myPos := s.playerPos
	myPos.Z--

	var tiles []data.Tile
	if myPos.Z == 7 {
		skip := 0
		for _, floor := range []struct{ z, offset int }{{5, 3}, {4, 4}, {3, 5}, {2, 6}, {1, 7}, {0, 8}} {
			t, err := m.parseFloorDescription(ignore, myPos.X-8, myPos.Y-6, floor.z, 18, 14, floor.offset, &skip)
			if err != nil {
				return nil, err
			}
			if !ignore {
				tiles = append(tiles, t...)
			}

		}
	} else if myPos.Z > 7 {
		skip := 0
		t, err := m.parseFloorDescription(ignore, myPos.X-8, myPos.Y-6, myPos.Z-2, 18, 14, 3, &skip)
		if err != nil {
			return nil, err
		}
		if !ignore {
			tiles = append(tiles, t...)
		}
	}

	s.playerPos = data.Location{X: myPos.X + 1, Y: myPos.Y + 1, Z: myPos.Z}
	return data.MoveFloorUp{TimeOffset: offset, PlayerPos: s.playerPos, Tiles: tiles}, nil
}

func parseMoveFloorDown(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	myPos := s.playerPos
	myPos.Z++

	var tiles []data.Tile
	skip := 0
	if myPos.Z == 8 {
		for i, j := myPos.Z, -1; i < myPos.Z+3; i, j = i+1, j-1 {
			t, err := m.parseFloorDescription(ignore, myPos.X-8, myPos.Y-6, i, 18, 14, j, &skip)
			if err != nil {
				return nil, err
			}
			tiles = append(tiles, t...)
		}
	} else if myPos.Z > 8 && myPos.Z < 14 {
		t, err := m.parseFloorDescription(ignore, myPos.X-8, myPos.Y-6, myPos.Z+2, 18, 14, -3, &skip)
		if err != nil {
			return nil, err
		}
		tiles = append(tiles, t...)
	}

	s.playerPos = data.Location{X: myPos.X - 1, Y: myPos.Y - 1, Z: myPos.Z}
	return data.MoveFloorDown{TimeOffset: offset, PlayerPos: s.playerPos, Tiles: tiles}, nil
}

func parsePromptChooseOutfit(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	outfit, err := m.getOutfit()
	if err != nil {
		return nil, err
	}
	start, err := m.getU16()
	if err != nil {
		return nil, err
	}
	end, err := m.getU16()
	if err != nil {
		return nil, err
	}
	return data.PromptChooseOutfit{TimeOffset: offset, PlayerPos: s.playerPos, Outfit: outfit, OutfitStart: start, OutfitEnd: end}, nil
}

func parseVIPState(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	id, err := m.getU32()
	if err != nil {
		return nil, err
	}
	var name string
	err = m.getString(&name, ignore)
	if err != nil {
		return nil, err
	}
	online, err := m.getByte()
	if err != nil {
		return nil, err
	}
	return data.VIPState{TimeOffset: offset, PlayerPos: s.playerPos, ID: id, Name: name, Online: online}, nil
}

func parseVIPLogin(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	id, err := m.getU32()
	if err != nil {
		return nil, err
	}
	return data.VIPLogin{TimeOffset: offset, PlayerPos: s.playerPos, ID: id}, nil
}

func parseVIPLogout(m *message, s *parseState, ignore bool, offset time.Duration) (data.Operation, error) {
	id, err := m.getU32()
	if err != nil {
		return nil, err
	}
	return data.VIPLogout{TimeOffset: offset, PlayerPos: s.playerPos, ID: id}, nil
}
