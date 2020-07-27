package bmftype

import (
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestSttsBox_Version(t *testing.T) {
	mb := SttsBox{
		version: 11,
	}

	if mb.Version() != 11 {
		t.Fatalf("Version() not correct.")
	}
}

func TestSttsBox_Flags(t *testing.T) {
	mb := SttsBox{
		flags: 22,
	}

	if mb.Flags() != 22 {
		t.Fatalf("Flags() not correct.")
	}
}

func TestSttsBox_SampleCounts(t *testing.T) {
	d := []uint32{44, 55}

	mb := SttsBox{
		sampleCounts: d,
	}

	if reflect.DeepEqual(mb.SampleCounts(), d) != true {
		t.Fatalf("SampleCounts() not correct.")
	}
}

func TestSttsBox_SampleDeltas(t *testing.T) {
	d := []uint32{66, 77}

	mb := SttsBox{
		sampleDeltas: d,
	}

	if reflect.DeepEqual(mb.SampleDeltas(), d) != true {
		t.Fatalf("SampleDeltas() not correct.")
	}
}

func TestSttsBoxFactory_Name(t *testing.T) {
	name := sttsBoxFactory{}.Name()
	if name != "stts" {
		t.Fatalf("Name() not correct.")
	}
}

func TestSttsBoxFactory_New(t *testing.T) {
	data := []byte{}

	// flags
	bmfcommon.PushBytes(&data, uint32(0x11223344))

	// count
	bmfcommon.PushBytes(&data, uint32(3))

	// record 1 count
	bmfcommon.PushBytes(&data, uint32(11))

	// record 1 delta
	bmfcommon.PushBytes(&data, uint32(22))

	// record 2 count
	bmfcommon.PushBytes(&data, uint32(33))

	// record 2 delta
	bmfcommon.PushBytes(&data, uint32(44))

	// record 3 count
	bmfcommon.PushBytes(&data, uint32(55))

	// record 3 delta
	bmfcommon.PushBytes(&data, uint32(66))

	b := []byte{}
	bmfcommon.PushBox(&b, "stts", data)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := bmfcommon.NewBmfResource(sb, int64(len(b)))

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := sttsBoxFactory{}.New(box)
	log.PanicIf(err)

	sd := cb.(*SttsBox)

	if sd.Version() != 0x11 {
		t.Fatalf("Version() not correct: (0x%02x)", sd.Version())
	}

	if sd.Flags() != 0x11223344 {
		t.Fatalf("Flags() not correct: (0x%08x)", sd.Flags())
	}

	expectedCounts := []uint32{
		11,
		33,
		55,
	}

	if reflect.DeepEqual(sd.SampleCounts(), expectedCounts) != true {
		t.Fatalf("SampleCounts() not correct.")
	}

	expectedDeltas := []uint32{
		22,
		44,
		66,
	}

	if reflect.DeepEqual(sd.SampleDeltas(), expectedDeltas) != true {
		t.Fatalf("SampleDeltas not correct.")
	}
}
