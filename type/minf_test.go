package bmftype

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestMinfBoxFactory_Name(t *testing.T) {
	name := minfBoxFactory{}.Name()

	if name != "minf" {
		t.Fatalf("Name() not correct.")
	}
}

func TestMinfBoxFactory_New(t *testing.T) {
	b := []byte{}
	bmfcommon.PushBox(&b, "minf", nil)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := bmfcommon.NewBmfResource(sb, int64(len(b)))

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := minfBoxFactory{}.New(box)
	log.PanicIf(err)

	// Nothing else we can validate.
	_, ok := cb.(*MinfBox)

	if ok != true {
		t.Fatalf("Expected an 'minf' box.")
	}
}
