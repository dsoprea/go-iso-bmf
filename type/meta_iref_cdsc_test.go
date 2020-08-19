package bmftype

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestCdscBoxFactory_Name(t *testing.T) {
	factory := cdscBoxFactory{}

	if factory.Name() != "cdsc" {
		t.Fatalf("Name() not correct.")
	}
}

func TestCdscBox_New_Version0_WithReferences(t *testing.T) {
	var b []byte

	// Push CDSC box.

	dataCdsc := make([]byte, 0)

	// Push from-ID.
	bmfcommon.PushBytes(&dataCdsc, uint16(11))

	// Reference count.
	bmfcommon.PushBytes(&dataCdsc, uint16(2))

	// References
	bmfcommon.PushBytes(&dataCdsc, uint16(22))
	bmfcommon.PushBytes(&dataCdsc, uint16(33))

	bmfcommon.PushBox(&b, "cdsc", dataCdsc)

	// Parse stream. We do a formal parse so that other supporting boxes that
	// might've been pushed will be indexed.

	sb := rifs.NewSeekableBufferWithBytes(b)

	// Use zero length to prevent immediate parsing.
	resource, err := bmfcommon.NewResource(sb, 0)
	log.PanicIf(err)

	iref := &IrefBox{
		version: 0,
	}

	ibe := bmfcommon.IndexedBoxEntry{"meta.iref", 0}
	resource.Index()[ibe] = iref

	headerSize := int64(8)
	box := bmfcommon.NewBox("cdsc", 0, headerSize+int64(len(dataCdsc)), headerSize, resource)

	factory := cdscBoxFactory{}

	cb, childrenOffset, err := factory.New(box)
	log.PanicIf(err)

	if childrenOffset != -1 {
		t.Fatalf("Children offset not correct.")
	}

	cdsc := cb.(*CdscBox)

	if cdsc.InlineString() != "NAME=[cdsc] PARENT=[ROOT] START=(0x0000000000000000) SIZE=(16) VER=(0) FROM-ITEM-ID=(11) TO-ITEM-IDS=(2)[22,33]" {
		t.Fatalf("InlineString() not correct: [%s]", cdsc.InlineString())
	}
}

func TestCdscBox_New_Version1_WithReferences(t *testing.T) {
	var b []byte

	// Push CDSC box.

	dataCdsc := make([]byte, 0)

	// Push from-ID.
	bmfcommon.PushBytes(&dataCdsc, uint32(11))

	// Reference count.
	bmfcommon.PushBytes(&dataCdsc, uint16(2))

	// References
	bmfcommon.PushBytes(&dataCdsc, uint32(22))
	bmfcommon.PushBytes(&dataCdsc, uint32(33))

	bmfcommon.PushBox(&b, "cdsc", dataCdsc)

	// Parse stream. We do a formal parse so that other supporting boxes that
	// might've been pushed will be indexed.

	sb := rifs.NewSeekableBufferWithBytes(b)

	// Use zero length to prevent immediate parsing.
	resource, err := bmfcommon.NewResource(sb, 0)
	log.PanicIf(err)

	iref := &IrefBox{
		version: 1,
	}

	ibe := bmfcommon.IndexedBoxEntry{"meta.iref", 0}
	resource.Index()[ibe] = iref

	headerSize := int64(8)
	box := bmfcommon.NewBox("cdsc", 0, headerSize+int64(len(dataCdsc)), headerSize, resource)

	factory := cdscBoxFactory{}

	cb, childrenOffset, err := factory.New(box)
	log.PanicIf(err)

	if childrenOffset != -1 {
		t.Fatalf("Children offset not correct.")
	}

	cdsc := cb.(*CdscBox)

	if cdsc.InlineString() != "NAME=[cdsc] PARENT=[ROOT] START=(0x0000000000000000) SIZE=(22) VER=(0) FROM-ITEM-ID=(11) TO-ITEM-IDS=(2)[22,33]" {
		t.Fatalf("InlineString() not correct: [%s]", cdsc.InlineString())
	}
}
