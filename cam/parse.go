package cam

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"time"

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

			ops, err := parsePacket(state, packet.Data, packet.TimeOffset)
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

func parsePacket(state *parseState, buf []byte, timeOffset time.Duration) ([]data.Operation, error) {
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
			op, err = parseLogin(m, false, timeOffset)
		case 0x14:
			op, err = parseDisconnectClient(m, false, timeOffset)
		case 0x16:
			op, err = parseWaitList(m, false, timeOffset)
		case 0x1E:
			op, err = parsePing(m, false, timeOffset)
		case 0x64:
			op, err = parseMapDescription(m, state, false, timeOffset)
		case 0x65:
			op, err = parseMoveNorth(m, state, false, timeOffset)
		case 0x66:
			op, err = parseMoveEast(m, state, false, timeOffset)
		case 0x67:
			op, err = parseMoveSouth(m, state, false, timeOffset)
		case 0x68:
			op, err = parseMoveWest(m, state, false, timeOffset)
		case 0x69:
			op, err = parseUpdateTile(m, state, false, timeOffset)
		case 0x6A:
			op, err = parseAddTileItem(m, state, false, timeOffset)
		case 0x6B:
			op, err = parseUpdateTileItem(m, state, false, timeOffset)
		case 0x6C:
			op, err = parseRemoveTileItem(m, state, false, timeOffset)
		case 0x6D:
			op, err = parseMoveCreature(m, state, false, timeOffset)
		case 0x6E:
			op, err = parseContainer(m, false, timeOffset)
		case 0x6F:
			op, err = parseCloseContainer(m, false, timeOffset)
		case 0x70:
			op, err = parseAddContainerItem(m, false, timeOffset)
		case 0x71:
			op, err = parseUpdateContainerItem(m, false, timeOffset)
		case 0x72:
			op, err = parseRemoveContainerItem(m, false, timeOffset)
		case 0x78, 0x79:
			op, err = parseInventoryItem(m, head, false, timeOffset)
		case 0x7D, 0x7E:
			op, err = parseTradeItemRequest(m, false, timeOffset)
		case 0x7F:
			op, err = parseCloseTrade(m, false, timeOffset)
		case 0x82:
			op, err = parseWorldLight(m, false, timeOffset)
		case 0x83:
			op, err = parseMagicEffect(m, false, timeOffset)
		case 0x84:
			op, err = parseAnimatedText(m, false, timeOffset)
		case 0x85:
			op, err = parseDistanceShoot(m, false, timeOffset)
		case 0x86:
			op, err = parseCreatureSquare(m, false, timeOffset)
		case 0x8C:
			op, err = parseCreatureHealth(m, false, timeOffset)
		case 0x8D:
			op, err = parseCreatureLight(m, false, timeOffset)
		case 0x8E:
			op, err = parseCreatureOutfit(m, false, timeOffset)
		case 0x8F:
			op, err = parseChangeSpeed(m, false, timeOffset)
		case 0x90:
			op, err = parseCreatureSkull(m, false, timeOffset)
		case 0x91:
			op, err = parseCreatureShield(m, false, timeOffset)
		case 0x96:
			op, err = parseTextWindow(m, false, timeOffset)
		case 0x97:
			op, err = parseHouseWindow(m, false, timeOffset)
		case 0xA0:
			op, err = parsePlayerStats(m, false, timeOffset)
		case 0xA1:
			op, err = parsePlayerSkills(m, false, timeOffset)
		case 0xA2:
			op, err = parsePlayerIcons(m, false, timeOffset)
		case 0xA3:
			op, err = parseCancelTarget(m, false, timeOffset)
		case 0xAA:
			op, err = parseCreatureSpeak(m, false, timeOffset)
		case 0xAB:
			op, err = parseChannelsDialog(m, false, timeOffset)
		case 0xAC:
			op, err = parseChannel(m, false, timeOffset)
		case 0xAD:
			op, err = parseOpenPrivateChannel(m, false, timeOffset)
		case 0xAE:
			op, err = parseRuleViolationsChannel(m, false, timeOffset)
		case 0xAF:
			op, err = parseRemoveReport(m, false, timeOffset)
		case 0xB0:
			op, err = parseRuleViolationCancel(m, false, timeOffset)
		case 0xB1:
			op, err = parseLockRuleViolation(m, false, timeOffset)
		case 0xB2:
			op, err = parseCreatePrivateChannel(m, false, timeOffset)
		case 0xB3:
			op, err = parseClosePrivate(m, false, timeOffset)
		case 0xB4:
			op, err = parseTextMessage(m, false, timeOffset)
		case 0xB5:
			op, err = parseCancelWalk(m, false, timeOffset)
		case 0xBE:
			op, err = parseFloorChangeUp(m, state, false, timeOffset)
		case 0xBF:
			op, err = parseFloorChangeDown(m, state, false, timeOffset)
		case 0xC8:
			op, err = parseOutfitWindow(m, false, timeOffset)
		case 0xD2:
			op, err = parseVIP(m, false, timeOffset)
		case 0xD3:
			op, err = parseVIPLogin(m, false, timeOffset)
		case 0xD4:
			op, err = parseVIPLogout(m, false, timeOffset)
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
