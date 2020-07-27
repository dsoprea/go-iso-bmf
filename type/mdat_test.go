package bmftype

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestMdatBoxFactory_Name(t *testing.T) {
	name := mdatBoxFactory{}.Name()

	if name != "mdat" {
		t.Fatalf("Name() not correct.")
	}
}

func TestMdatBoxFactory_New(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			log.PrintError(err)

			t.Fatalf("Test failed.")
		}
	}()

	var b []byte
	bmfcommon.PushBox(&b, "mdat", nil)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := bmfcommon.NewBmfResource(sb, int64(len(b)))

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := mdatBoxFactory{}.New(box)
	log.PanicIf(err)

	// Nothing else we can validate.
	_, ok := cb.(*MdatBox)

	if ok != true {
		t.Fatalf("Expected an 'mdat' box.")
	}
}
