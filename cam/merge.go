package cam

import (
	"encoding/binary"
	"errors"
	"io"
)

type MergeOpts struct{}

// Merge merges multiple CAM files into a single file.
// Still needs work -- we're seeing crashes on the end of the second cam.
func Merge(w io.WriteSeeker, opts *MergeOpts, r ...io.ReadSeeker) error {
	if len(r) == 0 {
		return errors.New("no CAM files provided")
	}

	// Clone the header from the first file.
	header := make([]byte, 4)
	if _, err := io.ReadFull(r[0], header); err != nil {
		return err
	}
	if _, err := w.Write(header); err != nil {
		return err
	}

	if _, err := r[0].Seek(0, io.SeekStart); err != nil {
		return err
	}

	// Placeholder for checksum.
	if _, err := w.Write(make([]byte, 8)); err != nil {
		return err
	}

	tick := uint64(0)
	for _, r := range r {
		tickOffset := tick
		for packet, err := range Read(r) {
			tick = tickOffset + uint64(packet.TimeOffset.Milliseconds())
			if err != nil {
				return err
			}
			if err := binary.Write(w, binary.LittleEndian, tick); err != nil {
				return err
			}
			if err := binary.Write(w, binary.LittleEndian, uint16(len(packet.Data))); err != nil {
				return err
			}
			if _, err := w.Write(packet.Data); err != nil {
				return err
			}
		}
	}

	// Skip to the start of the checksum.
	if _, err := w.Seek(4, io.SeekStart); err != nil {
		return err
	}

	// TODO: Calculate and write the checksum.
	// checksum := make([]byte, 8)
	// w.Write(checksum)

	return nil
}
