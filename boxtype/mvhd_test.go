package boxtype

import (
	"testing"
	"time"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"
)

func TestMvhdBox_Flags(t *testing.T) {
	mb := MvhdBox{
		flags: 11,
	}

	if mb.Flags() != 11 {
		t.Fatalf("Flags() is incorrect.")
	}
}

func TestMvhdBox_Version(t *testing.T) {
	mb := MvhdBox{
		version: 22,
	}

	if mb.Version() != 22 {
		t.Fatalf("Version() is incorrect.")
	}
}

func TestMvhdBox_Rate(t *testing.T) {
	mb := MvhdBox{
		rate: 33,
	}

	if mb.Rate() != 33 {
		t.Fatalf("Rate() is incorrect.")
	}
}

func TestMvhdBox_Volume(t *testing.T) {
	mb := MvhdBox{
		volume: 44,
	}

	if mb.Volume() != 44 {
		t.Fatalf("Volume() is incorrect.")
	}
}

func TestMvhdBox_InlineString(t *testing.T) {
	timeScale := uint32(1)
	duration := uint32(60)

	sts := newStandard32TimeSupport(
		0,
		0,
		duration,
		timeScale)

	mb := MvhdBox{
		standard32TimeSupport: sts,
		flags:                 11,
		version:               22,
		rate:                  33,
		volume:                44,
	}

	if mb.InlineString() != "NAME=[] START=(0) SIZE=(0) VER=(0x16) FLAGS=(0x0000000b) RATE=(33]) VOLUME=(44) DUR-S=[60.00]" {
		t.Fatalf("InlineString() is incorrect: [%s]", mb.InlineString())
	}
}

func TestMvhdBoxFactory_Name(t *testing.T) {
	mbf := mvhdBoxFactory{}
	if mbf.Name() != "mvhd" {
		t.Fatalf("Name() not correct.")
	}
}

func TestMvhdBoxFactory_New(t *testing.T) {
	var data []byte

	// flags
	pushBytes(&data, uint32(0x00000011))

	// creation and modified epochs

	epoch := uint32(3677725917)
	baseTime := EpochToTime(epoch)

	// creation epoch
	pushBytes(&data, epoch)

	// modification epoch
	pushBytes(&data, epoch+1)

	// timeScale
	pushBytes(&data, uint32(30))

	// scaledDuration
	pushBytes(&data, uint32(300))

	// rate
	pushBytes(&data, uint32(22))

	// volume
	pushBytes(&data, uint16(33))

	b := []byte{}
	pushBox(&b, "mvhd", data)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := NewFile(sb, int64(len(b)))

	box, err := file.readBaseBox(0)
	log.PanicIf(err)

	cb, err := mvhdBoxFactory{}.New(box)
	log.PanicIf(err)

	mb := cb.(*MvhdBox)

	if mb.Version() != 0 {
		t.Fatalf("Version() not correct: (0x%02x)", mb.Version())
	}

	if mb.CreationTime() != baseTime {
		t.Fatalf("CreationTime() not correct: [%s] != [%s]", mb.CreationTime(), baseTime)
	}

	if mb.ModificationTime() != baseTime.Add(1*time.Second) {
		t.Fatalf("ModificationTime() not correct: %s", mb.ModificationTime())
	}

	if mb.TimeScale() != 30 {
		t.Fatalf("TimeScale() not correct.")
	}

	if mb.ScaledDuration() != 300 {
		t.Fatalf("ScaledDuration() not correct.")
	}

	if mb.Rate() != 22 {
		t.Fatalf("Rate() not correct: (%d)", mb.Rate())
	}

	if mb.Volume() != 33 {
		t.Fatalf("Volume() not correct.")
	}
}