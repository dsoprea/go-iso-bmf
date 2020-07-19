package bmf

import (
	"fmt"
	"testing"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
	"github.com/dsoprea/go-iso-bmf/type"
)

const (
	testMp4Filepath  = "assets/tears-of-steel.mp4"
	testHeicFilepath = "assets/image.heic"
)

func TestOpen_Mp4(t *testing.T) {
	s, err := Open(testMp4Filepath)
	log.PanicIf(err)

	ftypBoxes := bmfcommon.ChildBoxes(s, "ftyp")
	ftyp := ftypBoxes[0].(*bmftype.FtypBox)

	if ftyp.Name() != "ftyp" {
		t.Fatalf("ftyp name not correct: [%s]", ftyp.Name())
	}

	if ftyp.MajorBrand() != "isom" {
		t.Fatalf("ftyp MajorBrand is not correct: [%s]", ftyp.MajorBrand())
	}
}

func TestOpen_Heic(t *testing.T) {
	s, err := Open(testHeicFilepath)
	log.PanicIf(err)

	ftypBoxes := bmfcommon.ChildBoxes(s, "ftyp")
	ftyp := ftypBoxes[0].(*bmftype.FtypBox)

	if ftyp.Name() != "ftyp" {
		t.Fatalf("ftyp name not correct: [%s]", ftyp.Name())
	}

	if ftyp.MajorBrand() != "heic" {
		t.Fatalf("ftyp MajorBrand is not correct: [%s]", ftyp.MajorBrand())
	}
}

func ExampleOpen() {
	s, err := Open(testMp4Filepath)
	log.PanicIf(err)

	ftypBoxes := bmfcommon.ChildBoxes(s, "ftyp")
	ftyp := ftypBoxes[0].(*bmftype.FtypBox)

	fmt.Printf("ftyp Name: [%s]\n", ftyp.Name())
	fmt.Printf("ftyp MajorBrand: [%s]\n", ftyp.MajorBrand())
	fmt.Printf("ftyp MinorVersion: (%d)\n", ftyp.MinorVersion())
	fmt.Printf("ftyp CompatibleBrands: %v\n", ftyp.CompatibleBrands())

	moovBoxes := bmfcommon.ChildBoxes(s, "moov")
	moov := moovBoxes[0].(*bmftype.MoovBox)

	fmt.Printf("moov Name: [%s]\n", moov.Name())
	fmt.Printf("moov Size: (%d)\n", moov.Size())

	mvhdBoxes := bmfcommon.ChildBoxes(moov, "mvhd")
	mvhd := mvhdBoxes[0].(*bmftype.MvhdBox)

	fmt.Printf("mvhd Name: [%s]\n", mvhd.Name())
	fmt.Printf("mvhd Version: (%d)\n", mvhd.Version())
	fmt.Printf("mvhd Volume: (%d)\n", mvhd.Volume())

	trakBoxes := bmfcommon.ChildBoxes(moov, "trak")
	trak0 := trakBoxes[0].(*bmftype.TrakBox)
	trak1 := trakBoxes[1].(*bmftype.TrakBox)

	fmt.Printf("trak (0) Size: (%d)\n", trak0.Size())
	fmt.Printf("trak (1) Size: (%d)\n", trak1.Size())

	mdatBoxes := bmfcommon.ChildBoxes(s, "mdat")
	mdat := mdatBoxes[0].(*bmftype.MdatBox)

	fmt.Printf("mdat Size: (%d)\n", mdat.Size())

	// Output:
	// ftyp Name: [ftyp]
	// ftyp MajorBrand: [isom]
	// ftyp MinorVersion: (512)
	// ftyp CompatibleBrands: [isom iso2 avc1 mp41]
	// moov Name: [moov]
	// moov Size: (3170)
	// mvhd Name: [mvhd]
	// mvhd Version: (0)
	// mvhd Volume: (1)
	// trak (0) Size: (1517)
	// trak (1) Size: (1439)
	// mdat Size: (2872360)
}
