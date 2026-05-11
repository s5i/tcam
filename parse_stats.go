package main

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/s5i/tcam/cam"
	"github.com/s5i/tcam/data"
	"golang.org/x/sync/errgroup"
)

// ParseStats prints aggregate parsing performance statistics for a directory.
func ParseStats(ctx context.Context, dirPath string, w io.Writer, noFilter bool) error {
	totalStats := cam.NewParseStats()
	var totalMu sync.Mutex

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(runtime.NumCPU())
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

			opts := &cam.ParseOpts{
				Stats: &cam.ParseStats{
					Count:    map[data.OpType]int{},
					Duration: map[data.OpType]time.Duration{},
				},
			}
			if !noFilter {
				opts.TFilter = map[data.OpType]bool{}
			}
			for _, err := range cam.Parse(f, opts) {
				if err != nil {
					return err
				}
			}

			totalMu.Lock()
			totalStats.Merge(opts.Stats)
			totalMu.Unlock()

			return nil
		})

		return nil
	}); err != nil {
		return err
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	totalStats.Write(os.Stderr)

	return nil
}
