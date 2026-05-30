package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/s5i/tcam/cam"
	"github.com/s5i/tcam/data"
	"golang.org/x/sync/errgroup"
)

func DialogueTree(ctx context.Context, dirPath string, w io.Writer) error {
	type message struct {
		Name string
		Text string
		NPC  bool
	}

	merged := map[message]map[message]bool{}
	var mergedMu sync.Mutex

	eg, ctx := errgroup.WithContext(ctx)
	if err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		eg.Go(func() error {
			if d.IsDir() {
				return nil
			}

			if strings.ToLower(filepath.Ext(path)) != ".cam" {
				return nil
			}

			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			var messages []message
			players := map[string]bool{}

			for op, err := range cam.Parse(f, &cam.ParseOpts{
				TFilter: map[data.OpType]bool{
					data.TMap:             true,
					data.TMoveNorth:       true,
					data.TMoveSouth:       true,
					data.TMoveEast:        true,
					data.TMoveWest:        true,
					data.TMoveFloorDown:   true,
					data.TMoveFloorUp:     true,
					data.TTileItemAdd:     true,
					data.TCreatureMessage: true,
				},
			}) {
				if err != nil {
					return err
				}

				var tiles []data.Tile
				switch msg := op.(type) {
				case data.Map:
					tiles = msg.Tiles
				case data.MoveNorth:
					tiles = msg.Tiles
				case data.MoveSouth:
					tiles = msg.Tiles
				case data.MoveEast:
					tiles = msg.Tiles
				case data.MoveWest:
					tiles = msg.Tiles
				case data.MoveFloorDown:
					tiles = msg.Tiles
				case data.MoveFloorUp:
					tiles = msg.Tiles
				case data.TileItemAdd:
					tiles = []data.Tile{{Location: msg.Location, Things: []data.Thing{msg.Thing}}}
				case data.CreatureMessage:
					if msg.Type != 1 {
						continue
					}

					npc := !players[msg.Name]
					text := msg.Text
					text = numberRE.ReplaceAllString(text, "[num]")
					if npc {
						for player := range players {
							text = strings.ReplaceAll(text, player, "[player]")
						}
					}

					messages = append(messages, message{
						Name: msg.Name,
						Text: text,
						NPC:  npc,
					})
				}

				for _, t := range tiles {
					for _, th := range t.Things {
						if !th.HasCreature {
							continue
						}
						if th.Creature.Name == "" {
							continue
						}
						if th.Creature.ID >= 1<<30 {
							continue
						}
						players[th.Creature.Name] = true
					}
				}
			}

			mergedMu.Lock()
			for i, msg := range messages {
				if i == 0 {
					continue
				}
				if !msg.NPC {
					continue
				}
				if merged[msg] == nil {
					merged[msg] = map[message]bool{}
				}
				prev := messages[i-1]
				if !prev.NPC {
					prev.Name = "Player"
				}
				merged[msg][prev] = true
			}
			mergedMu.Unlock()

			return nil
		})

		return nil
	}); err != nil {
		return err
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	for _, resp := range slices.SortedFunc(maps.Keys(merged), func(a, b message) int {
		return strings.Compare(a.Text, b.Text)
	}) {

		if resp.Name != "Gabel" {
			continue
		}
		reqs := merged[resp]
		fmt.Fprintf(w, "%s: %s\n", resp.Name, resp.Text)
		for _, req := range slices.SortedFunc(maps.Keys(reqs), func(a, b message) int {
			return strings.Compare(a.Text, b.Text)
		}) {
			fmt.Fprintf(w, "\t%s\n", req.Text)
		}
	}

	return nil
}

var numberRE = regexp.MustCompile(`\d+`)
