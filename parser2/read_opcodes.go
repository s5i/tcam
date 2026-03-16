package parser2

import (
	"bytes"
)

// ReadInitGame reads a SV_CMD_INIT_GAME packet payload (after opcode).
func ReadInitGame(r *bytes.Reader) (TInitGame, error) {
	creatureID, err := ReadQuad(r)
	if err != nil {
		return TInitGame{}, err
	}
	beat, err := ReadWord(r)
	if err != nil {
		return TInitGame{}, err
	}
	canReport, err := ReadByte(r)
	if err != nil {
		return TInitGame{}, err
	}
	return TInitGame{
		CreatureID:    creatureID,
		Beat:          beat,
		CanReportBugs: canReport != 0,
	}, nil
}

// ReadRights reads a SV_CMD_RIGHTS packet payload (after opcode).
// Returns 32 bytes of action bitmasks.
func ReadRights(r *bytes.Reader) ([32]TByte, error) {
	var rights [32]TByte
	for i := 0; i < 32; i++ {
		b, err := ReadByte(r)
		if err != nil {
			return rights, err
		}
		rights[i] = b
	}
	return rights, nil
}

// ReadFullScreen reads a SV_CMD_FULLSCREEN packet payload (after opcode).
func ReadFullScreen(r *bytes.Reader) (TPosition, []TMapObject, error) {
	pos, err := ReadPosition(r)
	if err != nil {
		return TPosition{}, nil, err
	}
	objects, err := ReadMapDescription(r)
	if err != nil {
		return TPosition{}, nil, err
	}
	return pos, objects, nil
}

// ReadRow reads a SV_CMD_ROW_* packet payload (after opcode).
func ReadRow(r *bytes.Reader) ([]TMapObject, error) {
	return ReadMapDescription(r)
}

// ReadFloorChange reads a SV_CMD_FLOOR_UP or SV_CMD_FLOOR_DOWN packet payload (after opcode).
func ReadFloorChange(r *bytes.Reader) ([]TMapObject, error) {
	return ReadMapDescription(r)
}

// ReadFieldData reads a SV_CMD_FIELD_DATA packet payload (after opcode).
func ReadFieldData(r *bytes.Reader) (TPosition, []TMapObject, error) {
	pos, err := ReadPosition(r)
	if err != nil {
		return TPosition{}, nil, err
	}
	objects, err := ReadMapDescription(r)
	if err != nil {
		return TPosition{}, nil, err
	}
	return pos, objects, nil
}

// ReadAddField reads a SV_CMD_ADD_FIELD packet payload (after opcode).
func ReadAddField(r *bytes.Reader) (TFieldChange, error) {
	pos, err := ReadPosition(r)
	if err != nil {
		return TFieldChange{}, err
	}
	obj, err := ReadMapObject(r)
	if err != nil {
		return TFieldChange{}, err
	}
	return TFieldChange{Position: pos, Object: &obj}, nil
}

// ReadChangeField reads a SV_CMD_CHANGE_FIELD packet payload (after opcode).
func ReadChangeField(r *bytes.Reader) (TFieldChange, error) {
	pos, err := ReadPosition(r)
	if err != nil {
		return TFieldChange{}, err
	}
	stackIdx, err := ReadByte(r)
	if err != nil {
		return TFieldChange{}, err
	}
	obj, err := ReadMapObject(r)
	if err != nil {
		return TFieldChange{}, err
	}
	return TFieldChange{Position: pos, StackIndex: stackIdx, Object: &obj}, nil
}

// ReadDeleteField reads a SV_CMD_DELETE_FIELD packet payload (after opcode).
func ReadDeleteField(r *bytes.Reader) (TFieldChange, error) {
	pos, err := ReadPosition(r)
	if err != nil {
		return TFieldChange{}, err
	}
	stackIdx, err := ReadByte(r)
	if err != nil {
		return TFieldChange{}, err
	}
	return TFieldChange{Position: pos, StackIndex: stackIdx}, nil
}

// ReadMoveCreature reads a SV_CMD_MOVE_CREATURE packet payload (after opcode).
func ReadMoveCreature(r *bytes.Reader) (TMoveCreature, error) {
	origPos, err := ReadPosition(r)
	if err != nil {
		return TMoveCreature{}, err
	}
	origIdx, err := ReadByte(r)
	if err != nil {
		return TMoveCreature{}, err
	}
	destPos, err := ReadPosition(r)
	if err != nil {
		return TMoveCreature{}, err
	}
	return TMoveCreature{
		OrigPosition: origPos,
		OrigIndex:    origIdx,
		DestPosition: destPos,
	}, nil
}

// ReadContainer reads a SV_CMD_CONTAINER packet payload (after opcode).
func ReadContainer(r *bytes.Reader) (TContainer, error) {
	nr, err := ReadByte(r)
	if err != nil {
		return TContainer{}, err
	}
	typeID, err := ReadWord(r)
	if err != nil {
		return TContainer{}, err
	}
	name, err := ReadString(r)
	if err != nil {
		return TContainer{}, err
	}
	capacity, err := ReadByte(r)
	if err != nil {
		return TContainer{}, err
	}
	hasParent, err := ReadByte(r)
	if err != nil {
		return TContainer{}, err
	}
	itemCount, err := ReadByte(r)
	if err != nil {
		return TContainer{}, err
	}

	items := make([]TItem, 0, int(itemCount))
	for i := 0; i < int(itemCount); i++ {
		item, err := ReadItem(r)
		if err != nil {
			return TContainer{}, err
		}
		items = append(items, item)
	}

	return TContainer{
		ContainerNr: nr,
		TypeID:      typeID,
		Name:        name,
		Capacity:    capacity,
		HasParent:   hasParent != 0,
		Items:       items,
	}, nil
}

// ReadCloseContainer reads a SV_CMD_CLOSE_CONTAINER packet payload (after opcode).
func ReadCloseContainer(r *bytes.Reader) (TByte, error) {
	return ReadByte(r)
}

// ReadCreateInContainer reads a SV_CMD_CREATE_IN_CONTAINER packet payload (after opcode).
func ReadCreateInContainer(r *bytes.Reader) (TByte, TItem, error) {
	nr, err := ReadByte(r)
	if err != nil {
		return 0, TItem{}, err
	}
	item, err := ReadItem(r)
	if err != nil {
		return 0, TItem{}, err
	}
	return nr, item, nil
}

// ReadChangeInContainer reads a SV_CMD_CHANGE_IN_CONTAINER packet payload (after opcode).
func ReadChangeInContainer(r *bytes.Reader) (TByte, TByte, TItem, error) {
	nr, err := ReadByte(r)
	if err != nil {
		return 0, 0, TItem{}, err
	}
	idx, err := ReadByte(r)
	if err != nil {
		return 0, 0, TItem{}, err
	}
	item, err := ReadItem(r)
	if err != nil {
		return 0, 0, TItem{}, err
	}
	return nr, idx, item, nil
}

// ReadDeleteInContainer reads a SV_CMD_DELETE_IN_CONTAINER packet payload (after opcode).
func ReadDeleteInContainer(r *bytes.Reader) (TByte, TByte, error) {
	nr, err := ReadByte(r)
	if err != nil {
		return 0, 0, err
	}
	idx, err := ReadByte(r)
	if err != nil {
		return 0, 0, err
	}
	return nr, idx, nil
}

// ReadSetInventory reads a SV_CMD_SET_INVENTORY packet payload (after opcode).
func ReadSetInventory(r *bytes.Reader) (TByte, TItem, error) {
	pos, err := ReadByte(r)
	if err != nil {
		return 0, TItem{}, err
	}
	item, err := ReadItem(r)
	if err != nil {
		return 0, TItem{}, err
	}
	return pos, item, nil
}

// ReadDeleteInventory reads a SV_CMD_DELETE_INVENTORY packet payload (after opcode).
func ReadDeleteInventory(r *bytes.Reader) (TByte, error) {
	return ReadByte(r)
}

// ReadTradeOffer reads a SV_CMD_TRADE_OFFER_OWN or SV_CMD_TRADE_OFFER_PARTNER
// packet payload (after opcode).
func ReadTradeOffer(r *bytes.Reader) (TTradeOffer, error) {
	name, err := ReadString(r)
	if err != nil {
		return TTradeOffer{}, err
	}
	count, err := ReadByte(r)
	if err != nil {
		return TTradeOffer{}, err
	}

	items := make([]TItem, 0, int(count))
	for i := 0; i < int(count); i++ {
		item, err := ReadItem(r)
		if err != nil {
			return TTradeOffer{}, err
		}
		items = append(items, item)
	}

	return TTradeOffer{Name: name, Items: items}, nil
}

// ReadAmbiente reads a SV_CMD_AMBIENTE packet payload (after opcode).
func ReadAmbiente(r *bytes.Reader) (TAmbiente, error) {
	brightness, err := ReadByte(r)
	if err != nil {
		return TAmbiente{}, err
	}
	color, err := ReadByte(r)
	if err != nil {
		return TAmbiente{}, err
	}
	return TAmbiente{Brightness: brightness, Color: color}, nil
}

// ReadGraphicalEffect reads a SV_CMD_GRAPHICAL_EFFECT packet payload (after opcode).
func ReadGraphicalEffect(r *bytes.Reader) (TGraphicalEffect, error) {
	pos, err := ReadPosition(r)
	if err != nil {
		return TGraphicalEffect{}, err
	}
	typ, err := ReadByte(r)
	if err != nil {
		return TGraphicalEffect{}, err
	}
	return TGraphicalEffect{Position: pos, Type: typ}, nil
}

// ReadTextualEffect reads a SV_CMD_TEXTUAL_EFFECT packet payload (after opcode).
func ReadTextualEffect(r *bytes.Reader) (TTextualEffect, error) {
	pos, err := ReadPosition(r)
	if err != nil {
		return TTextualEffect{}, err
	}
	color, err := ReadByte(r)
	if err != nil {
		return TTextualEffect{}, err
	}
	text, err := ReadString(r)
	if err != nil {
		return TTextualEffect{}, err
	}
	return TTextualEffect{Position: pos, Color: color, Text: text}, nil
}

// ReadMissileEffect reads a SV_CMD_MISSILE_EFFECT packet payload (after opcode).
func ReadMissileEffect(r *bytes.Reader) (TMissileEffect, error) {
	orig, err := ReadPosition(r)
	if err != nil {
		return TMissileEffect{}, err
	}
	dest, err := ReadPosition(r)
	if err != nil {
		return TMissileEffect{}, err
	}
	typ, err := ReadByte(r)
	if err != nil {
		return TMissileEffect{}, err
	}
	return TMissileEffect{Origin: orig, Destination: dest, Type: typ}, nil
}

// ReadMarkCreature reads a SV_CMD_MARK_CREATURE packet payload (after opcode).
func ReadMarkCreature(r *bytes.Reader) (TMarkCreature, error) {
	id, err := ReadQuad(r)
	if err != nil {
		return TMarkCreature{}, err
	}
	color, err := ReadByte(r)
	if err != nil {
		return TMarkCreature{}, err
	}
	return TMarkCreature{CreatureID: id, Color: color}, nil
}

// ReadCreatureHealth reads a SV_CMD_CREATURE_HEALTH packet payload (after opcode).
func ReadCreatureHealth(r *bytes.Reader) (TCreatureHealthUpdate, error) {
	id, err := ReadQuad(r)
	if err != nil {
		return TCreatureHealthUpdate{}, err
	}
	health, err := ReadByte(r)
	if err != nil {
		return TCreatureHealthUpdate{}, err
	}
	return TCreatureHealthUpdate{CreatureID: id, Health: health}, nil
}

// ReadCreatureLight reads a SV_CMD_CREATURE_LIGHT packet payload (after opcode).
func ReadCreatureLight(r *bytes.Reader) (TCreatureLightUpdate, error) {
	id, err := ReadQuad(r)
	if err != nil {
		return TCreatureLightUpdate{}, err
	}
	brightness, err := ReadByte(r)
	if err != nil {
		return TCreatureLightUpdate{}, err
	}
	color, err := ReadByte(r)
	if err != nil {
		return TCreatureLightUpdate{}, err
	}
	return TCreatureLightUpdate{CreatureID: id, Brightness: brightness, Color: color}, nil
}

// ReadCreatureOutfit reads a SV_CMD_CREATURE_OUTFIT packet payload (after opcode).
func ReadCreatureOutfit(r *bytes.Reader) (TCreatureOutfitUpdate, error) {
	id, err := ReadQuad(r)
	if err != nil {
		return TCreatureOutfitUpdate{}, err
	}
	outfit, err := ReadOutfit(r)
	if err != nil {
		return TCreatureOutfitUpdate{}, err
	}
	return TCreatureOutfitUpdate{CreatureID: id, Outfit: outfit}, nil
}

// ReadCreatureSpeed reads a SV_CMD_CREATURE_SPEED packet payload (after opcode).
func ReadCreatureSpeed(r *bytes.Reader) (TCreatureSpeedUpdate, error) {
	id, err := ReadQuad(r)
	if err != nil {
		return TCreatureSpeedUpdate{}, err
	}
	speed, err := ReadWord(r)
	if err != nil {
		return TCreatureSpeedUpdate{}, err
	}
	return TCreatureSpeedUpdate{CreatureID: id, Speed: speed}, nil
}

// ReadCreatureSkull reads a SV_CMD_CREATURE_SKULL packet payload (after opcode).
func ReadCreatureSkull(r *bytes.Reader) (TCreatureUpdate, TByte, error) {
	id, err := ReadQuad(r)
	if err != nil {
		return TCreatureUpdate{}, 0, err
	}
	skull, err := ReadByte(r)
	if err != nil {
		return TCreatureUpdate{}, 0, err
	}
	return TCreatureUpdate{CreatureID: id}, skull, nil
}

// ReadCreatureParty reads a SV_CMD_CREATURE_PARTY packet payload (after opcode).
func ReadCreatureParty(r *bytes.Reader) (TCreatureUpdate, TByte, error) {
	id, err := ReadQuad(r)
	if err != nil {
		return TCreatureUpdate{}, 0, err
	}
	party, err := ReadByte(r)
	if err != nil {
		return TCreatureUpdate{}, 0, err
	}
	return TCreatureUpdate{CreatureID: id}, party, nil
}

// ReadEditText reads a SV_CMD_EDIT_TEXT packet payload (after opcode).
func ReadEditText(r *bytes.Reader) (TEditText, error) {
	objID, err := ReadQuad(r)
	if err != nil {
		return TEditText{}, err
	}
	typeID, err := ReadWord(r)
	if err != nil {
		return TEditText{}, err
	}
	maxLen, err := ReadWord(r)
	if err != nil {
		return TEditText{}, err
	}
	text, err := ReadString(r)
	if err != nil {
		return TEditText{}, err
	}
	editor, err := ReadString(r)
	if err != nil {
		return TEditText{}, err
	}
	return TEditText{
		ObjectID:  objID,
		TypeID:    typeID,
		MaxLength: maxLen,
		Text:      text,
		Editor:    editor,
	}, nil
}

// ReadEditList reads a SV_CMD_EDIT_LIST packet payload (after opcode).
func ReadEditList(r *bytes.Reader) (TEditList, error) {
	typ, err := ReadByte(r)
	if err != nil {
		return TEditList{}, err
	}
	id, err := ReadQuad(r)
	if err != nil {
		return TEditList{}, err
	}
	text, err := ReadString(r)
	if err != nil {
		return TEditList{}, err
	}
	return TEditList{Type: typ, ID: id, Text: text}, nil
}

// ReadPlayerData reads a SV_CMD_PLAYER_DATA packet payload (after opcode).
func ReadPlayerData(r *bytes.Reader) (TPlayerData, error) {
	var pd TPlayerData
	var err error

	if pd.HitPoints, err = ReadWord(r); err != nil {
		return TPlayerData{}, err
	}
	if pd.MaxHitPoints, err = ReadWord(r); err != nil {
		return TPlayerData{}, err
	}
	if pd.Capacity, err = ReadWord(r); err != nil {
		return TPlayerData{}, err
	}
	if pd.Experience, err = ReadQuad(r); err != nil {
		return TPlayerData{}, err
	}
	if pd.Level, err = ReadWord(r); err != nil {
		return TPlayerData{}, err
	}
	if pd.LevelPercent, err = ReadByte(r); err != nil {
		return TPlayerData{}, err
	}
	if pd.ManaPoints, err = ReadWord(r); err != nil {
		return TPlayerData{}, err
	}
	if pd.MaxManaPoints, err = ReadWord(r); err != nil {
		return TPlayerData{}, err
	}
	if pd.MagicLevel, err = ReadByte(r); err != nil {
		return TPlayerData{}, err
	}
	if pd.MagicLevelPercent, err = ReadByte(r); err != nil {
		return TPlayerData{}, err
	}
	if pd.SoulPoints, err = ReadByte(r); err != nil {
		return TPlayerData{}, err
	}
	return pd, nil
}

// ReadPlayerSkills reads a SV_CMD_PLAYER_SKILLS packet payload (after opcode).
func ReadPlayerSkills(r *bytes.Reader) (TPlayerSkills, error) {
	readSkill := func() (TSkillEntry, error) {
		level, err := ReadByte(r)
		if err != nil {
			return TSkillEntry{}, err
		}
		pct, err := ReadByte(r)
		if err != nil {
			return TSkillEntry{}, err
		}
		return TSkillEntry{Level: level, Percent: pct}, nil
	}

	var ps TPlayerSkills
	var err error
	if ps.Fist, err = readSkill(); err != nil {
		return TPlayerSkills{}, err
	}
	if ps.Club, err = readSkill(); err != nil {
		return TPlayerSkills{}, err
	}
	if ps.Sword, err = readSkill(); err != nil {
		return TPlayerSkills{}, err
	}
	if ps.Axe, err = readSkill(); err != nil {
		return TPlayerSkills{}, err
	}
	if ps.Distance, err = readSkill(); err != nil {
		return TPlayerSkills{}, err
	}
	if ps.Shielding, err = readSkill(); err != nil {
		return TPlayerSkills{}, err
	}
	if ps.Fishing, err = readSkill(); err != nil {
		return TPlayerSkills{}, err
	}
	return ps, nil
}

// ReadPlayerState reads a SV_CMD_PLAYER_STATE packet payload (after opcode).
func ReadPlayerState(r *bytes.Reader) (TByte, error) {
	return ReadByte(r)
}

// ReadTalk reads a SV_CMD_TALK packet payload (after opcode).
// The mode determines which extra fields follow.
//
// Talk modes (from sending.cc):
//
//	Say/Whisper/Yell/AnimalLow/AnimalLoud: position (5 bytes) + text
//	ChannelCall/GM channel/Highlight/Anonymous: channel (2 bytes) + text
//	PrivateMessage/GM request/GM answer/PlayerAnswer/GM broadcast/GM message: text only
//	  (GM request also has 4-byte ExtraData before text)
func ReadTalk(r *bytes.Reader) (TTalk, error) {
	var t TTalk
	var err error

	if t.StatementID, err = ReadQuad(r); err != nil {
		return TTalk{}, err
	}
	if t.Sender, err = ReadString(r); err != nil {
		return TTalk{}, err
	}
	if t.Mode, err = ReadByte(r); err != nil {
		return TTalk{}, err
	}

	switch t.Mode {
	case 1, 2, 3, 16, 17:
		// TALK_SAY=1, TALK_WHISPER=2, TALK_YELL=3, TALK_ANIMAL_LOW=16, TALK_ANIMAL_LOUD=17
		// Position: x(2) + y(2) + z(1)
		if t.Position, err = ReadPosition(r); err != nil {
			return TTalk{}, err
		}

	case 5, 10, 11, 12:
		// TALK_CHANNEL_CALL=5, TALK_GAMEMASTER_CHANNELCALL=10,
		// TALK_HIGHLIGHT_CHANNELCALL=11, TALK_ANONYMOUS_CHANNELCALL=12
		// Channel ID: 2 bytes
		if t.ChannelID, err = ReadWord(r); err != nil {
			return TTalk{}, err
		}

	case 4, 9, 13, 14:
		// TALK_PRIVATE_MESSAGE=4, TALK_GAMEMASTER_BROADCAST=9,
		// TALK_GAMEMASTER_ANSWER=13, TALK_PLAYER_ANSWER=14
		// No extra data before text.

	case 6:
		// TALK_GAMEMASTER_MESSAGE=6 — no extra data.

	case 7:
		// TALK_GAMEMASTER_REQUEST=7 — 4-byte extra data.
		if t.ExtraData, err = ReadQuad(r); err != nil {
			return TTalk{}, err
		}

	default:
		// Unknown mode — try to continue by reading text only.
	}

	if t.Text, err = ReadString(r); err != nil {
		return TTalk{}, err
	}

	return t, nil
}

// ReadChannels reads a SV_CMD_CHANNELS packet payload (after opcode).
func ReadChannels(r *bytes.Reader) ([]TChannel, error) {
	count, err := ReadByte(r)
	if err != nil {
		return nil, err
	}

	channels := make([]TChannel, 0, int(count))
	for i := 0; i < int(count); i++ {
		id, err := ReadWord(r)
		if err != nil {
			return nil, err
		}
		name, err := ReadString(r)
		if err != nil {
			return nil, err
		}
		channels = append(channels, TChannel{ID: id, Name: name})
	}

	return channels, nil
}

// ReadOpenChannel reads a SV_CMD_OPEN_CHANNEL or SV_CMD_OPEN_OWN_CHANNEL
// packet payload (after opcode).
func ReadOpenChannel(r *bytes.Reader) (TChannel, error) {
	id, err := ReadWord(r)
	if err != nil {
		return TChannel{}, err
	}
	name, err := ReadString(r)
	if err != nil {
		return TChannel{}, err
	}
	return TChannel{ID: id, Name: name}, nil
}

// ReadPrivateChannel reads a SV_CMD_PRIVATE_CHANNEL packet payload (after opcode).
func ReadPrivateChannel(r *bytes.Reader) (TString, error) {
	return ReadString(r)
}

// ReadOpenRequestQueue reads a SV_CMD_OPEN_REQUEST_QUEUE packet payload (after opcode).
func ReadOpenRequestQueue(r *bytes.Reader) (TWord, error) {
	return ReadWord(r)
}

// ReadDeleteRequest reads a SV_CMD_DELETE_REQUEST packet payload (after opcode).
func ReadDeleteRequest(r *bytes.Reader) (TString, error) {
	return ReadString(r)
}

// ReadFinishRequest reads a SV_CMD_FINISH_REQUEST packet payload (after opcode).
func ReadFinishRequest(r *bytes.Reader) (TString, error) {
	return ReadString(r)
}

// ReadCloseChannel reads a SV_CMD_CLOSE_CHANNEL packet payload (after opcode).
func ReadCloseChannel(r *bytes.Reader) (TWord, error) {
	return ReadWord(r)
}

// ReadMessage reads a SV_CMD_MESSAGE packet payload (after opcode).
func ReadMessage(r *bytes.Reader) (TMessage, error) {
	mode, err := ReadByte(r)
	if err != nil {
		return TMessage{}, err
	}
	text, err := ReadString(r)
	if err != nil {
		return TMessage{}, err
	}
	return TMessage{Mode: mode, Text: text}, nil
}

// ReadSnapback reads a SV_CMD_SNAPBACK packet payload (after opcode).
func ReadSnapback(r *bytes.Reader) (TByte, error) {
	return ReadByte(r)
}

// ReadOutfitWindow reads a SV_CMD_OUTFIT packet payload (after opcode).
func ReadOutfitWindow(r *bytes.Reader) (TOutfitWindow, error) {
	outfit, err := ReadOutfit(r)
	if err != nil {
		return TOutfitWindow{}, err
	}
	first, err := ReadWord(r)
	if err != nil {
		return TOutfitWindow{}, err
	}
	last, err := ReadWord(r)
	if err != nil {
		return TOutfitWindow{}, err
	}
	return TOutfitWindow{CurrentOutfit: outfit, FirstOutfit: first, LastOutfit: last}, nil
}

// ReadBuddyData reads a SV_CMD_BUDDY_DATA packet payload (after opcode).
func ReadBuddyData(r *bytes.Reader) (TBuddyData, error) {
	id, err := ReadQuad(r)
	if err != nil {
		return TBuddyData{}, err
	}
	name, err := ReadString(r)
	if err != nil {
		return TBuddyData{}, err
	}
	online, err := ReadByte(r)
	if err != nil {
		return TBuddyData{}, err
	}
	return TBuddyData{CharacterID: id, Name: name, Online: online != 0}, nil
}

// ReadBuddyOnline reads a SV_CMD_BUDDY_ONLINE packet payload (after opcode).
func ReadBuddyOnline(r *bytes.Reader) (TQuad, error) {
	return ReadQuad(r)
}

// ReadBuddyOffline reads a SV_CMD_BUDDY_OFFLINE packet payload (after opcode).
func ReadBuddyOffline(r *bytes.Reader) (TQuad, error) {
	return ReadQuad(r)
}
