package bmftype

import (
	"bytes"
	"testing"

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

	if tb.InlineString() != "NAME=[] PARENT=[ROOT] START=(0) SIZE=(0) VER=(0x11) FLAGS=(0x11223344) TRACK-ID=(85]) LAYER=(102) ALT-GROUP=(119) VOLUME=[OFF] MATRIX=(3) W=(153) H=(265) DUR-S=[60.00]" {
		t.Fatalf("InlineString() not correct: [%s]", tb.InlineString())
	}
}

func TestTkhdBoxFactory_Name(t *testing.T) {
	tbf := tkhdBoxFactory{}
	if tbf.Name() != "tkhd" {
		t.Fatalf("Name() not correct.")
	}
}
