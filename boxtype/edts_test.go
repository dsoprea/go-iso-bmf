package boxtype

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"
)

func TestEdtsBoxFactory_Name(t *testing.T) {
	ebf := edtsBoxFactory{}

	if ebf.Name() != "edts" {
		t.Fatalf("Name() not correct.")
	}
}

func TestEdtsBoxFactory_New(t *testing.T) {
	var b []byte
	pushBox(&b, "edts", nil)

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := NewFile(sb, int64(len(b)))

	box, err := file.readBaseBox(0)
	log.PanicIf(err)

	cb, err := edtsBoxFactory{}.New(box)
	log.PanicIf(err)

	// Nothing else we can validate.
	_, ok := cb.(*EdtsBox)

	if ok != true {
		t.Fatalf("Expected an 'edts' box.")
	}
}
