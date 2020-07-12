package mp4

import (
	"fmt"
	"testing"

	"github.com/dsoprea/go-logging"
)

const (
	testMp4Filepath  = "assets/tears-of-steel.mp4"
	testHeicFilepath = "assets/image.heic"
)

func TestOpen_Mp4(t *testing.T) {
	s, err := Open(testMp4Filepath)
	log.PanicIf(err)

	if s.Ftyp().Name() != "ftyp" {
		t.Fatalf("ftyp name not correct: [%s]", s.Ftyp().Name())
	}

	if s.Ftyp().MajorBrand != "isom" {
		t.Fatalf("ftyp MajorBrand is not correct: [%s]", s.Ftyp().MajorBrand)
	}
}

func TestOpen_Heic(t *testing.T) {
	s, err := Open(testHeicFilepath)
	log.PanicIf(err)

	if s.Ftyp().Name() != "ftyp" {
		t.Fatalf("ftyp name not correct: [%s]", s.Ftyp().Name())
	}

	if s.Ftyp().MajorBrand != "heic" {
		t.Fatalf("ftyp MajorBrand is not correct: [%s]", s.Ftyp().MajorBrand)
	}
}

func ExampleOpen() {
	s, err := Open(testMp4Filepath)
	log.PanicIf(err)

	fmt.Println(s.Ftyp().Name())
	fmt.Println(s.Ftyp().MajorBrand)
	fmt.Println(s.Ftyp().MinorVersion)
	fmt.Println(s.Ftyp().CompatibleBrands)

	fmt.Println(s.Moov().Name(), s.Moov().Size())
	fmt.Println(s.Moov().Mvhd.Name())
	fmt.Println(s.Moov().Mvhd.Version)
	fmt.Println(s.Moov().Mvhd.Volume)

	fmt.Println("trak size: ", s.Moov().Traks[0].Size())
	fmt.Println("trak size: ", s.Moov().Traks[1].Size())
	fmt.Println("mdat size: ", s.Mdat().Size())

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
