package bmftype

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestAvc1Box_Version(t *testing.T) {
	ab := &Avc1Box{
		version: 123,
	}

	if ab.Version() != 123 {
		t.Fatalf("Version() not correct.")
	}
}

func TestAvc1BoxFactory_Name(t *testing.T) {
	abf := avc1BoxFactory{}

	if abf.Name() != "avc1" {
		t.Fatalf("Name() not correct.")
	}
}

func TestAvc1BoxFactory_New(t *testing.T) {
	data := []byte{
		11,
	}

	// bmfcommon.PushBytes(&data, x interface{})

	var b []byte
	bmfcommon.PushBox(&b, "avc1", data)

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := bmfcommon.NewFile(sb, int64(len(b)))

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, err := avc1BoxFactory{}.New(box)
	log.PanicIf(err)

	avc1 := cb.(*Avc1Box)

	if avc1.Version() != 11 {
		t.Fatalf("Version() not correct.")
	}
}
