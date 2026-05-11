package cam

import (
	"time"

	"github.com/s5i/tcam/data"
)

func parseLogin(m *message, ignore bool, offset time.Duration) (data.Login, error) {
	op := data.Login{TimeOffset: offset}
	var err error
	op.PlayerID, err = m.getU32()
	if err != nil {
		return op, err
	}
	if _, err := m.getByte(); err != nil { // unknown
		return op, err
	}
	if _, err := m.getByte(); err != nil { // unknown
		return op, err
	}
	op.AccessLevel, err = m.getByte()
	if err != nil {
		return op, err
	}
	if op.AccessLevel == 1 {
		if _, err := m.getByte(); err != nil { // "loop" byte (unused)
			return op, err
		}
		for i := 0; i < 32; i++ {
			if _, err := m.getByte(); err != nil {
				return op, err
			}
		}
	}
	return op, nil
}

func parseDisconnectClient(m *message, ignore bool, offset time.Duration) (data.DisconnectClient, error) {
	msg, err := m.getString()
	return data.DisconnectClient{TimeOffset: offset, Message: msg}, err
}

func parseWaitList(m *message, ignore bool, offset time.Duration) (data.WaitList, error) {
	msg, err := m.getString()
	if err != nil {
		return data.WaitList{}, err
	}
	t, err := m.getByte()
	return data.WaitList{TimeOffset: offset, Message: msg, Time: t}, err
}

func parsePing(m *message, ignore bool, offset time.Duration) (data.Ping, error) {
	return data.Ping{TimeOffset: offset}, nil
}

func parseMapDescription(m *message, s *parseState, ignore bool, offset time.Duration) (data.MapDescription, error) {
	loc, err := m.getLocation()
	if err != nil {
		return data.MapDescription{}, err
	}
	s.playerPos = loc
	tiles, err := getMapDescription(m, loc.X-8, loc.Y-6, loc.Z, 18, 14)
	if err != nil {
		return data.MapDescription{}, err
	}
	s.updateTiles(tiles)
	return data.MapDescription{TimeOffset: offset, PlayerPos: loc, Tiles: tiles}, nil
}

func parseMoveNorth(m *message, s *parseState, ignore bool, offset time.Duration) (data.MoveNorth, error) {
	s.playerPos.Y--
	loc := s.playerPos
	tiles, err := getMapDescription(m, loc.X-8, loc.Y-6, loc.Z, 18, 1)
	if err != nil {
		return data.MoveNorth{}, err
	}
	s.updateTiles(tiles)
	return data.MoveNorth{TimeOffset: offset, Tiles: tiles}, nil
}

func parseMoveEast(m *message, s *parseState, ignore bool, offset time.Duration) (data.MoveEast, error) {
	s.playerPos.X++
	loc := s.playerPos
	tiles, err := getMapDescription(m, loc.X+9, loc.Y-6, loc.Z, 1, 14)
	if err != nil {
		return data.MoveEast{}, err
	}
	s.updateTiles(tiles)
	return data.MoveEast{TimeOffset: offset, Tiles: tiles}, nil
}

func parseMoveSouth(m *message, s *parseState, ignore bool, offset time.Duration) (data.MoveSouth, error) {
	s.playerPos.Y++
	loc := s.playerPos
	tiles, err := getMapDescription(m, loc.X-8, loc.Y+7, loc.Z, 18, 1)
	if err != nil {
		return data.MoveSouth{}, err
	}
	s.updateTiles(tiles)
	return data.MoveSouth{TimeOffset: offset, Tiles: tiles}, nil
}

func parseMoveWest(m *message, s *parseState, ignore bool, offset time.Duration) (data.MoveWest, error) {
	s.playerPos.X--
	loc := s.playerPos
	tiles, err := getMapDescription(m, loc.X-8, loc.Y-6, loc.Z, 1, 14)
	if err != nil {
		return data.MoveWest{}, err
	}
	s.updateTiles(tiles)
	return data.MoveWest{TimeOffset: offset, Tiles: tiles}, nil
}

func parseUpdateTile(m *message, s *parseState, ignore bool, offset time.Duration) (data.UpdateTile, error) {
	loc, err := m.getLocation()
	if err != nil {
		return data.UpdateTile{}, err
	}
	thingID, err := m.peekU16()
	if err != nil {
		return data.UpdateTile{}, err
	}
	if thingID == 0xFF01 {
		if _, err := m.getU16(); err != nil {
			return data.UpdateTile{}, err
		}
		return data.UpdateTile{TimeOffset: offset, Location: loc}, nil
	}
	tile, err := parseTileDescription(m, loc)
	if err != nil {
		return data.UpdateTile{}, err
	}
	if _, err := m.getU16(); err != nil {
		return data.UpdateTile{}, err
	}
	s.updateTiles([]data.Tile{tile})
	return data.UpdateTile{TimeOffset: offset, Location: loc, Tile: &tile}, nil
}

func parseAddTileItem(m *message, s *parseState, ignore bool, offset time.Duration) (data.AddTileItem, error) {
	loc, err := m.getLocation()
	if err != nil {
		return data.AddTileItem{}, err
	}
	thing, err := getThing(m)
	if err != nil {
		return data.AddTileItem{}, err
	}
	s.addThing(loc, thing)
	return data.AddTileItem{TimeOffset: offset, Location: loc, Thing: thing}, nil
}

func parseUpdateTileItem(m *message, s *parseState, ignore bool, offset time.Duration) (data.UpdateTileItem, error) {
	loc, err := m.getLocation()
	if err != nil {
		return data.UpdateTileItem{}, err
	}
	stack, err := m.getByte()
	if err != nil {
		return data.UpdateTileItem{}, err
	}
	thing, err := getThing(m)
	if err != nil {
		return data.UpdateTileItem{}, err
	}
	s.replaceThing(loc, int(stack), thing)
	return data.UpdateTileItem{TimeOffset: offset, Location: loc, StackIndex: stack, Thing: thing}, nil
}

func parseRemoveTileItem(m *message, s *parseState, ignore bool, offset time.Duration) (data.RemoveTileItem, error) {
	loc, err := m.getLocation()
	if err != nil {
		return data.RemoveTileItem{}, err
	}
	stack, err := m.getByte()
	if err != nil {
		return data.RemoveTileItem{}, err
	}
	if !loc.IsCreature() {
		s.removeThing(loc, int(stack))
	}
	return data.RemoveTileItem{TimeOffset: offset, Location: loc, StackIndex: stack}, nil
}

func parseMoveCreature(m *message, s *parseState, ignore bool, offset time.Duration) (data.MoveCreature, error) {
	oldLoc, err := m.getLocation()
	if err != nil {
		return data.MoveCreature{}, err
	}
	oldStack, err := m.getByte()
	if err != nil {
		return data.MoveCreature{}, err
	}
	newLoc, err := m.getLocation()
	if err != nil {
		return data.MoveCreature{}, err
	}

	if oldLoc.IsCreature() {
		cID := oldLoc.CreatureID(oldStack)
		s.addThing(newLoc, data.Thing{Creature: &data.Creature{ID: cID}})
	} else {
		thing := s.getThing(oldLoc, int(oldStack))
		if thing != nil && thing.Creature != nil {
			s.removeThing(oldLoc, int(oldStack))
			s.addThing(newLoc, *thing)
		}
	}
	return data.MoveCreature{TimeOffset: offset, OldLocation: oldLoc, OldStack: oldStack, NewLocation: newLoc}, nil
}

func parseContainer(m *message, ignore bool, offset time.Duration) (data.Container, error) {
	op := data.Container{TimeOffset: offset}
	var err error
	op.ContainerID, err = m.getByte()
	if err != nil {
		return op, err
	}
	op.ItemID, err = m.getU16()
	if err != nil {
		return op, err
	}
	op.Name, err = m.getString()
	if err != nil {
		return op, err
	}
	op.Volume, err = m.getByte()
	if err != nil {
		return op, err
	}
	op.HasParent, err = m.getByte()
	if err != nil {
		return op, err
	}
	size, err := m.getByte()
	if err != nil {
		return op, err
	}
	for i := 0; i < int(size); i++ {
		thing, err := getThing(m)
		if err != nil {
			return op, err
		}
		op.Items = append(op.Items, thing)
	}
	return op, nil
}

func parseCloseContainer(m *message, ignore bool, offset time.Duration) (data.CloseContainer, error) {
	id, err := m.getByte()
	return data.CloseContainer{TimeOffset: offset, ContainerID: id}, err
}

func parseAddContainerItem(m *message, ignore bool, offset time.Duration) (data.AddContainerItem, error) {
	id, err := m.getByte()
	if err != nil {
		return data.AddContainerItem{}, err
	}
	thing, err := getThing(m)
	return data.AddContainerItem{TimeOffset: offset, ContainerID: id, Thing: thing}, err
}

func parseUpdateContainerItem(m *message, ignore bool, offset time.Duration) (data.UpdateContainerItem, error) {
	id, err := m.getByte()
	if err != nil {
		return data.UpdateContainerItem{}, err
	}
	slot, err := m.getByte()
	if err != nil {
		return data.UpdateContainerItem{}, err
	}
	thing, err := getThing(m)
	return data.UpdateContainerItem{TimeOffset: offset, ContainerID: id, Slot: slot, Thing: thing}, err
}

func parseRemoveContainerItem(m *message, ignore bool, offset time.Duration) (data.RemoveContainerItem, error) {
	id, err := m.getByte()
	if err != nil {
		return data.RemoveContainerItem{}, err
	}
	slot, err := m.getByte()
	return data.RemoveContainerItem{TimeOffset: offset, ContainerID: id, Slot: slot}, err
}

func parseInventoryItem(m *message, packetHead byte, ignore bool, offset time.Duration) (data.Operation, error) {
	if packetHead == 0x78 {
		slot, err := m.getByte()
		if err != nil {
			return nil, err
		}
		item, err := getItem(m, 0xFFFF)
		if err != nil {
			return nil, err
		}
		return data.InventorySetItem{TimeOffset: offset, Slot: slot, Item: item}, nil
	}
	// packetHead == 0x79
	slot, err := m.getByte()
	if err != nil {
		return nil, err
	}
	return data.InventoryClearItem{TimeOffset: offset, Slot: slot}, nil
}

func parseTradeItemRequest(m *message, ignore bool, offset time.Duration) (data.TradeItemRequest, error) {
	name, err := m.getString()
	if err != nil {
		return data.TradeItemRequest{}, err
	}
	size, err := m.getByte()
	if err != nil {
		return data.TradeItemRequest{}, err
	}
	op := data.TradeItemRequest{TimeOffset: offset, Name: name}
	for i := 0; i < int(size); i++ {
		thing, err := getThing(m)
		if err != nil {
			return op, err
		}
		op.Items = append(op.Items, thing)
	}
	return op, nil
}

func parseCloseTrade(m *message, ignore bool, offset time.Duration) (data.CloseTrade, error) {
	return data.CloseTrade{TimeOffset: offset}, nil
}

func parseWorldLight(m *message, ignore bool, offset time.Duration) (data.WorldLight, error) {
	level, err := m.getByte()
	if err != nil {
		return data.WorldLight{}, err
	}
	color, err := m.getByte()
	return data.WorldLight{TimeOffset: offset, Level: level, Color: color}, err
}

func parseMagicEffect(m *message, ignore bool, offset time.Duration) (data.MagicEffect, error) {
	loc, err := m.getLocation()
	if err != nil {
		return data.MagicEffect{}, err
	}
	effect, err := m.getByte()
	return data.MagicEffect{TimeOffset: offset, Location: loc, Effect: effect}, err
}

func parseAnimatedText(m *message, ignore bool, offset time.Duration) (data.AnimatedText, error) {
	loc, err := m.getLocation()
	if err != nil {
		return data.AnimatedText{}, err
	}
	color, err := m.getByte()
	if err != nil {
		return data.AnimatedText{}, err
	}
	text, err := m.getString()
	return data.AnimatedText{TimeOffset: offset, Location: loc, Color: color, Text: text}, err
}

func parseDistanceShoot(m *message, ignore bool, offset time.Duration) (data.DistanceShoot, error) {
	from, err := m.getLocation()
	if err != nil {
		return data.DistanceShoot{}, err
	}
	to, err := m.getLocation()
	if err != nil {
		return data.DistanceShoot{}, err
	}
	effect, err := m.getByte()
	return data.DistanceShoot{TimeOffset: offset, From: from, To: to, Effect: effect}, err
}

func parseCreatureSquare(m *message, ignore bool, offset time.Duration) (data.CreatureSquare, error) {
	id, err := m.getU32()
	if err != nil {
		return data.CreatureSquare{}, err
	}
	color, err := m.getByte()
	return data.CreatureSquare{TimeOffset: offset, CreatureID: id, Color: color}, err
}

func parseCreatureHealth(m *message, ignore bool, offset time.Duration) (data.CreatureHealth, error) {
	id, err := m.getU32()
	if err != nil {
		return data.CreatureHealth{}, err
	}
	health, err := m.getByte()
	return data.CreatureHealth{TimeOffset: offset, CreatureID: id, Health: health}, err
}

func parseCreatureLight(m *message, ignore bool, offset time.Duration) (data.CreatureLight, error) {
	id, err := m.getU32()
	if err != nil {
		return data.CreatureLight{}, err
	}
	level, err := m.getByte()
	if err != nil {
		return data.CreatureLight{}, err
	}
	color, err := m.getByte()
	return data.CreatureLight{TimeOffset: offset, CreatureID: id, Level: level, Color: color}, err
}

func parseCreatureOutfit(m *message, ignore bool, offset time.Duration) (data.CreatureOutfit, error) {
	id, err := m.getU32()
	if err != nil {
		return data.CreatureOutfit{}, err
	}
	outfit, err := m.getOutfit()
	return data.CreatureOutfit{TimeOffset: offset, CreatureID: id, Outfit: outfit}, err
}

func parseChangeSpeed(m *message, ignore bool, offset time.Duration) (data.ChangeSpeed, error) {
	id, err := m.getU32()
	if err != nil {
		return data.ChangeSpeed{}, err
	}
	speed, err := m.getU16()
	return data.ChangeSpeed{TimeOffset: offset, CreatureID: id, Speed: speed}, err
}

func parseCreatureSkull(m *message, ignore bool, offset time.Duration) (data.CreatureSkull, error) {
	id, err := m.getU32()
	if err != nil {
		return data.CreatureSkull{}, err
	}
	skull, err := m.getByte()
	return data.CreatureSkull{TimeOffset: offset, CreatureID: id, Skull: skull}, err
}

func parseCreatureShield(m *message, ignore bool, offset time.Duration) (data.CreatureShield, error) {
	id, err := m.getU32()
	if err != nil {
		return data.CreatureShield{}, err
	}
	shield, err := m.getByte()
	return data.CreatureShield{TimeOffset: offset, CreatureID: id, Shield: shield}, err
}

func parseTextWindow(m *message, ignore bool, offset time.Duration) (data.TextWindow, error) {
	op := data.TextWindow{TimeOffset: offset}
	var err error
	op.WindowID, err = m.getU32()
	if err != nil {
		return op, err
	}
	op.ItemID, err = m.getU16()
	if err != nil {
		return op, err
	}
	op.MaxLen, err = m.getU16()
	if err != nil {
		return op, err
	}
	op.Text, err = m.getString()
	if err != nil {
		return op, err
	}
	op.Author, err = m.getString()
	return op, err
}

func parseHouseWindow(m *message, ignore bool, offset time.Duration) (data.HouseWindow, error) {
	op := data.HouseWindow{TimeOffset: offset}
	var err error
	op.Unknown, err = m.getByte()
	if err != nil {
		return op, err
	}
	op.ID, err = m.getU32()
	if err != nil {
		return op, err
	}
	op.Text, err = m.getString()
	return op, err
}

func parsePlayerStats(m *message, ignore bool, offset time.Duration) (data.PlayerStats, error) {
	op := data.PlayerStats{TimeOffset: offset}
	var err error
	op.HP, err = m.getU16()
	if err != nil {
		return op, err
	}
	op.MaxHP, err = m.getU16()
	if err != nil {
		return op, err
	}
	op.Capacity, err = m.getU16()
	if err != nil {
		return op, err
	}
	op.Exp, err = m.getU32()
	if err != nil {
		return op, err
	}
	op.Level, err = m.getByte()
	if err != nil {
		return op, err
	}
	op.LevelPct, err = m.getByte()
	if err != nil {
		return op, err
	}
	op.Mana, err = m.getU16()
	if err != nil {
		return op, err
	}
	op.MaxMana, err = m.getU16()
	if err != nil {
		return op, err
	}
	op.MagicLvl, err = m.getByte()
	if err != nil {
		return op, err
	}
	op.MagicLvlPct, err = m.getByte()
	if err != nil {
		return op, err
	}
	op.Soul, err = m.getU16()
	return op, err
}

func parsePlayerSkills(m *message, ignore bool, offset time.Duration) (data.PlayerSkills, error) {
	op := data.PlayerSkills{TimeOffset: offset}
	for i := 0; i < 7; i++ {
		level, err := m.getByte()
		if err != nil {
			return op, err
		}
		pct, err := m.getByte()
		if err != nil {
			return op, err
		}
		op.Skills[i] = data.SkillValue{Level: level, Percent: pct}
	}
	return op, nil
}

func parsePlayerIcons(m *message, ignore bool, offset time.Duration) (data.PlayerIcons, error) {
	icons, err := m.getByte()
	return data.PlayerIcons{TimeOffset: offset, Icons: icons}, err
}

func parseCancelTarget(m *message, ignore bool, offset time.Duration) (data.CancelTarget, error) {
	return data.CancelTarget{TimeOffset: offset}, nil
}

func parseCreatureSpeak(m *message, ignore bool, offset time.Duration) (data.CreatureSpeak, error) {
	op := data.CreatureSpeak{TimeOffset: offset}
	var err error
	op.StatementID, err = m.getU32()
	if err != nil {
		return op, err
	}
	op.Name, err = m.getString()
	if err != nil {
		return op, err
	}
	op.Type, err = m.getByte()
	if err != nil {
		return op, err
	}
	switch op.Type {
	case 1, 2, 3, 0x10, 0x11:
		loc, err := m.getLocation()
		if err != nil {
			return op, err
		}
		op.Location = &loc
	case 5, 6, 0xA, 0xC, 0xE:
		ch, err := m.getU16()
		if err != nil {
			return op, err
		}
		op.ChannelID = &ch
	}
	op.Text, err = m.getString()
	return op, err
}

func parseChannelsDialog(m *message, ignore bool, offset time.Duration) (data.ChannelsDialog, error) {
	size, err := m.getByte()
	if err != nil {
		return data.ChannelsDialog{}, err
	}
	op := data.ChannelsDialog{TimeOffset: offset}
	for i := 0; i < int(size); i++ {
		id, err := m.getU16()
		if err != nil {
			return op, err
		}
		name, err := m.getString()
		if err != nil {
			return op, err
		}
		op.Channels = append(op.Channels, data.ChannelEntry{ID: id, Name: name})
	}
	return op, nil
}

func parseChannel(m *message, ignore bool, offset time.Duration) (data.Channel, error) {
	id, err := m.getU16()
	if err != nil {
		return data.Channel{}, err
	}
	name, err := m.getString()
	return data.Channel{TimeOffset: offset, ID: id, Name: name}, err
}

func parseOpenPrivateChannel(m *message, ignore bool, offset time.Duration) (data.OpenPrivateChannel, error) {
	name, err := m.getString()
	return data.OpenPrivateChannel{TimeOffset: offset, Name: name}, err
}

func parseRuleViolationsChannel(m *message, ignore bool, offset time.Duration) (data.RuleViolationsChannel, error) {
	size, err := m.getU16()
	return data.RuleViolationsChannel{TimeOffset: offset, Size: size}, err
}

func parseRemoveReport(m *message, ignore bool, offset time.Duration) (data.RemoveReport, error) {
	name, err := m.getString()
	return data.RemoveReport{TimeOffset: offset, Name: name}, err
}

func parseRuleViolationCancel(m *message, ignore bool, offset time.Duration) (data.RuleViolationCancel, error) {
	name, err := m.getString()
	return data.RuleViolationCancel{TimeOffset: offset, Name: name}, err
}

func parseLockRuleViolation(m *message, ignore bool, offset time.Duration) (data.LockRuleViolation, error) {
	return data.LockRuleViolation{TimeOffset: offset}, nil
}

func parseCreatePrivateChannel(m *message, ignore bool, offset time.Duration) (data.CreatePrivateChannel, error) {
	id, err := m.getU16()
	if err != nil {
		return data.CreatePrivateChannel{}, err
	}
	name, err := m.getString()
	return data.CreatePrivateChannel{TimeOffset: offset, ID: id, Name: name}, err
}

func parseClosePrivate(m *message, ignore bool, offset time.Duration) (data.ClosePrivate, error) {
	id, err := m.getU16()
	return data.ClosePrivate{TimeOffset: offset, ChannelID: id}, err
}

func parseTextMessage(m *message, ignore bool, offset time.Duration) (data.TextMessage, error) {
	t, err := m.getByte()
	if err != nil {
		return data.TextMessage{}, err
	}
	text, err := m.getString()
	return data.TextMessage{TimeOffset: offset, Type: t, Text: text}, err
}

func parseCancelWalk(m *message, ignore bool, offset time.Duration) (data.CancelWalk, error) {
	dir, err := m.getByte()
	return data.CancelWalk{TimeOffset: offset, Direction: data.Direction(dir)}, err
}

func parseFloorChangeUp(m *message, s *parseState, ignore bool, offset time.Duration) (data.FloorChangeUp, error) {
	myPos := s.playerPos
	myPos.Z--

	var tiles []data.Tile
	if myPos.Z == 7 {
		skip := 0
		for _, floor := range []struct{ z, offset int }{{5, 3}, {4, 4}, {3, 5}, {2, 6}, {1, 7}, {0, 8}} {
			t, err := parseFloorDescription(m, myPos.X-8, myPos.Y-6, floor.z, 18, 14, floor.offset, &skip)
			if err != nil {
				return data.FloorChangeUp{}, err
			}
			tiles = append(tiles, t...)
		}
	} else if myPos.Z > 7 {
		skip := 0
		t, err := parseFloorDescription(m, myPos.X-8, myPos.Y-6, myPos.Z-2, 18, 14, 3, &skip)
		if err != nil {
			return data.FloorChangeUp{}, err
		}
		tiles = append(tiles, t...)
	}

	s.playerPos = data.Location{X: myPos.X + 1, Y: myPos.Y + 1, Z: myPos.Z}
	s.updateTiles(tiles)
	return data.FloorChangeUp{TimeOffset: offset, Tiles: tiles}, nil
}

func parseFloorChangeDown(m *message, s *parseState, ignore bool, offset time.Duration) (data.FloorChangeDown, error) {
	myPos := s.playerPos
	myPos.Z++

	var tiles []data.Tile
	skip := 0
	if myPos.Z == 8 {
		for i, j := myPos.Z, -1; i < myPos.Z+3; i, j = i+1, j-1 {
			t, err := parseFloorDescription(m, myPos.X-8, myPos.Y-6, i, 18, 14, j, &skip)
			if err != nil {
				return data.FloorChangeDown{}, err
			}
			tiles = append(tiles, t...)
		}
	} else if myPos.Z > 8 && myPos.Z < 14 {
		t, err := parseFloorDescription(m, myPos.X-8, myPos.Y-6, myPos.Z+2, 18, 14, -3, &skip)
		if err != nil {
			return data.FloorChangeDown{}, err
		}
		tiles = append(tiles, t...)
	}

	s.playerPos = data.Location{X: myPos.X - 1, Y: myPos.Y - 1, Z: myPos.Z}
	s.updateTiles(tiles)
	return data.FloorChangeDown{TimeOffset: offset, Tiles: tiles}, nil
}

func parseOutfitWindow(m *message, ignore bool, offset time.Duration) (data.OutfitWindow, error) {
	outfit, err := m.getOutfit()
	if err != nil {
		return data.OutfitWindow{}, err
	}
	start, err := m.getU16()
	if err != nil {
		return data.OutfitWindow{}, err
	}
	end, err := m.getU16()
	return data.OutfitWindow{TimeOffset: offset, Outfit: outfit, OutfitStart: start, OutfitEnd: end}, err
}

func parseVIP(m *message, ignore bool, offset time.Duration) (data.VIP, error) {
	id, err := m.getU32()
	if err != nil {
		return data.VIP{}, err
	}
	name, err := m.getString()
	if err != nil {
		return data.VIP{}, err
	}
	online, err := m.getByte()
	return data.VIP{TimeOffset: offset, ID: id, Name: name, Online: online}, err
}

func parseVIPLogin(m *message, ignore bool, offset time.Duration) (data.VIPLogin, error) {
	id, err := m.getU32()
	return data.VIPLogin{TimeOffset: offset, ID: id}, err
}

func parseVIPLogout(m *message, ignore bool, offset time.Duration) (data.VIPLogout, error) {
	id, err := m.getU32()
	return data.VIPLogout{TimeOffset: offset, ID: id}, err
}
