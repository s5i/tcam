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

func Dialogues(ctx context.Context, dirPath string, w io.Writer, target string, ctxSize time.Duration) error {
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

			var messages []data.CreatureMessage
			for op, err := range cam.Parse(f, nil) {
				if err != nil {
					return err
				}

				if msg, ok := op.(data.CreatureMessage); ok {
					if msg.Type != 1 {
						continue
					}

					messages = append(messages, msg)
				}
			}

			target := strings.ToLower(target)
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

			if !seen {
				return nil
			}

			b := bytes.NewBuffer(nil)
			fmt.Fprintf(b, "## %s\n", filepath.Base(path))

			var hasPrev bool
			var lastMsg time.Duration
			for i, msg := range messages {
				if show[i] {
					if !hasPrev || lastMsg+ctxSize < msg.TimeOffset {
						fmt.Fprintf(b, "\n### %v\n\n", msg.TimeOffset)
					}

					hasPrev = true
					lastMsg = msg.TimeOffset
					fmt.Fprintf(b, "* %s: %s\n", msg.Name, msg.Text)
				}
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

	fmt.Fprintf(w, "# Dialogues for %s\n", target)
	for _, r := range results {
		fmt.Fprintf(w, "\n%s", r.text)
	}

	return nil
}
