package bmftype

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestVmhdBox_Version(t *testing.T) {
	vb := VmhdBox{
		version: 0x11,
	}

	if vb.Version() != 0x11 {
		t.Fatalf("Version() not correct.")
	}
}

func TestVmhdBox_Flags(t *testing.T) {
	vb := VmhdBox{
		flags: 0x11223344,
	}

	if vb.Flags() != 0x11223344 {
		t.Fatalf("Flags() not correct.")
	}
}

func TestVmhdBox_GraphicsMode(t *testing.T) {
	vb := VmhdBox{
		graphicsMode: 0x22,
	}

	if vb.GraphicsMode() != 0x22 {
		t.Fatalf("GraphicsMode() not correct.")
	}
}

func TestVmhdBox_OpColor(t *testing.T) {
	vb := VmhdBox{
		opColor: 0x33,
	}

	if vb.OpColor() != 0x33 {
		t.Fatalf("OpColor() not correct.")
	}
}

func TestVmhdBoxFactory_Name(t *testing.T) {
	vbf := vmhdBoxFactory{}
	if vbf.Name() != "vmhd" {
		t.Fatalf("Name() not correct.")
	}
}

func TestVmhdBoxFactory_New(t *testing.T) {
	var data []byte

	// flags
	bmfcommon.PushBytes(&data, uint32(0x11223344))

	// graphicsMode
	bmfcommon.PushBytes(&data, uint16(0x55))

	// opColor
	bmfcommon.PushBytes(&data, uint16(0x66))

	b := []byte{}
	bmfcommon.PushBox(&b, "vmhd", data)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := bmfcommon.NewFile(sb, int64(len(b)))

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := vmhdBoxFactory{}.New(box)
	log.PanicIf(err)

	vb := cb.(*VmhdBox)

	if vb.Version() != 0x11 {
		t.Fatalf("Version() not correct: (0x%02x)", vb.Version())
	}

	if vb.Flags() != 0x11223344 {
		t.Fatalf("Flags() not correct.")
	}

	if vb.GraphicsMode() != 0x55 {
		t.Fatalf("GraphicsMode() not correct.")
	}

	if vb.OpColor() != 0x66 {
		t.Fatalf("OpColor() not correct.")
	}
}
