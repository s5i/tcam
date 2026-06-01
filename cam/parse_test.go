package cam

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/s5i/tcam/data"
)

func TestParse(t *testing.T) {
	r := bytes.NewReader(input)
	w := bytes.NewBuffer(nil)

	for op, err := range Parse(r, nil) {
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
	if diff := cmp.Diff(string(parseGolden), string(out)); diff != "" {
		if *updateGolden {
			goldenF := "testdata/1.parse.golden.txt"
			if err := os.WriteFile(goldenF, out, 0644); err != nil {
				t.Fatalf("os.WriteFile(%q) error when updating golden: %v", goldenF, err)
			}
			t.Logf("Updated %q.", goldenF)
			return
		}
		t.Errorf("Output diff; -want +got:\n%v", diff)
	}
}

func BenchmarkParse(b *testing.B) {
	for b.Loop() {
		r := bytes.NewReader(input)
		for _, err := range Parse(r, nil) {
			if err != nil {
				b.Fatalf("Parse() error: %v", err)
			}
		}
	}
}

func BenchmarkParseIgnore(b *testing.B) {
	for b.Loop() {
		r := bytes.NewReader(input)
		for _, err := range Parse(r, &ParseOpts{
			TFilter: map[data.OpType]bool{},
		}) {
			if err != nil {
				b.Fatalf("Parse() error: %v", err)
			}
		}
	}
}
