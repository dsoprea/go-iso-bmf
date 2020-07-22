package bmftype

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
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
	box := bmfcommon.NewBox("abcd", 1234, 5678, nil)

	hb := HdlrBox{
		Box:      box,
		version:  11,
		flags:    11,
		handler:  "handler_test",
		hdlrName: "name_test",
	}

	if hb.String() != "hdlr<NAME=[abcd] PARENT=[ROOT] START=(0x00000000000004d2) SIZE=(5678) VER=(0x0b) FLAGS=(0x0000000b) HANDLER=[handler_test] HDLR-NAME=[name_test]>" {
		t.Fatalf("String() not correct: [%s]", hb.String())
	}
}

func TestHdlrBox_InlineString(t *testing.T) {
	box := bmfcommon.NewBox("abcd", 1234, 5678, nil)

	hb := HdlrBox{
		Box:      box,
		version:  11,
		flags:    11,
		handler:  "handler_test",
		hdlrName: "name_test",
	}

	if hb.InlineString() != "NAME=[abcd] PARENT=[ROOT] START=(0x00000000000004d2) SIZE=(5678) VER=(0x0b) FLAGS=(0x0000000b) HANDLER=[handler_test] HDLR-NAME=[name_test]" {
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

	var hdlrBoxData []byte

	// Version and flags.
	bmfcommon.PushBytes(&hdlrBoxData, uint32(0x11223344))

	// Reserved spacing.
	bmfcommon.PushBytes(&hdlrBoxData, uint32(0))

	// Handler name
	bmfcommon.PushBytes(&hdlrBoxData, []byte{'a', 'b', 'c', 'd'})

	// Reserved spacing.
	// TODO(dustin): This is probably data that we need to add support for.
	bmfcommon.PushBytes(&hdlrBoxData, uint32(0))
	bmfcommon.PushBytes(&hdlrBoxData, uint32(0))
	bmfcommon.PushBytes(&hdlrBoxData, uint32(0))

	// handler name (all remaining)
	// TODO(dustin): Update this comment to not be a duplicate.
	bmfcommon.PushBytes(&hdlrBoxData, []byte{'t', 'e', 's', 't', 'n', 'a', 'm', 'e'})

	b := []byte{}
	bmfcommon.PushBox(&b, "hdlr", hdlrBoxData)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := bmfcommon.NewFile(sb, int64(len(b)))

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, err := hdlrBoxFactory{}.New(box)
	log.PanicIf(err)

	hb := cb.(*HdlrBox)

	if hb.Version() != 0x11 {
		t.Fatalf("Version() not correct.")
	}

	if hb.Flags() != 0x11223344 {
		t.Fatalf("Flags() not correct: (0x%x)", hb.Flags())
	}

	if hb.Handler() != "abcd" {
		t.Fatalf("Handler() not correct: [%s]", hb.Handler())
	}

	if hb.HdlrName() != "testname" {
		t.Fatalf("HdlrName() not correct.")
	}

	if hb.String() != "hdlr<NAME=[hdlr] PARENT=[ROOT] START=(0x0000000000000000) SIZE=(40) VER=(0x11) FLAGS=(0x11223344) HANDLER=[abcd] HDLR-NAME=[testname]>" {
		t.Fatalf("String() not correct: [%s]", hb.String())
	}

	if hb.InlineString() != "NAME=[hdlr] PARENT=[ROOT] START=(0x0000000000000000) SIZE=(40) VER=(0x11) FLAGS=(0x11223344) HANDLER=[abcd] HDLR-NAME=[testname]" {
		t.Fatalf("InlineString() not correct: [%s]", hb.InlineString())
	}
}
