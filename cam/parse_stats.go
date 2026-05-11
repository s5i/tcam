package cam

import (
	"fmt"
	"io"
	"maps"
	"slices"
	"time"

	"github.com/s5i/tcam/data"
)

// NewParseStats initializes ParseStats.
func NewParseStats() *ParseStats {
	return &ParseStats{
		Count:    map[data.OpType]int{},
		Duration: map[data.OpType]time.Duration{},
	}
}

// ParseStats contains parsing statistics.
type ParseStats struct {
	Count    map[data.OpType]int
	Duration map[data.OpType]time.Duration
}

// Merge adds data from the argument into the receiver.
func (s *ParseStats) Merge(from *ParseStats) {
	for op := range from.Count {
		s.Count[op] += from.Count[op]
		s.Duration[op] += from.Duration[op]
	}
}

// Write outputs human-readable ParseStats.
func (s *ParseStats) Write(w io.Writer) {
	for _, op := range slices.SortedFunc(maps.Keys(s.Count), func(a, b data.OpType) int {
		return int(s.Duration[b]/time.Duration(s.Count[b]+1) - s.Duration[a]/time.Duration(s.Count[a]+1))
	}) {
		fmt.Fprintf(w, "%s -> %v / op (%d ops, %v elapsed)\n", data.OpName[op], s.Duration[op]/time.Duration(s.Count[op]), s.Count[op], s.Duration[op].Truncate(time.Microsecond))
	}
}
