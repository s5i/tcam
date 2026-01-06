package parser

import (
	"context"
	"io"
	"strings"

	"github.com/s5i/tcam/enum"
	"github.com/s5i/tcam/network"
)

func ParsePackets(ctx context.Context, packetsCh <-chan *network.Packet) (<-chan any, <-chan error) {
	retCh := make(chan any)
	errCh := make(chan error, 1)

	go func() (retErr error) {
		defer func() {
			if retErr != nil && retErr != io.EOF {
				errCh <- retErr
			} else {
				close(retCh)
			}
		}()

		var pkt *network.Packet
		var ret any
		var err error
		for {
			if pkt == nil {
				select {
				case <-ctx.Done():
					return ctx.Err()

				case p, ok := <-packetsCh:
					if !ok {
						return nil
					}

					pkt = p
				}
			}

			oldPkt := pkt
			switch pkt.OpCode() {
			case enum.OpCodeTalk:
				ret, pkt, err = parseTalk(pkt)
			case enum.OpCodeMoveCreature:
				ret, pkt, err = parseMoveCreature(pkt)
			case enum.OpCodeChangeOnMap:
				ret, pkt, err = parseChangeOnMap(pkt)
				if pkt != nil {
					if strings.HasPrefix(pkt.OpCode().String(), "Unknown") {
						Logger.Printf("BAD PACKET - %v", oldPkt)
					}
				}
			default:
				ret, pkt, err = pkt.OpCode(), nil, nil
			}

			if err != nil {
				return err
			}

			select {
			case retCh <- ret:
			case <-ctx.Done():
				return ctx.Err()
			}

		}
	}()

	return retCh, errCh
}
