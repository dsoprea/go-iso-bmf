package mp4

import (
	"fmt"
	"testing"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-mp4/atom"
)

const (
	testMp4Filepath  = "assets/tears-of-steel.mp4"
	testHeicFilepath = "assets/image.heic"
)

func TestOpen_Mp4(t *testing.T) {
	s, err := Open(testMp4Filepath)
	log.PanicIf(err)

	ftypBoxes := atom.ChildBoxes(s, "ftyp")
	ftyp := ftypBoxes[0].(*atom.FtypBox)

	if ftyp.Name() != "ftyp" {
		t.Fatalf("ftyp name not correct: [%s]", ftyp.Name())
	}

	if ftyp.MajorBrand != "isom" {
		t.Fatalf("ftyp MajorBrand is not correct: [%s]", ftyp.MajorBrand)
	}
}

func TestOpen_Heic(t *testing.T) {
	s, err := Open(testHeicFilepath)
	log.PanicIf(err)

	ftypBoxes := atom.ChildBoxes(s, "ftyp")
	ftyp := ftypBoxes[0].(*atom.FtypBox)

	if ftyp.Name() != "ftyp" {
		t.Fatalf("ftyp name not correct: [%s]", ftyp.Name())
	}

	if ftyp.MajorBrand != "heic" {
		t.Fatalf("ftyp MajorBrand is not correct: [%s]", ftyp.MajorBrand)
	}
}

func ExampleOpen() {
	s, err := Open(testMp4Filepath)
	log.PanicIf(err)

	ftypBoxes := atom.ChildBoxes(s, "ftyp")
	ftyp := ftypBoxes[0].(*atom.FtypBox)

	fmt.Println(ftyp.Name())
	fmt.Println(ftyp.MajorBrand)
	fmt.Println(ftyp.MinorVersion)
	fmt.Println(ftyp.CompatibleBrands)

	moovBoxes := atom.ChildBoxes(s, "moov")
	moov := moovBoxes[0].(*atom.MoovBox)

	fmt.Println(moov.Name(), moov.Size())

	mvhdBoxes := atom.ChildBoxes(moov, "mvhd")
	mvhd := mvhdBoxes[0].(*atom.MvhdBox)

	fmt.Println(mvhd.Name())
	fmt.Println(mvhd.Version)
	fmt.Println(mvhd.Volume)

	trakBoxes := atom.ChildBoxes(moov, "trak")
	trak0 := trakBoxes[0].(*atom.TrakBox)
	trak1 := trakBoxes[1].(*atom.TrakBox)

	fmt.Println("trak size: ", trak0.Size())
	fmt.Println("trak size: ", trak1.Size())

	mdatBoxes := atom.ChildBoxes(s, "mdat")
	mdat := mdatBoxes[0].(*atom.MdatBox)

	fmt.Println("mdat size: ", mdat.Size())

	// Output:
	// ftyp
	// isom
	// 512
	// [isom iso2 avc1 mp41]
	// moov 3170
	// mvhd
	// 0
	// 1
	// trak size:  1517
	// trak size:  1439
	// mdat size:  2872360
}
