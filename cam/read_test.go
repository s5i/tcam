package cam

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestRead(t *testing.T) {
	r := bytes.NewReader(input)
	w := bytes.NewBuffer(nil)

	for packet, err := range Read(r) {
		if err != nil {
			t.Fatalf("Read() error: %v", err)
		}
		fmt.Fprintf(w, "o:%d t:%d l:%d\n", packet.FileOffset, packet.TimeOffset/time.Second, len(packet.Data))
	}

	out := w.Bytes()
	if diff := cmp.Diff(string(golden), string(out)); diff != "" {
		if *updateGolden {
			goldenF := "testdata/1.golden.txt"
			if err := os.WriteFile(goldenF, out, 0644); err != nil {
				t.Fatalf("os.WriteFile(%q) error when updating golden: %v", goldenF, err)
			}
			t.Logf("Updated %q.", goldenF)
			return
		}
		t.Errorf("Output diff; -want +got:\n%v", diff)
	}
}

var (
	updateGolden = flag.Bool("cam_update_golden", false, "Whether to update the golden files.")

	//go:embed testdata/1.cam
	input []byte
	//go:embed testdata/1.golden.txt
	golden []byte
)
