package boxtype

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"
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
	pushBox(&b, "mdat", nil)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := NewFile(sb, int64(len(b)))

	box, err := file.readBaseBox(0)
	log.PanicIf(err)

	cb, err := mdatBoxFactory{}.New(box)
	log.PanicIf(err)

	// Nothing else we can validate.
	_, ok := cb.(*MdatBox)

	if ok != true {
		t.Fatalf("Expected an 'mdat' box.")
	}
}
