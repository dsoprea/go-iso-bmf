package boxtype

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"
)

func TestMdiaBoxFactory_Name(t *testing.T) {
	name := mdiaBoxFactory{}.Name()

	if name != "mdia" {
		t.Fatalf("Name() not correct.")
	}
}

func TestMdiaBoxFactory_New(t *testing.T) {
	b := []byte{}
	pushBox(&b, "mdia", nil)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := NewFile(sb, int64(len(b)))

	box, err := file.readBaseBox(0)
	log.PanicIf(err)

	cb, err := mdiaBoxFactory{}.New(box)
	log.PanicIf(err)

	// Nothing else we can validate.
	_, ok := cb.(*MdiaBox)

	if ok != true {
		t.Fatalf("Expected an 'mdia' box.")
	}
}
