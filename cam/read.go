package cam

import (
	"encoding/binary"
	"errors"
	"io"
	"iter"
	"time"

	"github.com/s5i/tcam/data"
)

// Read returns an iterator over the provided io.ReadSeeker that returns subsequent data.RawPackets.
func Read(r io.ReadSeeker) iter.Seq2[data.RawPacket, error] {
		return func(yield func(data.RawPacket, error) bool) {
		yieldVal := func(p data.RawPacket) bool { return yield(p, nil) }
		yieldErr := func(err error) {
			if !errors.Is(err, io.EOF) {
				yield(data.RawPacket{}, err)
			}
		}

		var headerSize uint32
		if err := binary.Read(r, binary.LittleEndian, &headerSize); err != nil {
			yieldErr(err)
			return
		}
		dataOffset := int64(headerSize) + 4
		if _, err := r.Seek(dataOffset, io.SeekStart); err != nil {
			yieldErr(err)
			return
		}

		// Read first tick count (8 bytes), rewind.
		var startTick uint64
		if err := binary.Read(r, binary.LittleEndian, &startTick); err != nil {
			yieldErr(err)
			return
		}
		if _, err := r.Seek(-8, io.SeekCurrent); err != nil {
			yieldErr(err)
			return
		}

		for {
			// Read tick count (8 bytes).
			var curTick uint64
			if err := binary.Read(r, binary.LittleEndian, &curTick); err != nil {
				yieldErr(err)
				return
			}

			// Read packet length (2 bytes).
			var pktLen uint16
			if err := binary.Read(r, binary.LittleEndian, &pktLen); err != nil {
				yieldErr(err)
				return
			}

			cur, _ := r.Seek(0, io.SeekCurrent)
			packetData := make([]byte, pktLen)
			if _, err := io.ReadFull(r, packetData); err != nil {
				yieldErr(err)
				return
			}

			if !yieldVal(data.RawPacket{
				FileOffset: int(cur),
				TimeOffset: time.Duration(curTick-startTick) * time.Millisecond,
				Data:       packetData,
			}) {
				return
			}
		}
	}
}
