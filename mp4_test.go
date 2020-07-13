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

	ftyp := atom.MustGetChildBox(s, "ftyp").(*atom.FtypBox)

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

	ftyp := atom.MustGetChildBox(s, "ftyp").(*atom.FtypBox)

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

	ftyp := atom.MustGetChildBox(s, "ftyp").(*atom.FtypBox)

	fmt.Println(ftyp.Name())
	fmt.Println(ftyp.MajorBrand)
	fmt.Println(ftyp.MinorVersion)
	fmt.Println(ftyp.CompatibleBrands)

	moov := atom.MustGetChildBox(s, "moov").(*atom.MoovBox)

	fmt.Println(moov.Name(), moov.Size())

	mvhd := atom.MustGetChildBox(moov, "mvhd").(*atom.MvhdBox)

	fmt.Println(mvhd.Name())
	fmt.Println(mvhd.Version)
	fmt.Println(mvhd.Volume)

	// TODO(dustin): !! Finish this. A sequence of box types may include the same box-type multiple times. We need to update the index to have slices instead of single values. Then, we can update this to access that slice.
	// fmt.Println("trak size: ", moov.Traks[0].Size())
	// fmt.Println("trak size: ", moov.Traks[1].Size())

	mdat := atom.MustGetChildBox(s, "mdat").(*atom.MdatBox)

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
