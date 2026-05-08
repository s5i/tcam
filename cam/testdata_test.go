package cam

import (
	"flag"

	_ "embed"
)

var (
	updateGolden = flag.Bool("cam_update_golden", false, "Whether to update the golden files.")

	//go:embed testdata/1.cam
	input []byte
	//go:embed testdata/1.read.golden.txt
	golden []byte

	//go:embed testdata/1.parse.golden.txt
	parseGolden []byte
)
