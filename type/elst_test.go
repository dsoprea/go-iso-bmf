package bmftype

import (
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
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
	bmfcommon.PushBytes(&data, uint32(11))

	// Entry count.
	bmfcommon.PushBytes(&data, uint32(2))

	// Push entry (1).
	bmfcommon.PushBytes(&data, uint32(11))
	bmfcommon.PushBytes(&data, uint32(22))
	bmfcommon.PushBytes(&data, uint16(33))
	bmfcommon.PushBytes(&data, uint16(44))

	// Push entry (2).
	bmfcommon.PushBytes(&data, uint32(55))
	bmfcommon.PushBytes(&data, uint32(66))
	bmfcommon.PushBytes(&data, uint16(77))
	bmfcommon.PushBytes(&data, uint16(88))

	var b []byte
	bmfcommon.PushBox(&b, "elst", data)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := bmfcommon.NewBmfResource(sb, int64(len(b)))

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := elstBoxFactory{}.New(box)
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
