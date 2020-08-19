package bmftype

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestPitmBox_ItemId(t *testing.T) {
	pitm := &PitmBox{
		itemId: 11,
	}

	if pitm.ItemId() != 11 {
		t.Fatalf("ItemId() not correct.")
	}
}

func TestPitmBox_InlineString(t *testing.T) {
	pitm := &PitmBox{
		itemId: 11,
	}

	if pitm.InlineString() != "NAME=[] PARENT=[ROOT] START=(0x0000000000000000) SIZE=(0) ID=(11)" {
		t.Fatalf("InlineString() not correct: [%s]", pitm.InlineString())
	}
}

func TestPitmBoxFactory_Name(t *testing.T) {
	factory := pitmBoxFactory{}
	if factory.Name() != "pitm" {
		t.Fatalf("Name() not correct.")
	}
}

func TestPitmBoxFactory_New_Version0(t *testing.T) {
	var data []byte

	version := uint8(0)
	bmfcommon.PushBytes(&data, []byte{version, 0, 0, 0})

	itemId := uint16(11)
	bmfcommon.PushBytes(&data, itemId)

	b := []byte{}
	bmfcommon.PushBox(&b, "pitm", data)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file, err := bmfcommon.NewResource(sb, 0)
	log.PanicIf(err)

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := pitmBoxFactory{}.New(box)
	log.PanicIf(err)

	pitm := cb.(*PitmBox)

	if pitm.ItemId() != uint32(itemId) {
		t.Fatalf("ItemId() not correct.")
	}
}

func TestPitmBoxFactory_New_Version1(t *testing.T) {
	var data []byte

	version := uint8(1)
	bmfcommon.PushBytes(&data, []byte{version, 0, 0, 0})

	itemId := uint32(11)
	bmfcommon.PushBytes(&data, itemId)

	b := []byte{}
	bmfcommon.PushBox(&b, "pitm", data)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file, err := bmfcommon.NewResource(sb, 0)
	log.PanicIf(err)

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := pitmBoxFactory{}.New(box)
	log.PanicIf(err)

	pitm := cb.(*PitmBox)

	if pitm.ItemId() != itemId {
		t.Fatalf("ItemId() not correct.")
	}
}
