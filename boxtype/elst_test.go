package boxtype

import (
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"
)

func TestElstBox_Version(t *testing.T) {
	eb := &ElstBox{
		version: 11,
	}

	if eb.Version() != 11 {
		t.Fatalf("Version() not correct.")
	}
}

func TestElstBox_Entries(t *testing.T) {
	entries := []elstEntry{
		elstEntry{segmentDuration: 11, mediaTime: 22, mediaRate: 33, mediaRateFraction: 44},
		elstEntry{segmentDuration: 55, mediaTime: 66, mediaRate: 77, mediaRateFraction: 88},
	}

	eb := &ElstBox{
		entries: entries,
	}

	if reflect.DeepEqual(eb.Entries(), entries) != true {
		t.Fatalf("Entries() not correct.")
	}
}

func TestElstBoxFactory_Name(t *testing.T) {
	name := elstBoxFactory{}.Name()

	if name != "elst" {
		t.Fatalf("Name() not correct.")
	}
}

func TestElstBoxFactory_New(t *testing.T) {
	// Load

	var data []byte

	// Version.
	pushBytes(&data, uint32(11))

	// Entry count.
	pushBytes(&data, uint32(2))

	// Push entry (1).
	pushBytes(&data, uint32(11))
	pushBytes(&data, uint32(22))
	pushBytes(&data, uint16(33))
	pushBytes(&data, uint16(44))

	// Push entry (2).
	pushBytes(&data, uint32(55))
	pushBytes(&data, uint32(66))
	pushBytes(&data, uint16(77))
	pushBytes(&data, uint16(88))

	var b []byte
	pushBox(&b, "elst", data)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := NewFile(sb, int64(len(b)))

	box, err := file.readBaseBox(0)
	log.PanicIf(err)

	cb, err := elstBoxFactory{}.New(box)
	log.PanicIf(err)

	elst := cb.(*ElstBox)

	if elst.Version() != 11 {
		t.Fatalf("Version() not correct.")
	}

	entries := []elstEntry{
		elstEntry{segmentDuration: 11, mediaTime: 22, mediaRate: 33, mediaRateFraction: 44},
		elstEntry{segmentDuration: 55, mediaTime: 66, mediaRate: 77, mediaRateFraction: 88},
	}

	if reflect.DeepEqual(elst.Entries(), entries) != true {
		t.Fatalf("Entries() not correct.")
	}
}
