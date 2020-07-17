package mp4box

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"
)

func TestTrakBoxFactory_Name(t *testing.T) {
	name := trakBoxFactory{}.Name()

	if name != "trak" {
		t.Fatalf("Name() not correct.")
	}
}

func TestTrakBoxFactory_New(t *testing.T) {
	b := []byte{}
	pushBox(&b, "trak", nil)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := NewFile(sb, int64(len(b)))

	box, err := file.readBaseBox(0)
	log.PanicIf(err)

	cb, err := trakBoxFactory{}.New(box)
	log.PanicIf(err)

	// Nothing else we can validate.
	_, ok := cb.(*TrakBox)

	if ok != true {
		t.Fatalf("Expected an 'trak' box.")
	}
}
