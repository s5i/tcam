package parser

import (
	"context"
	"flag"
	"io"
	"strings"

	"github.com/s5i/tcam/enum"
	"github.com/s5i/tcam/network"
)

var (
	naiveTalkSearch = flag.Bool("naive_talk_search", false, "Whether to perform a naive talk search in unparsed packets trailings.")
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
				ret, pkt, err = parseTalk(pkt, false)
			case enum.OpCodeMoveCreature:
				ret, pkt, err = parseMoveCreature(pkt)
			case enum.OpCodeChangeOnMap:
				ret, pkt, err = parseChangeOnMap(pkt)
			case enum.OpCodePlayerData:
				ret, pkt, err = parsePlayerData(pkt)
			case enum.OpCodeMarkCreature:
				ret, pkt, err = parseMarkCreature(pkt)
			case enum.OpCodeCreatureHealth:
				ret, pkt, err = parseCreatureHealth(pkt)

			// Broken.
			// case enum.OpCodeMapLeftRow:
			// ret, pkt, err = parseMapLeftRow(pkt)
			default:
				if *naiveTalkSearch {
					basePkt := pkt
					for {
						basePkt = basePkt.Next(1)
						pkt = basePkt

						if pkt == nil {
							break
						}
						if pkt.OpCode() == enum.OpCodeTalk {
							ret, pkt, err = parseTalk(pkt, true)
							if err != nil {
								ret, pkt, err = nil, oldPkt.Next(1), nil
							} else {
								Logger.Printf("Found talk: %v", ret)
							}
						}
					}
				} else {
					ret, pkt, err = pkt.OpCode(), nil, nil
				}
			}

			if pkt != nil && strings.HasPrefix(pkt.OpCode().String(), "Unknown") {
				Logger.Printf("BAD PACKET - %v", oldPkt)
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
