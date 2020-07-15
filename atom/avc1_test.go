package atom

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"
)

func TestAvc1Box_Version(t *testing.T) {
	ab := &Avc1Box{
		version: 123,
	}

	if ab.Version() != 123 {
		t.Fatalf("Version() not correct.")
	}
}

func TestAvc1BoxFactory_Name(t *testing.T) {
	abf := avc1BoxFactory{}

	if abf.Name() != "avc1" {
		t.Fatalf("Name() not correct.")
	}
}

func TestAvc1BoxFactory_New(t *testing.T) {
	data := []byte{
		11,
	}

	// pushBytes(&data, x interface{})

	var b []byte
	pushBox(&b, "avc1", data)

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := NewFile(sb, int64(len(b)))

	box, err := file.readBaseBox(0)
	log.PanicIf(err)

	cb, err := avc1BoxFactory{}.New(box)
	log.PanicIf(err)

	avc1 := cb.(*Avc1Box)

	if avc1.Version() != 11 {
		t.Fatalf("Version() not correct.")
	}
}
