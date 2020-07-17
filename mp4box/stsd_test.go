package mp4box

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"
)

func TestStsdBox_Version(t *testing.T) {
	sb := StsdBox{
		version: 0x11,
	}

	if sb.Version() != 0x11 {
		t.Fatalf("Version() not correct.")
	}
}

func TestStsdBox_Flags(t *testing.T) {
	sb := StsdBox{
		flags: 0x22,
	}

	if sb.Flags() != 0x22 {
		t.Fatalf("Flags() not correct.")
	}
}

func TestStsdBoxFactory_Name(t *testing.T) {
	name := stsdBoxFactory{}.Name()

	if name != "stsd" {
		t.Fatalf("Name() not correct.")
	}
}

func TestStsdBoxFactory_New(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			log.PrintError(err)

			t.Fatalf("Test failed.")
		}
	}()

	var data []byte

	// flags
	pushBytes(&data, uint32(0x11223344))

	// The child boxes are read at offset (8) in the content. Since this test
	// doesn't add any actual boxes, the content will be exactly eight bytes.
	data = append(data, 0, 0, 0, 0)

	var b []byte
	pushBox(&b, "stsd", data)

	DumpBytes(b)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := NewFile(sb, int64(len(b)))

	box, err := file.readBaseBox(0)
	log.PanicIf(err)

	cb, err := stsdBoxFactory{}.New(box)
	log.PanicIf(err)

	mb := cb.(*StsdBox)

	if mb.Version() != 0x11 {
		t.Fatalf("Version() not correct: (0x%02x)", mb.Version())
	}

	if mb.Flags() != 0x11223344 {
		t.Fatalf("Flags() not correct: (0x%08x)", mb.Flags())
	}
}
