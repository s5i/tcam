package cam

import (
	"errors"
	"fmt"
	"io"
	"iter"

	"github.com/s5i/tcam/data"
)

// ParseOpts controls the behavior of Parse.
type ParseOpts struct {
	// If set, only yield the specified operation types.
	TFilter map[data.OpType]bool

	// If set, Parse will populate the maps.
	Stats *ParseStats
}

// Parse returns an iterator over the provided io.ReadSeeker that returns subsequent data.Operations.
func Parse(r io.ReadSeeker, opts *ParseOpts) iter.Seq2[data.Operation, error] {
	return func(yield func(data.Operation, error) bool) {
		if opts == nil {
			opts = &ParseOpts{}
		}

		yieldVal := func(p data.Operation) bool { return yield(p, nil) }
		yieldErr := func(err error) {
			if !errors.Is(err, io.EOF) {
				yield(nil, err)
			}
		}

		state := &parseState{
			tiles: make(map[tileKey][]data.Thing),
			stats: opts.Stats,
		}

		for packet, err := range Read(r) {
			if err != nil {
				yieldErr(err)
				return
			}

			ops, err := parsePacket(state, packet.Data, packet.TimeOffset, opts)
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
