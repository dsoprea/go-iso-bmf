package bmftype

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestTrakBoxFactory_Name(t *testing.T) {
	name := trakBoxFactory{}.Name()

	if name != "trak" {
		t.Fatalf("Name() not correct.")
	}
}

func TestTrakBoxFactory_New(t *testing.T) {
	b := []byte{}
	bmfcommon.PushBox(&b, "trak", nil)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file, err := bmfcommon.NewBmfResource(sb, int64(len(b)))
	log.PanicIf(err)

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := trakBoxFactory{}.New(box)
	log.PanicIf(err)

	// Nothing else we can validate.
	_, ok := cb.(*TrakBox)

	if ok != true {
		t.Fatalf("Expected an 'trak' box.")
	}
}
