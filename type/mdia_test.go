package bmftype

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestMdiaBoxFactory_Name(t *testing.T) {
	name := mdiaBoxFactory{}.Name()

	if name != "mdia" {
		t.Fatalf("Name() not correct.")
	}
}

func TestMdiaBoxFactory_New(t *testing.T) {
	b := []byte{}
	bmfcommon.PushBox(&b, "mdia", nil)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := bmfcommon.NewFile(sb, int64(len(b)))

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, err := mdiaBoxFactory{}.New(box)
	log.PanicIf(err)

	// Nothing else we can validate.
	_, ok := cb.(*MdiaBox)

	if ok != true {
		t.Fatalf("Expected an 'mdia' box.")
	}
}
