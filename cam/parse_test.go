package cam

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	r := bytes.NewReader(input)
	w := bytes.NewBuffer(nil)

	for op, err := range Parse(r) {
		if err != nil {
			t.Fatalf("Parse() error: %v", err)
		}
		fmt.Fprintf(w, "%s\n", reflect.TypeOf(op).Name())
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
		for _, err := range Parse(r) {
			if err != nil {
				b.Fatalf("Parse() error: %v", err)
			}
		}
	}
}
