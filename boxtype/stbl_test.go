package boxtype

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"
)

func TestStblBoxFactory_Name(t *testing.T) {
	name := stblBoxFactory{}.Name()

	if name != "stbl" {
		t.Fatalf("Name() not correct.")
	}
}

func TestStblBoxFactory_New(t *testing.T) {
	b := []byte{}
	pushBox(&b, "stbl", nil)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := NewFile(sb, int64(len(b)))

	box, err := file.readBaseBox(0)
	log.PanicIf(err)

	cb, err := stblBoxFactory{}.New(box)
	log.PanicIf(err)

	// Nothing else we can validate.
	_, ok := cb.(*StblBox)

	if ok != true {
		t.Fatalf("Expected an 'stbl' box.")
	}
}
