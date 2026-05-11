package cam

import (
	"errors"
	"fmt"
	"io"
	"iter"

	"github.com/s5i/tcam/data"
)

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
