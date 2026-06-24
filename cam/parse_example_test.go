package cam

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/s5i/tcam/data"
)

func ExampleParse_dialogues() {
	r := bytes.NewReader(input)
	target := "Muzir"
	ctxSize := time.Minute

	var messages []data.CreatureMessage
	for op, err := range Parse(r, &ParseOpts{
		DATFile: testDat,
		TFilter: map[data.OpType]bool{
			data.TCreatureMessage: true,
		},
	}) {
		if err != nil {
			panic(err)
		}

		if msg, ok := op.(data.CreatureMessage); ok && msg.Type == 1 {
			messages = append(messages, msg)
		}
	}

	target = strings.ToLower(target)
	var seen bool
	var lastSeen time.Duration

	show := make([]bool, len(messages))
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]

		if strings.ToLower(msg.Name) == target {
			show[i] = true
			seen = true
			lastSeen = msg.TimeOffset
			continue
		}

		if seen && lastSeen < msg.TimeOffset+ctxSize {
			show[i] = true
		}
	}

	for i, msg := range messages {
		if show[i] {
			fmt.Printf("%s: %s\n", msg.Name, msg.Text)
		}
	}

	// Output:
	// Shy Teddy: hi
	// Muzir: Welcome Shy Teddy! Daraman's blessings.
	// Shy Teddy: change gold
	// Muzir: How many platinum coins do you want to get?
	// Shy Teddy: 32
	// Shy Teddy: yes
	// Muzir: Here you are.
	// Muzir: Daraman's blessings.
}

func ExampleParse_location() {
	r := bytes.NewReader(input)
	x, y, z := 33175, 32524, 7
	radius := 7
	ctxSize := time.Minute

	dist := func(a, b int) int {
		if a < b {
			return b - a
		}
		return a - b
	}

	var seen []time.Duration
	for op, err := range Parse(r, &ParseOpts{
		DATFile: testDat,
		TFilter: map[data.OpType]bool{
			data.TMoveNorth:     true,
			data.TMoveSouth:     true,
			data.TMoveEast:      true,
			data.TMoveWest:      true,
			data.TMoveFloorDown: true,
			data.TMoveFloorUp:   true,
		},
	}) {
		if err != nil {
			panic(err)
		}

		var loc data.Location
		var ts time.Duration

		switch msg := op.(type) {
		case data.MoveNorth:
			loc = msg.PlayerPos
			ts = msg.TimeOffset
		case data.MoveSouth:
			loc = msg.PlayerPos
			ts = msg.TimeOffset
		case data.MoveEast:
			loc = msg.PlayerPos
			ts = msg.TimeOffset
		case data.MoveWest:
			loc = msg.PlayerPos
			ts = msg.TimeOffset
		case data.MoveFloorDown:
			loc = msg.PlayerPos
			ts = msg.TimeOffset
		case data.MoveFloorUp:
			loc = msg.PlayerPos
			ts = msg.TimeOffset
		}

		if loc.Z != z || dist(loc.X, x) > radius || dist(loc.Y, y) > radius {
			continue
		}

		if len(seen) > 0 && seen[len(seen)-1]+ctxSize >= ts {
			continue
		}

		seen = append(seen, ts)
	}

	for _, ts := range seen {
		fmt.Println(ts.Truncate(time.Second))
	}

	// Output:
	// 52s
	// 3m23s
}

func ExampleParse_creature() {
	r := bytes.NewReader(input)
	name := "Muzir"

	for op, err := range Parse(r, &ParseOpts{
		DATFile: testDat,
		TFilter: map[data.OpType]bool{
			data.TMap:           true,
			data.TMoveNorth:     true,
			data.TMoveSouth:     true,
			data.TMoveEast:      true,
			data.TMoveWest:      true,
			data.TMoveFloorDown: true,
			data.TMoveFloorUp:   true,
			data.TTileItemAdd:   true,
		},
	}) {
		if err != nil {
			panic(err)
		}

		var ts time.Duration
		var tiles []data.Tile

		switch msg := op.(type) {
		case data.Map:
			ts = msg.TimeOffset
			tiles = msg.Tiles
		case data.MoveNorth:
			ts = msg.TimeOffset
			tiles = msg.Tiles
		case data.MoveSouth:
			ts = msg.TimeOffset
			tiles = msg.Tiles
		case data.MoveEast:
			ts = msg.TimeOffset
			tiles = msg.Tiles
		case data.MoveWest:
			ts = msg.TimeOffset
			tiles = msg.Tiles
		case data.MoveFloorDown:
			ts = msg.TimeOffset
			tiles = msg.Tiles
		case data.MoveFloorUp:
			ts = msg.TimeOffset
			tiles = msg.Tiles
		case data.TileItemAdd:
			ts = msg.TimeOffset
			tiles = []data.Tile{{Location: msg.Location, Things: []data.Thing{msg.Thing}}}
		}

		for _, t := range tiles {
			for _, th := range t.Things {
				if !th.HasCreature {
					continue
				}

				if !strings.EqualFold(name, th.Creature.Name) {
					continue
				}

				fmt.Printf("%s @ (%d, %d, %d) at %v\n", th.Creature.Name, t.Location.X, t.Location.Y, t.Location.Z, ts.Truncate(time.Millisecond))
			}
		}
	}

	// Output:
	// Muzir @ (33221, 32389, 7) at 2m2.109s
}
