package mp4

import (
	"fmt"
	"testing"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-mp4/mp4box"
)

const (
	testMp4Filepath  = "assets/tears-of-steel.mp4"
	testHeicFilepath = "assets/image.heic"
)

func TestOpen_Mp4(t *testing.T) {
	s, err := Open(testMp4Filepath)
	log.PanicIf(err)

	ftypBoxes := mp4box.ChildBoxes(s, "ftyp")
	ftyp := ftypBoxes[0].(*mp4box.FtypBox)

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

	ftypBoxes := mp4box.ChildBoxes(s, "ftyp")
	ftyp := ftypBoxes[0].(*mp4box.FtypBox)

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

	ftypBoxes := mp4box.ChildBoxes(s, "ftyp")
	ftyp := ftypBoxes[0].(*mp4box.FtypBox)

	fmt.Println(ftyp.Name())
	fmt.Println(ftyp.MajorBrand())
	fmt.Println(ftyp.MinorVersion())
	fmt.Println(ftyp.CompatibleBrands())

	moovBoxes := mp4box.ChildBoxes(s, "moov")
	moov := moovBoxes[0].(*mp4box.MoovBox)

	fmt.Println(moov.Name(), moov.Size())

	mvhdBoxes := mp4box.ChildBoxes(moov, "mvhd")
	mvhd := mvhdBoxes[0].(*mp4box.MvhdBox)

	fmt.Println(mvhd.Name())
	fmt.Println(mvhd.Version())
	fmt.Println(mvhd.Volume())

	trakBoxes := mp4box.ChildBoxes(moov, "trak")
	trak0 := trakBoxes[0].(*mp4box.TrakBox)
	trak1 := trakBoxes[1].(*mp4box.TrakBox)

	fmt.Println("trak size: ", trak0.Size())
	fmt.Println("trak size: ", trak1.Size())

	mdatBoxes := mp4box.ChildBoxes(s, "mdat")
	mdat := mdatBoxes[0].(*mp4box.MdatBox)

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
