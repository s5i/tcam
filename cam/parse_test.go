package cam

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/s5i/tcam/dat"
	"github.com/s5i/tcam/data"
)

type camFixture struct {
	name        string
	cam         []byte
	dat         *dat.File
	readGolden  []byte
	parseGolden []byte
}

func TestRead(t *testing.T) {
	for _, fx := range camFixtures() {
		t.Run(fx.name, func(t *testing.T) {
			r := bytes.NewReader(fx.cam)
			w := bytes.NewBuffer(nil)

			for packet, err := range Read(r) {
				if err != nil {
					t.Fatalf("Read() error: %v", err)
				}
				fmt.Fprintf(w, "o:%d t:%d l:%d\n", packet.FileOffset, packet.TimeOffset/time.Second, len(packet.Data))
			}

			out := w.Bytes()
			if diff := cmp.Diff(string(fx.readGolden), string(out)); diff != "" {
				if *updateGolden {
					goldenF := fmt.Sprintf("testdata/%s.read.golden.txt", fx.name)
					if err := os.WriteFile(goldenF, out, 0644); err != nil {
						t.Fatalf("os.WriteFile(%q) error when updating golden: %v", goldenF, err)
					}
					t.Logf("Updated %q.", goldenF)
					return
				}
				t.Errorf("Output diff; -want +got:\n%v", diff)
			}
		})
	}
}

func TestParse(t *testing.T) {
	for _, fx := range camFixtures() {
		t.Run(fx.name, func(t *testing.T) {
			r := bytes.NewReader(fx.cam)
			w := bytes.NewBuffer(nil)

			for op, err := range Parse(r, &ParseOpts{DATFile: fx.dat}) {
				if err != nil {
					t.Fatalf("Parse() error: %v", err)
				}
				switch op := op.(type) {
				case data.CamMetadata:
				default:
					t := time.Duration(reflect.ValueOf(op).FieldByName("TimeOffset").Int()).Truncate(time.Second)
					x := reflect.ValueOf(op).FieldByName("PlayerPos").FieldByName("X").Int()
					y := reflect.ValueOf(op).FieldByName("PlayerPos").FieldByName("Y").Int()
					z := reflect.ValueOf(op).FieldByName("PlayerPos").FieldByName("Z").Int()
					n := reflect.TypeOf(op).Name()
					fmt.Fprintf(w, "%v - (%d,%d,%d) - %s\n", t, x, y, z, n)
				}
			}

			out := w.Bytes()
			if diff := cmp.Diff(string(fx.parseGolden), string(out)); diff != "" {
				if *updateGolden {
					goldenF := fmt.Sprintf("testdata/%s.parse.golden.txt", fx.name)
					if err := os.WriteFile(goldenF, out, 0644); err != nil {
						t.Fatalf("os.WriteFile(%q) error when updating golden: %v", goldenF, err)
					}
					t.Logf("Updated %q.", goldenF)
					return
				}
				t.Errorf("Output diff; -want +got:\n%v", diff)
			}
		})
	}
}

func TestParse_CamMetadata(t *testing.T) {
	tests := []struct {
		name       string
		cam        []byte
		dat        *dat.File
		duration   time.Duration
		playerName string
		serverName string
		lastVisit  time.Time
	}{
		{
			name:       "tibiantis",
			cam:        tibiantisCam,
			dat:        tibiantisDAT,
			duration:   1635234 * time.Millisecond,
			playerName: "Shy Teddy",
			serverName: "Tibiantis",
			lastVisit:  time.Date(2025, 11, 28, 14, 37, 22, 0, time.FixedZone("CET", 3600)),
		},
		{
			name:       "relic",
			cam:        relicCam,
			dat:        tibiaRelicDAT,
			duration:   1658906 * time.Millisecond,
			playerName: "Golden",
			serverName: "Tibia Relic",
			lastVisit:  time.Date(2026, 4, 17, 11, 58, 40, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.cam)

			var meta data.CamMetadata
			for op, err := range Parse(r, &ParseOpts{
				DATFile: tt.dat,
				TFilter: map[data.OpType]bool{
					data.TCamMetadata: true,
				},
			}) {
				if err != nil {
					t.Fatalf("Parse() error: %v", err)
				}
				if m, ok := op.(data.CamMetadata); ok {
					meta = m
				}
			}

			if got, want := meta.Duration, tt.duration; got != want {
				t.Fatalf("CamMetadata.Duration = %v, want %v", got, want)
			}

			if got, want := meta.PlayerName, tt.playerName; got != want {
				t.Fatalf("CamMetadata.PlayerName = %q, want %q", got, want)
			}

			if got, want := meta.ServerName, tt.serverName; got != want {
				t.Fatalf("CamMetadata.ServerName = %q, want %q", got, want)
			}

			if got, want := meta.LastVisit, tt.lastVisit; !got.Equal(want) {
				t.Fatalf("CamMetadata.LastVisit = %v, want %v", got, want)
			}
		})
	}
}

func TestParse_MissingDat(t *testing.T) {
	r := bytes.NewReader(tibiantisCam)

	for _, err := range Parse(r, nil) {
		if err == nil {
			t.Fatal("Parse() error = nil, want error")
		}
		return
	}
}

func BenchmarkParse(b *testing.B) {
	for b.Loop() {
		r := bytes.NewReader(tibiantisCam)
		for _, err := range Parse(r, testParseOpts()) {
			if err != nil {
				b.Fatalf("Parse() error: %v", err)
			}
		}
	}
}

func BenchmarkParseIgnore(b *testing.B) {
	for b.Loop() {
		r := bytes.NewReader(tibiantisCam)
		for _, err := range Parse(r, &ParseOpts{
			TFilter: map[data.OpType]bool{},
			DATFile: tibiantisDAT,
		}) {
			if err != nil {
				b.Fatalf("Parse() error: %v", err)
			}
		}
	}
}

func BenchmarkRead(b *testing.B) {
	for b.Loop() {
		r := bytes.NewReader(tibiantisCam)
		for _, err := range Read(r) {
			if err != nil {
				b.Fatalf("Read() error: %v", err)
			}
		}
	}
}
