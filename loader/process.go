package loader

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/s5i/tcam/enum"
	"github.com/s5i/tcam/network"
)

func ReadFile(ctx context.Context, path string) (<-chan network.Packet, <-chan error) {
	retCh := make(chan network.Packet)
	errCh := make(chan error, 1)

	go func() (retErr error) {
		defer func() {
			if retErr != nil && retErr != io.EOF {
				errCh <- retErr
			} else {
				close(retCh)
			}
		}()

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		// Skip header (4 bytes) + checksum (8 bytes).
		if _, err := f.Seek(4+8, io.SeekCurrent); err != nil {
			return err
		}

		// Read first tick count (8 bytes), rewind.
		var startTick uint64
		if err := binary.Read(f, binary.LittleEndian, &startTick); err != nil {
			return err
		}
		if _, err := f.Seek(-8, io.SeekCurrent); err != nil {
			return err
		}

		for {
			cur, _ := f.Seek(0, io.SeekCurrent)
			pos := fmt.Sprintf("%8X+%2d", cur-cur%16, cur%16)

			var ticks uint64
			if err := binary.Read(f, binary.LittleEndian, &ticks); err != nil {
				return err
			}

			var pktLen uint16
			if err := binary.Read(f, binary.LittleEndian, &pktLen); err != nil {
				return err
			}

			packetData := make([]byte, pktLen)
			if _, err := io.ReadFull(f, packetData); err != nil {
				return err
			}

			opCode := enum.OpCode(packetData[0])

			Logger.Printf("%s | OP = %s | LEN = %d | T = %v", pos, opCode, pktLen, time.Duration(ticks-startTick)*time.Millisecond)

			select {
			case <-ctx.Done():
				return ctx.Err()

			case retCh <- network.Packet{
				Offset: time.Duration(ticks-startTick) * time.Millisecond,
				Data:   packetData,
			}:
			}

		}
	}()

	return retCh, errCh
}
