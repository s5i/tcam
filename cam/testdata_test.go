package cam

import (
	"bytes"
	"flag"

	"github.com/s5i/tcam/dat"

	_ "embed"
)

var (
	updateGolden = flag.Bool("cam_update_golden", false, "Whether to update the golden files.")

	//go:embed testdata/tibiantis.cam
	tibiantisCam []byte
	//go:embed testdata/tibiantis.read.golden.txt
	tibiantisReadGolden []byte
	//go:embed testdata/tibiantis.parse.golden.txt
	tibiantisParseGolden []byte

	//go:embed testdata/relic.cam
	relicCam []byte
	//go:embed testdata/relic.read.golden.txt
	relicReadGolden []byte
	//go:embed testdata/relic.parse.golden.txt
	relicParseGolden []byte

	//go:embed testdata/Tibiantis.dat
	tibiantisDat []byte
	//go:embed testdata/TibiaRelic.dat
	tibiaRelicDat []byte

	tibiantisDAT *dat.File
	tibiaRelicDAT *dat.File
)

func init() {
	var err error
	tibiantisDAT, err = dat.Read(bytes.NewReader(tibiantisDat))
	if err != nil {
		panic(err)
	}
	tibiaRelicDAT, err = dat.Read(bytes.NewReader(tibiaRelicDat))
	if err != nil {
		panic(err)
	}
}

func testParseOpts() *ParseOpts {
	return &ParseOpts{DATFile: tibiantisDAT}
}

func testMergeOpts() *MergeOpts {
	return &MergeOpts{Dat: tibiantisDAT}
}

func camFixtures() []camFixture {
	return []camFixture{
		{
			name:        "tibiantis",
			cam:         tibiantisCam,
			dat:         tibiantisDAT,
			readGolden:  tibiantisReadGolden,
			parseGolden: tibiantisParseGolden,
		},
		{
			name:        "relic",
			cam:         relicCam,
			dat:         tibiaRelicDAT,
			readGolden:  relicReadGolden,
			parseGolden: relicParseGolden,
		},
	}
}
