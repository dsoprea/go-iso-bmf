package atom

import (
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"
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
	box := Box{
		name:  "abcd",
		start: 1234,
		size:  5678,
	}

	fb := FtypBox{
		Box:              box,
		majorBrand:       "efgh",
		minorVersion:     11,
		compatibleBrands: []string{"abc", "def"},
	}

	if fb.String() != "ftyp<NAME=[abcd] START=(1234) SIZE=(5678) MAJOR-BRAND=[efgh] MINOR-VER=(0x0000000b) COMPAT-BRANDS=[abc,def]>" {
		t.Fatalf("String() not correct: [%s]", fb.String())
	}
}

func TestFtypBox_InlineString(t *testing.T) {
	box := Box{
		name:  "abcd",
		start: 1234,
		size:  5678,
	}

	fb := FtypBox{
		Box:              box,
		majorBrand:       "efgh",
		minorVersion:     11,
		compatibleBrands: []string{"abc", "def"},
	}

	if fb.InlineString() != "NAME=[abcd] START=(1234) SIZE=(5678) MAJOR-BRAND=[efgh] MINOR-VER=(0x0000000b) COMPAT-BRANDS=[abc,def]" {
		t.Fatalf("InlineString() not correct: [%s]", fb.String())
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
	pushBytes(&data, uint32(11))

	// Add brands

	brands := []string{"efgh", "ijkl"}
	for _, brand := range brands {
		data = append(data, []byte(brand)...)
	}

	var b []byte
	pushBox(&b, "elst", data)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := NewFile(sb, int64(len(b)))

	box, err := file.readBaseBox(0)
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
	} else if fb.String() != "ftyp<NAME=[elst] START=(0) SIZE=(24) MAJOR-BRAND=[abcd] MINOR-VER=(0x0000000b) COMPAT-BRANDS=[efgh,ijkl]>" {
		t.Fatalf("String() not correct: [%s]", fb.String())
	}
}
