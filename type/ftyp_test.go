package bmftype

import (
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestFtypBox_MajorBrand(t *testing.T) {
	fb := FtypBox{
		majorBrand: "abcdef",
	}

	if fb.MajorBrand() != "abcdef" {
		t.Fatalf("MajorBrand() not correct.")
	}
}

func TestFtypBox_MinorVersion(t *testing.T) {
	fb := FtypBox{
		minorVersion: 11,
	}

	if fb.MinorVersion() != 11 {
		t.Fatalf("MinorVersion() not correct.")
	}
}

func TestFtypBox_CompatibleBrands(t *testing.T) {
	brands := []string{"abc", "def"}

	fb := FtypBox{
		compatibleBrands: brands,
	}

	if reflect.DeepEqual(fb.CompatibleBrands(), brands) != true {
		t.Fatalf("CompatibleBrands() not correct.")
	}
}

func TestFtypBox_String(t *testing.T) {
	box := bmfcommon.NewBox("abcd", 1234, 5678, nil)

	fb := FtypBox{
		Box:              box,
		majorBrand:       "efgh",
		minorVersion:     11,
		compatibleBrands: []string{"abc", "def"},
	}

	if fb.String() != "ftyp<NAME=[abcd] PARENT=[ROOT] START=(0x00000000000004d2) SIZE=(5678) MAJOR-BRAND=[efgh] MINOR-VER=(0x0000000b) COMPAT-BRANDS=[abc,def]>" {
		t.Fatalf("String() not correct: [%s]", fb.String())
	}
}

func TestFtypBox_InlineString(t *testing.T) {
	box := bmfcommon.NewBox("abcd", 1234, 5678, nil)

	fb := FtypBox{
		Box:              box,
		majorBrand:       "efgh",
		minorVersion:     11,
		compatibleBrands: []string{"abc", "def"},
	}

	if fb.InlineString() != "NAME=[abcd] PARENT=[ROOT] START=(0x00000000000004d2) SIZE=(5678) MAJOR-BRAND=[efgh] MINOR-VER=(0x0000000b) COMPAT-BRANDS=[abc,def]" {
		t.Fatalf("InlineString() not correct: [%s]", fb.InlineString())
	}
}

func TestFtypBoxFactory_Name(t *testing.T) {
	name := ftypBoxFactory{}.Name()

	if name != "ftyp" {
		t.Fatalf("Name() not correct.")
	}
}

func TestFtypBoxFactory_New(t *testing.T) {
	// Load

	var data []byte

	// majorBrand
	data = append(data, 'a', 'b', 'c', 'd')

	// minorVersion
	bmfcommon.PushBytes(&data, uint32(11))

	// Add brands

	brands := []string{"efgh", "ijkl"}
	for _, brand := range brands {
		data = append(data, []byte(brand)...)
	}

	var b []byte
	bmfcommon.PushBox(&b, "elst", data)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := bmfcommon.NewFile(sb, int64(len(b)))

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, err := ftypBoxFactory{}.New(box)
	log.PanicIf(err)

	fb := cb.(*FtypBox)

	if fb.MajorBrand() != "abcd" {
		t.Fatalf("MajorBrand() not correct.")
	} else if fb.MinorVersion() != 11 {
		t.Fatalf("MinorVersion() not correct.")
	} else if reflect.DeepEqual(fb.CompatibleBrands(), brands) != true {
		t.Fatalf("CompatibleBrands() not correct.")
	} else if fb.String() != "ftyp<NAME=[elst] PARENT=[ROOT] START=(0x0000000000000000) SIZE=(24) MAJOR-BRAND=[abcd] MINOR-VER=(0x0000000b) COMPAT-BRANDS=[efgh,ijkl]>" {
		t.Fatalf("String() not correct: [%s]", fb.String())
	}
}
