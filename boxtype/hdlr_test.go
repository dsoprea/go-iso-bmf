package boxtype

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

	var hdlrBoxData []byte

	// Version and flags.
	pushBytes(&hdlrBoxData, uint32(0x11223344))

	// Reserved spacing.
	pushBytes(&hdlrBoxData, uint32(0))

	// Handler name
	pushBytes(&hdlrBoxData, []byte{'a', 'b', 'c', 'd'})

	// Reserved spacing.
	// TODO(dustin): This is probably data that we need to add support for.
	pushBytes(&hdlrBoxData, uint32(0))
	pushBytes(&hdlrBoxData, uint32(0))
	pushBytes(&hdlrBoxData, uint32(0))

	// handler name (all remaining)
	// TODO(dustin): Update this comment to not be a duplicate.
	pushBytes(&hdlrBoxData, []byte{'t', 'e', 's', 't', 'n', 'a', 'm', 'e'})

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

	if hb.String() != "hdlr<NAME=[hdlr] START=(0) SIZE=(40) VER=(0x11) FLAGS=(0x11223344) HANDLER=[abcd] HDLR-NAME=[testname]>" {
		t.Fatalf("String() not correct: [%s]", hb.String())
	}

	if hb.InlineString() != "NAME=[hdlr] START=(0) SIZE=(40) VER=(0x11) FLAGS=(0x11223344) HANDLER=[abcd] HDLR-NAME=[testname]" {
		t.Fatalf("InlineString() not correct: [%s]", hb.InlineString())
	}
}
