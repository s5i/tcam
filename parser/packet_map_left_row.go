package parser

import (
	"bytes"
	"fmt"

	"github.com/s5i/tcam/enum"
	"github.com/s5i/tcam/network"
)

type MapLeftRow struct{}

// [104
// 214 003 008 005 000 255
// 214 003 001 005 000 255
// 214 003 199 006 000 255
// 214 003 001 005 092 011 000 255
// 214 003 006 005 000 255
// 214 003 000 255
// 214 003 000 255
// 214 003 000 255 216 003 000 255
// 231 000 000 255
// 231 000 000 255
// 231 000 000 255
// 231 000 000 255
// 231 000 000 255
// 173 001 134 009 000 255
// 173 001 135 009 000 255
// 173 001 136 009 000 255
// 173 001 002 005 010 255
// 128 004 000 255
// 128 004 000 255
// 128 004 081 255]
func parseMapLeftRow(p *network.Packet) (*MapLeftRow, *network.Packet, error) {
	if p.OpCode() != enum.OpCodeMapLeftRow {
		return nil, nil, fmt.Errorf("expected op code %s, got %s", enum.OpCodeMapLeftRow, p.OpCode())
	}

	r := bytes.NewReader(p.Data)

	if err := opcode(r); err != nil {
		return nil, nil, err
	}

	//  Position pos;
	// if (g_game.getFeature(Otc::GameMapMovePosition))
	//     pos = getPosition(msg);
	// else
	//     pos = g_map.getCentralPosition();
	// pos.x--;

	// g_map.setCentralPosition(pos);

	// g_map.cleanTile(position);
	// if (msg->peekU16() >= 0xff00)
	//     return msg->getU16() & 0xff;

	// if (g_game.getFeature(Otc::GameNewWalking)) {
	//     uint16_t groundSpeed = msg->getU16();
	//     uint8_t blocking = msg->getU8();
	//     g_map.setTileSpeed(position, groundSpeed, blocking);
	// }

	// if (g_game.getFeature(Otc::GameEnvironmentEffect) && !g_game.getFeature(Otc::GameTibia12Protocol)) {
	//     msg->getU16();
	// }

	// for (int stackPos = 0; stackPos < 256; stackPos++) {
	//     if (msg->peekU16() >= 0xff00)
	//         return msg->getU16() & 0xff;

	//     if (!g_game.getFeature(Otc::GameNewCreatureStacking) && stackPos > Tile::MAX_THINGS)
	//         g_logger.traceError(stdext::format("too many things, pos=%s, stackpos=%d", stdext::to_string(position), stackPos));

	//     ThingPtr thing = getThing(msg);
	//     g_map.addThing(thing, position, stackPos);
	// }

	var x uint16
	for stackPos := 0; stackPos < 256; stackPos++ {
		if err := read(r, &x); err != nil {
			return nil, nil, err
		}

		if x >= 0xff00 {
			return &MapLeftRow{}, p.Next(cur(r)), nil
		}

		if err := skip(r, -2); err != nil {
			return nil, nil, err
		}

		if err := thing(r); err != nil {
			return nil, nil, err
		}
	}

	return &MapLeftRow{}, p.Next(cur(r)), nil
}
