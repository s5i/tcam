package cam

import (
	"errors"
	"fmt"
	"io"
	"iter"

	"github.com/s5i/tcam/data"
)

type tileKey struct {
	x, y, z int
}

type parseState struct {
	playerPos data.Location
	tiles     map[tileKey][]data.Thing
}

func (s *parseState) updateTiles(tiles []data.Tile) {
	for _, t := range tiles {
		k := tileKey{t.Location.X, t.Location.Y, t.Location.Z}
		s.tiles[k] = append([]data.Thing(nil), t.Things...)
	}
}

func (s *parseState) getThing(loc data.Location, stack int) *data.Thing {
	k := tileKey{loc.X, loc.Y, loc.Z}
	things := s.tiles[k]
	if stack < 0 || stack >= len(things) {
		return nil
	}
	return &things[stack]
}

func (s *parseState) addThing(loc data.Location, thing data.Thing) {
	k := tileKey{loc.X, loc.Y, loc.Z}
	s.tiles[k] = append(s.tiles[k], thing)
}

func (s *parseState) removeThing(loc data.Location, stack int) {
	k := tileKey{loc.X, loc.Y, loc.Z}
	things := s.tiles[k]
	if stack < 0 || stack >= len(things) {
		return
	}
	s.tiles[k] = append(things[:stack], things[stack+1:]...)
}

func (s *parseState) replaceThing(loc data.Location, stack int, thing data.Thing) {
	k := tileKey{loc.X, loc.Y, loc.Z}
	things := s.tiles[k]
	if stack < 0 || stack >= len(things) {
		return
	}
	things[stack] = thing
}

// Parse returns an iterator over the provided io.ReadSeeker that returns subsequent data.Operations.
func Parse(r io.ReadSeeker) iter.Seq2[data.Operation, error] {
	return func(yield func(data.Operation, error) bool) {
		yieldVal := func(p data.Operation) bool { return yield(p, nil) }
		yieldErr := func(err error) {
			if !errors.Is(err, io.EOF) {
				yield(nil, err)
			}
		}

		state := &parseState{
			tiles: make(map[tileKey][]data.Thing),
		}

		for packet, err := range Read(r) {
			if err != nil {
				yieldErr(err)
				return
			}

			ops, err := parsePacket(state, packet.Data)
			if err != nil {
				yieldErr(fmt.Errorf("at file offset %d: %w", packet.FileOffset, err))
				return
			}
			for _, op := range ops {
				if !yieldVal(op) {
					return
				}
			}
		}
	}
}

func parsePacket(state *parseState, buf []byte) ([]data.Operation, error) {
	m := newMessage(buf)
	var ops []data.Operation

	for m.remaining() > 0 {
		head, err := m.getByte()
		if err != nil {
			return ops, fmt.Errorf("reading packet head: %w", err)
		}

		var op data.Operation
		switch head {
		case 0x0A:
			op, err = parseLogin(m)
		case 0x14:
			op, err = parseDisconnectClient(m)
		case 0x16:
			op, err = parseWaitList(m)
		case 0x1E:
			op = data.Ping{}
		case 0x64:
			op, err = parseMapDescription(m, state)
		case 0x65:
			op, err = parseMoveNorth(m, state)
		case 0x66:
			op, err = parseMoveEast(m, state)
		case 0x67:
			op, err = parseMoveSouth(m, state)
		case 0x68:
			op, err = parseMoveWest(m, state)
		case 0x69:
			op, err = parseUpdateTile(m, state)
		case 0x6A:
			op, err = parseAddTileItem(m, state)
		case 0x6B:
			op, err = parseUpdateTileItem(m, state)
		case 0x6C:
			op, err = parseRemoveTileItem(m, state)
		case 0x6D:
			op, err = parseMoveCreature(m, state)
		case 0x6E:
			op, err = parseContainer(m)
		case 0x6F:
			op, err = parseCloseContainer(m)
		case 0x70:
			op, err = parseAddContainerItem(m)
		case 0x71:
			op, err = parseUpdateContainerItem(m)
		case 0x72:
			op, err = parseRemoveContainerItem(m)
		case 0x78, 0x79:
			op, err = parseInventoryItem(m, head)
		case 0x7D, 0x7E:
			op, err = parseTradeItemRequest(m)
		case 0x7F:
			op = data.CloseTrade{}
		case 0x82:
			op, err = parseWorldLight(m)
		case 0x83:
			op, err = parseMagicEffect(m)
		case 0x84:
			op, err = parseAnimatedText(m)
		case 0x85:
			op, err = parseDistanceShoot(m)
		case 0x86:
			op, err = parseCreatureSquare(m)
		case 0x8C:
			op, err = parseCreatureHealth(m)
		case 0x8D:
			op, err = parseCreatureLight(m)
		case 0x8E:
			op, err = parseCreatureOutfit(m)
		case 0x8F:
			op, err = parseChangeSpeed(m)
		case 0x90:
			op, err = parseCreatureSkull(m)
		case 0x91:
			op, err = parseCreatureShield(m)
		case 0x96:
			op, err = parseTextWindow(m)
		case 0x97:
			op, err = parseHouseWindow(m)
		case 0xA0:
			op, err = parsePlayerStats(m)
		case 0xA1:
			op, err = parsePlayerSkills(m)
		case 0xA2:
			op, err = parsePlayerIcons(m)
		case 0xA3:
			op = data.CancelTarget{}
		case 0xAA:
			op, err = parseCreatureSpeak(m)
		case 0xAB:
			op, err = parseChannelsDialog(m)
		case 0xAC:
			op, err = parseChannel(m)
		case 0xAD:
			op, err = parseOpenPrivateChannel(m)
		case 0xAE:
			op, err = parseRuleViolationsChannel(m)
		case 0xAF:
			op, err = parseRemoveReport(m)
		case 0xB0:
			op, err = parseRuleViolationCancel(m)
		case 0xB1:
			op = data.LockRuleViolation{}
		case 0xB2:
			op, err = parseCreatePrivateChannel(m)
		case 0xB3:
			op, err = parseClosePrivate(m)
		case 0xB4:
			op, err = parseTextMessage(m)
		case 0xB5:
			op, err = parseCancelWalk(m)
		case 0xBE:
			op, err = parseFloorChangeUp(m, state)
		case 0xBF:
			op, err = parseFloorChangeDown(m, state)
		case 0xC8:
			op, err = parseOutfitWindow(m)
		case 0xD2:
			op, err = parseVIP(m)
		case 0xD3:
			op, err = parseVIPLogin(m)
		case 0xD4:
			op, err = parseVIPLogout(m)
		default:
			return ops, fmt.Errorf("unknown packet head: 0x%02X", head)
		}

		if err != nil {
			return ops, fmt.Errorf("parsing 0x%02X: %w", head, err)
		}
		ops = append(ops, op)
	}
	return ops, nil
}
