package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/s5i/tcam/cam"
	"github.com/s5i/tcam/data"
	"golang.org/x/sync/errgroup"
)

// Creature prints timestamps and locations when a creature was spotted.
func Creature(ctx context.Context, dirPath string, w io.Writer, name string) error {
	type result struct {
		path string
		text []byte
	}

	var results []result
	var resultsMu sync.Mutex

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

			who := map[uint32]string{}
			where := map[uint32]data.Location{}
			when := map[uint32]time.Duration{}

			for op, err := range cam.Parse(f, &cam.ParseOpts{
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
					return err
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

						if _, ok := where[th.Creature.ID]; ok {
							continue
						}

						if name != "*" && !strings.EqualFold(name, th.Creature.Name) {
							continue
						}

						who[th.Creature.ID] = th.Creature.Name
						where[th.Creature.ID] = t.Location
						when[th.Creature.ID] = ts
					}
				}
			}

			if len(where) == 0 {
				return nil
			}

			b := bytes.NewBuffer(nil)
			fmt.Fprintf(b, "## %s\n\n", filepath.Base(path))

			order := slices.SortedFunc(maps.Keys(when), func(a, b uint32) int {
				return int(when[a] - when[b])
			})
			for _, id := range order {
				fmt.Fprintf(b, "* %s @ (%d, %d, %d) at %v\n", who[id], where[id].X, where[id].Y, where[id].Z, when[id].Truncate(time.Millisecond))
			}

			resultsMu.Lock()
			results = append(results, result{
				path: path,
				text: b.Bytes(),
			})
			resultsMu.Unlock()

			return nil
		})

		return nil
	}); err != nil {
		return err
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	slices.SortFunc(results, func(a, b result) int {
		return strings.Compare(a.path, b.path)
	})

	if name == "*" {
		fmt.Fprintf(w, "# All creatures\n")
	} else {
		fmt.Fprintf(w, "# Creature %s\n", name)
	}

	for _, r := range results {
		fmt.Fprintf(w, "\n%s", r.text)
	}

	return nil
}
