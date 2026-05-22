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

// PlayerName returns the CAM player name.
func PlayerName(r io.ReadSeeker) (string, error) {
	cursor, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return "", err
	}

	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	name, ok := func() (string, bool) {
		var id uint32
		for op, err := range Parse(r, &ParseOpts{
			TFilter: map[data.OpType]bool{
				data.TMap:              true,
				data.TLoginPlayerState: true,
			},
		}) {
			if err != nil {
				return "", false
			}

			switch msg := op.(type) {
			case data.LoginPlayerState:
				id = msg.PlayerID
			case data.Map:
				if id == 0 {
					continue
				}

				for _, tile := range msg.Tiles {
					for _, t := range tile.Things {
						if !t.HasCreature || t.Creature.ID != id {
							continue
						}
						return t.Creature.Name, true
					}
				}
			}
		}

		return "", false
	}()
	if !ok {
		return "", fmt.Errorf("name not found")
	}

	if _, err := r.Seek(cursor, io.SeekStart); err != nil {
		return "", err
	}

	return name, nil
}
