package boxtype

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"
)

func TestMoovBoxFactory_Name(t *testing.T) {
	name := moovBoxFactory{}.Name()

	if name != "moov" {
		t.Fatalf("Name() not correct.")
	}
}

func TestMoovBoxFactory_New(t *testing.T) {
	b := []byte{}
	pushBox(&b, "moov", nil)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := NewFile(sb, int64(len(b)))

	box, err := file.readBaseBox(0)
	log.PanicIf(err)

	cb, err := moovBoxFactory{}.New(box)
	log.PanicIf(err)

	// Nothing else we can validate.
	_, ok := cb.(*MoovBox)

	if ok != true {
		t.Fatalf("Expected an 'moov' box.")
	}
}
