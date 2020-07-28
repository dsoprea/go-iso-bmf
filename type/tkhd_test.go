package bmftype

import (
	"bytes"
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestTkhdBox_Version(t *testing.T) {
	tb := TkhdBox{
		version: 0x11,
	}

	if tb.Version() != 0x11 {
		t.Fatalf("Version() not correct.")
	}
}

func TestTkhdBox_Flags(t *testing.T) {
	tb := TkhdBox{
		flags: 0x11223344,
	}

	if tb.Flags() != 0x11223344 {
		t.Fatalf("Flags() not correct.")
	}
}

func TestTkhdBox_TrackId(t *testing.T) {
	tb := TkhdBox{
		trackId: 0x55,
	}

	if tb.TrackId() != 0x55 {
		t.Fatalf("TrackId() not correct.")
	}
}

func TestTkhdBox_Layer(t *testing.T) {
	tb := TkhdBox{
		layer: 0x66,
	}

	if tb.Layer() != 0x66 {
		t.Fatalf("Layer() not correct.")
	}
}

func TestTkhdBox_AlternateGroup(t *testing.T) {
	tb := TkhdBox{
		alternateGroup: 0x77,
	}

	if tb.AlternateGroup() != 0x77 {
		t.Fatalf("AlternateGroup() not correct.")
	}
}

func TestTkhdBox_Volume(t *testing.T) {
	tb := TkhdBox{
		volume: 0x88,
	}

	if tb.Volume() != 0x88 {
		t.Fatalf("Volume() not correct.")
	}
}

func TestTkhdBox_Matrix(t *testing.T) {
	data := []byte{1, 2, 3}

	tb := TkhdBox{
		matrix: data,
	}

	if bytes.Equal(tb.Matrix(), data) != true {
		t.Fatalf("Matrix() not correct.")
	}
}

func TestTkhdBox_Width(t *testing.T) {
	tb := TkhdBox{
		width: 0x99,
	}

	if tb.Width() != 0x99 {
		t.Fatalf("Width() not correct.")
	}
}

func TestTkhdBox_Height(t *testing.T) {
	tb := TkhdBox{
		height: 0x109,
	}

	if tb.Height() != 0x109 {
		t.Fatalf("Height() not correct.")
	}
}

func TestTkhdBox_InlineString(t *testing.T) {
	timeScale := uint32(1)
	duration := uint32(60)

	sts := bmfcommon.NewStandard32TimeSupport(
		0,
		0,
		duration,
		timeScale)

	tb := TkhdBox{
		Standard32TimeSupport: sts,

		version:        0x11,
		flags:          0x11223344,
		trackId:        0x55,
		layer:          0x66,
		alternateGroup: 0x77,
		volume:         0x88,
		matrix:         []byte{1, 2, 3},
		width:          0x99,
		height:         0x109,
	}

	if tb.InlineString() != "NAME=[] PARENT=[ROOT] START=(0x0000000000000000) SIZE=(0) VER=(0x11) FLAGS=(0x11223344) TRACK-ID=(85) LAYER=(102) ALT-GROUP=(119) VOLUME=[OFF] MATRIX=(3) W=(153) H=(265) DUR-S=[60.00]" {
		t.Fatalf("InlineString() not correct: [%s]", tb.InlineString())
	}
}

func TestTkhdBoxFactory_Name(t *testing.T) {
	tbf := tkhdBoxFactory{}
	if tbf.Name() != "tkhd" {
		t.Fatalf("Name() not correct.")
	}
}

func TestTkhdBoxFactory_New(t *testing.T) {
	// Construct the stream of TKHD data:

	var data []byte

	// flags
	bmfcommon.PushBytes(&data, uint32(0x00223344))

	// creationEpoch, modificationEpoch

	now := bmfcommon.NowTime()

	creationEpoch := bmfcommon.TimeToEpoch(now)
	bmfcommon.PushBytes(&data, creationEpoch)

	modificationEpoch := creationEpoch + 1
	bmfcommon.PushBytes(&data, modificationEpoch)

	// trackId
	bmfcommon.PushBytes(&data, uint32(0x11))

	// (reserved)
	bmfcommon.PushBytes(&data, uint32(0))

	// duration
	bmfcommon.PushBytes(&data, uint32(300))

	// (reserved)
	bmfcommon.PushBytes(&data, uint32(0))
	bmfcommon.PushBytes(&data, uint32(0))

	// layer
	bmfcommon.PushBytes(&data, uint16(0x22))

	// alternateGroup
	bmfcommon.PushBytes(&data, uint16(0x33))

	// volume: 0000 0100 0000 1000 -> 4/8
	bmfcommon.PushBytes(&data, uint16(0b0000010000001000))

	// (reserved)
	bmfcommon.PushBytes(&data, uint16(0))

	// matrix
	matrixData := make([]byte, 36)
	data = append(data, matrixData...)

	// width
	bmfcommon.PushBytes(&data, uint32(0x00110000))

	// height
	bmfcommon.PushBytes(&data, uint32(0x00220000))

	b := []byte{}
	bmfcommon.PushBox(&b, "tkhd", data)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := bmfcommon.NewBmfResource(sb, 0)

	// Register an MVHD so the TKHD factory can find it.

	timeScale := uint32(60)
	duration := uint32(60)

	sts := bmfcommon.NewStandard32TimeSupport(
		0,
		0,
		duration,
		timeScale)

	mvhd := &MvhdBox{
		Standard32TimeSupport: sts,
	}

	fbi := file.Index()
	fbi.Add(mvhd)

	mvhdIbe := bmfcommon.IndexedBoxEntry{"moov.mvhd", 0}
	fbi[mvhdIbe] = mvhd

	// Now, try to parse the TKHD.

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := tkhdBoxFactory{}.New(box)
	log.PanicIf(err)

	tb := cb.(*TkhdBox)

	if tb.Version() != 0 {
		t.Fatalf("Version() not correct: (0x%02x)", tb.Version())
	}

	if tb.Flags() != 0x00223344 {
		t.Fatalf("Flags() not correct.")
	}

	if tb.TrackId() != 0x11 {
		t.Fatalf("TrackId() not correct.")
	}

	if tb.Layer() != 0x22 {
		t.Fatalf("Layer() not correct.")
	}

	if tb.AlternateGroup() != 0x33 {
		t.Fatalf("AlternateGroup() not correct.")
	}

	// 0000 0100 0000 1000
	if tb.Volume() != 0b0000010000001000 {
		t.Fatalf("Volume() not correct.")
	}

	numerator, denominator := tb.Volume().Decode().Rational()

	if numerator != 4 || denominator != 8 {
		t.Fatalf("Volume rational not correct: (%d)/(%d)", numerator, denominator)
	}

	if len(tb.Matrix()) != 36 {
		t.Fatalf("Matrix() not correct.")
	}

	if tb.Width() != 0x11 {
		t.Fatalf("Width() not correct.")
	}

	if tb.Height() != 0x22 {
		t.Fatalf("Height() not correct.")
	}
}
