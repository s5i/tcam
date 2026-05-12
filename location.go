package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
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

// Location prints timestamps when a location was spotted.
func Location(ctx context.Context, dirPath string, w io.Writer, x, y, z, radius int, ctxSize time.Duration) error {
	dist := func(a, b int) int {
		if a < b {
			return b - a
		}
		return a - b
	}

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

			var seen []time.Duration
			for op, err := range cam.Parse(f, &cam.ParseOpts{
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
					return err
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

			if len(seen) == 0 {
				return nil
			}

			b := bytes.NewBuffer(nil)
			fmt.Fprintf(b, "## %s\n\n", filepath.Base(path))

			for _, ts := range seen {
				fmt.Fprintf(b, "* %v\n", ts.Truncate(time.Second))
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

	fmt.Fprintf(w, "# Location (x=%d, y=%d, z=%d, r=%d)\n", x, y, z, radius)
	for _, r := range results {
		fmt.Fprintf(w, "\n%s", r.text)
	}

	return nil
}
