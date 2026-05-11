package data

import "time"

// RawPacket contains the raw network packet and some metadata.
type RawPacket struct {
	FileOffset int
	TimeOffset time.Duration
	Data       []byte
}
