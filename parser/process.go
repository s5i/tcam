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
	naiveTalkSearch = flag.Bool("naive_talk_search", true, "Whether to perform a naive talk search in unparsed packets trailings.")
)

type UnhandledPacket struct {
	OpCode enum.OpCode
	Packet *network.Packet
}

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

		for {
			var ret any
			var err error

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

			switch opcode := pkt.OpCode(); opcode {
			case enum.OpCodeTalk:
				ret, pkt, err = parseTalk(pkt, false)
			case enum.OpCodeMoveCreature:
				ret, pkt, err = parseMoveCreature(pkt)
			case enum.OpCodePlayerData:
				ret, pkt, err = parsePlayerData(pkt)
			case enum.OpCodeMarkCreature:
				ret, pkt, err = parseMarkCreature(pkt)
			case enum.OpCodeCreatureHealth:
				ret, pkt, err = parseCreatureHealth(pkt)
			case enum.OpCodeLoginOrPendingState:
				ret, pkt, err = parseLoginOrPendingState(pkt)
			// Broken.
			// case enum.OpCodeChangeOnMap:
			// 	ret, pkt, err = parseChangeOnMap(pkt)
			// case enum.OpCodeMapLeftRow:
			// 	ret, pkt, err = parseMapLeftRow(pkt)
			default:
				if *naiveTalkSearch {
					for {
						pkt = pkt.Next(1)

						if pkt == nil {
							break
						}

						if pkt.OpCode() != enum.OpCodeTalk {
							continue
						}

						nRet, nPkt, nErr := parseTalk(pkt, true)
						if nErr != nil {
							continue
						}
						if nPkt != nil && strings.HasPrefix(pkt.OpCode().String(), "Unknown") {
							continue
						}

						ret, pkt, err = nRet, nPkt, nErr
						break
					}

					if ret != nil {
						break
					} else {
						pkt = oldPkt
					}
				}

				ret, pkt, err = &UnhandledPacket{OpCode: pkt.OpCode(), Packet: pkt}, nil, nil
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
