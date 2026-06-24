package cam

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"time"

	"github.com/s5i/tcam/dat"
	"github.com/s5i/tcam/data"
)

var errMissingDat = errors.New("opts.Dat is required")

// ParseOpts controls the behavior of Parse.
type ParseOpts struct {
	// DATFile holds item metadata from a Tibia client .dat file.
	DATFile *dat.File

	// If set, only yield the specified operation types.
	TFilter map[data.OpType]bool

	// If set, Parse will populate the maps.
	Stats *ParseStats
}

// Parse returns an iterator over the provided io.ReadSeeker that returns subsequent data.Operations.
func Parse(r io.ReadSeeker, opts *ParseOpts) iter.Seq2[data.Operation, error] {
	return func(yield func(data.Operation, error) bool) {
		if opts == nil || opts.DATFile == nil {
			yield(nil, errMissingDat)
			return
		}

		yieldVal := func(p data.Operation) bool { return yield(p, nil) }
		yieldErr := func(err error) {
			if !errors.Is(err, io.EOF) {
				yield(nil, err)
			}
		}

		state := &parseState{
			stats: opts.Stats,
		}

		var finalTimeOffset time.Duration
		for packet, err := range Read(r) {
			if err != nil {
				yieldErr(err)
				return
			}

			finalTimeOffset = packet.TimeOffset

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

		if opts.TFilter == nil || opts.TFilter[data.TCamMetadata] {
			if !yieldVal(data.CamMetadata{
				Duration:   finalTimeOffset,
				PlayerName: state.playerName,
			}) {
				return
			}
		}
	}
}
