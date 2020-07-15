package mp4box

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"
)

func TestHdlrBox_Version(t *testing.T) {
	hb := HdlrBox{
		version: 11,
	}

	if hb.Version() != 11 {
		t.Fatalf("Version() not correct.")
	}
}

func TestHdlrBox_Flags(t *testing.T) {
	hb := HdlrBox{
		flags: 11,
	}

	if hb.Flags() != 11 {
		t.Fatalf("Flags() not correct.")
	}
}

func TestHdlrBox_Handler(t *testing.T) {
	hb := HdlrBox{
		handler: "handler_test",
	}

	if hb.Handler() != "handler_test" {
		t.Fatalf("Handler() not correct.")
	}
}

func TestHdlrBox_HdlrName(t *testing.T) {
	hb := HdlrBox{
		hdlrName: "name_test",
	}

	if hb.HdlrName() != "name_test" {
		t.Fatalf("HdlrName() not correct.")
	}
}

func TestHdlrBox_String(t *testing.T) {
	box := Box{
		name:  "abcd",
		start: 1234,
		size:  5678,
	}

	hb := HdlrBox{
		Box:      box,
		version:  11,
		flags:    11,
		handler:  "handler_test",
		hdlrName: "name_test",
	}

	if hb.String() != "hdlr<NAME=[abcd] START=(1234) SIZE=(5678) VER=(0x0b) FLAGS=(0x0000000b) HANDLER=[handler_test] HDLR-NAME=[name_test]>" {
		t.Fatalf("String() not correct: [%s]", hb.String())
	}
}

func TestHdlrBox_InlineString(t *testing.T) {
	box := Box{
		name:  "abcd",
		start: 1234,
		size:  5678,
	}

	hb := HdlrBox{
		Box:      box,
		version:  11,
		flags:    11,
		handler:  "handler_test",
		hdlrName: "name_test",
	}

	if hb.InlineString() != "NAME=[abcd] START=(1234) SIZE=(5678) VER=(0x0b) FLAGS=(0x0000000b) HANDLER=[handler_test] HDLR-NAME=[name_test]" {
		t.Fatalf("InlineString() not correct: [%s]", hb.InlineString())
	}
}

func TestHdlrBoxFactory_Name(t *testing.T) {
	name := hdlrBoxFactory{}.Name()

	if name != "hdlr" {
		t.Fatalf("Name() not correct.")
	}
}

func TestHdlrBoxFactory_New(t *testing.T) {
	// Load

	var flagsBytes []byte
	pushBytes(&flagsBytes, uint32(0x11223344))

	hdlrBoxData := []byte{
		// 0: version (1)
		0x0b,

		// 1: flags (3)
		flagsBytes[0], flagsBytes[1], flagsBytes[2],

		// 4: (reserved spacing) (4)
		0, 0, 0, 0,

		// 8: handler (4)
		'a', 'b', 'c', 'd',

		// 12: (reserved spacing) (12)
		// TODO(dustin): This is probably data that we need to add support for.
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,

		// 24: hdlrName (all remaining)
		't', 'e', 's', 't', 'n', 'a', 'm', 'e',
	}

	b := []byte{}
	pushBox(&b, "hdlr", hdlrBoxData)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := NewFile(sb, int64(len(b)))

	box, err := file.readBaseBox(0)
	log.PanicIf(err)

	cb, err := hdlrBoxFactory{}.New(box)
	log.PanicIf(err)

	hb := cb.(*HdlrBox)

	if hb.Version() != 0x0b {
		t.Fatalf("Version() not correct.")
	}

	if hb.Flags() != 0x11223300 {
		t.Fatalf("Flags() not correct: (0x%x)", hb.Flags())
	}

	if hb.Handler() != "abcd" {
		t.Fatalf("Handler() not correct: [%s]", hb.Handler())
	}

	if hb.HdlrName() != "testname" {
		t.Fatalf("HdlrName() not correct.")
	}

	if hb.String() != "hdlr<NAME=[hdlr] START=(0) SIZE=(40) VER=(0x0b) FLAGS=(0x11223300) HANDLER=[abcd] HDLR-NAME=[testname]>" {
		t.Fatalf("String() not correct: [%s]", hb.String())
	}

	if hb.InlineString() != "NAME=[hdlr] START=(0) SIZE=(40) VER=(0x0b) FLAGS=(0x11223300) HANDLER=[abcd] HDLR-NAME=[testname]" {
		t.Fatalf("InlineString() not correct: [%s]", hb.InlineString())
	}
}
