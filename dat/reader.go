package dat

import (
	"encoding/binary"
	"fmt"
	"io"
)

type binReader struct {
	r io.Reader
}

func (br *binReader) u8() (byte, error) {
	var b [1]byte
	if _, err := io.ReadFull(br.r, b[:]); err != nil {
		return 0, err
	}
	return b[0], nil
}

func (br *binReader) u16() (uint16, error) {
	var b [2]byte
	if _, err := io.ReadFull(br.r, b[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b[:]), nil
}

func (br *binReader) u32() (uint32, error) {
	var b [4]byte
	if _, err := io.ReadFull(br.r, b[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b[:]), nil
}

func (br *binReader) skip(n int) error {
	if n <= 0 {
		return nil
	}
	_, err := io.CopyN(io.Discard, br.r, int64(n))
	return err
}

const (
	firstItemID = 100
	maxSprites  = 4096
)

type formatVersion int

const (
	formatV1 formatVersion = iota
	formatV2
	formatV3
	formatV4
	formatV5
	formatV6
)

func formatForSignature(signature uint32) formatVersion {
	// TODO(s5i): Map other signatures as they are discovered.
	switch signature {
	case 0x439D5A33, 0x6970EFAD:
		return formatV3
	default:
		return formatV1
	}
}

func readHeader(br *binReader) (signature uint32, itemCount, outfitCount, effectCount, missileCount uint16, err error) {
	signature, err = br.u32()
	if err != nil {
		return
	}
	itemCount, err = br.u16()
	if err != nil {
		return
	}
	outfitCount, err = br.u16()
	if err != nil {
		return
	}
	effectCount, err = br.u16()
	if err != nil {
		return
	}
	missileCount, err = br.u16()
	return
}

func readTexturePatterns(br *binReader, patternZFixed bool) error {
	width, err := br.u8()
	if err != nil {
		return err
	}
	height, err := br.u8()
	if err != nil {
		return err
	}
	if width > 1 || height > 1 {
		if _, err := br.u8(); err != nil {
			return err
		}
	}

	layers, err := br.u8()
	if err != nil {
		return err
	}
	patternX, err := br.u8()
	if err != nil {
		return err
	}
	patternY, err := br.u8()
	if err != nil {
		return err
	}
	var patternZ byte
	if patternZFixed {
		patternZ = 1
	} else {
		patternZ, err = br.u8()
		if err != nil {
			return err
		}
	}
	frames, err := br.u8()
	if err != nil {
		return err
	}

	totalSprites := uint32(width) * uint32(height) * uint32(layers) * uint32(patternX) * uint32(patternY) * uint32(patternZ) * uint32(frames)
	if totalSprites > maxSprites {
		return fmt.Errorf("dat: thing has %d sprites (max %d)", totalSprites, maxSprites)
	}
	return br.skip(int(totalSprites) * 2)
}
