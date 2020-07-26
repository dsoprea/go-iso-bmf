package bmftype

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestHmhdBox_Version(t *testing.T) {
	hb := HmhdBox{
		version: 11,
	}

	if hb.Version() != 11 {
		t.Fatalf("Version() not correct.")
	}
}

func TestHmhdBox_MaxPDUSize(t *testing.T) {
	hb := HmhdBox{
		maxPDUSize: 11,
	}

	if hb.MaxPDUSize() != 11 {
		t.Fatalf("MaxPDUSize() not correct.")
	}
}

func TestHmhdBox_AvgPDUSize(t *testing.T) {
	hb := HmhdBox{
		avgPDUSize: 11,
	}

	if hb.AvgPDUSize() != 11 {
		t.Fatalf("AvgPDUSize() not correct.")
	}
}

func TestHmhdBox_MaxBitrate(t *testing.T) {
	hb := HmhdBox{
		maxBitrate: 11,
	}

	if hb.MaxBitrate() != 11 {
		t.Fatalf("MaxBitrate() not correct.")
	}
}

func TestHmhdBox_AvgBitrate(t *testing.T) {
	hb := HmhdBox{
		avgBitrate: 11,
	}

	if hb.AvgBitrate() != 11 {
		t.Fatalf("AvgBitrate() not correct.")
	}
}

func TestHmhdBoxFactory_Name(t *testing.T) {
	name := hmhdBoxFactory{}.Name()

	if name != "hmhd" {
		t.Fatalf("Name() not correct.")
	}
}

func TestHmhdBoxFactory_New(t *testing.T) {
	data := []byte{
		0x11,
	}

	bmfcommon.PushBytes(&data, uint16(0x22))
	bmfcommon.PushBytes(&data, uint16(0x33))
	bmfcommon.PushBytes(&data, uint32(0x44))
	bmfcommon.PushBytes(&data, uint32(0x55))

	b := []byte{}
	bmfcommon.PushBox(&b, "hmnd", data)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := bmfcommon.NewFile(sb, int64(len(b)))

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := hmhdBoxFactory{}.New(box)
	log.PanicIf(err)

	hb := cb.(*HmhdBox)

	if hb.Version() != 0x11 {
		t.Fatalf("Version() not correct.")
	} else if hb.MaxPDUSize() != 0x22 {
		t.Fatalf("MaxPDUSize() not correct: (0x%04x)", hb.MaxPDUSize())
	} else if hb.AvgPDUSize() != 0x33 {
		t.Fatalf("AvgPDUSize() not correct: (0x%04x)", hb.AvgPDUSize())
	} else if hb.MaxBitrate() != 0x44 {
		t.Fatalf("MaxBitrate() not correct: (0x%08x)", hb.MaxBitrate())
	} else if hb.AvgBitrate() != 0x55 {
		t.Fatalf("AvgBitrate() not correct: (0x%08x)", hb.AvgBitrate())
	}
}
