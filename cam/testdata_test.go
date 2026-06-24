package cam

import (
	"bytes"
	"flag"

	"github.com/s5i/tcam/dat"

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

	//go:embed testdata/Tibiantis.dat
	tibiantisDat []byte

	testDat *dat.File
)

func init() {
	var err error
	testDat, err = dat.Read(bytes.NewReader(tibiantisDat))
	if err != nil {
		panic(err)
	}
}

func testParseOpts() *ParseOpts {
	return &ParseOpts{DATFile: testDat}
}

func testMergeOpts() *MergeOpts {
	return &MergeOpts{Dat: testDat}
}
