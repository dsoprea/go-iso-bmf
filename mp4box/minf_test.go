package mp4box

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"
)

func TestMinfBoxFactory_Name(t *testing.T) {
	name := minfBoxFactory{}.Name()

	if name != "minf" {
		t.Fatalf("Name() not correct.")
	}
}

func TestMinfBoxFactory_New(t *testing.T) {
	b := []byte{}
	pushBox(&b, "minf", nil)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := NewFile(sb, int64(len(b)))

	box, err := file.readBaseBox(0)
	log.PanicIf(err)

	cb, err := minfBoxFactory{}.New(box)
	log.PanicIf(err)

	// Nothing else we can validate.
	_, ok := cb.(*MinfBox)

	if ok != true {
		t.Fatalf("Expected an 'minf' box.")
	}
}
