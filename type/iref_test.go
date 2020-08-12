package bmftype

import (
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestIrefBox_Version(t *testing.T) {
	iref := &IrefBox{
		version: 11,
	}

	if iref.Version() != 11 {
		t.Fatalf("Version() not correct.")
	}
}

func TestIrefBox_InlineString(t *testing.T) {
	iref := &IrefBox{
		version: 11,
	}

	if iref.InlineString() != "NAME=[] PARENT=[ROOT] START=(0x0000000000000000) SIZE=(0)" {
		t.Fatalf("InlineString() not correct: [%s]", iref.InlineString())
	}
}

func TestIrefBox_SetLoadedBoxIndex(t *testing.T) {
	lbi := make(bmfcommon.LoadedBoxIndex)

	iref := new(IrefBox)
	iref.SetLoadedBoxIndex(lbi)

	if reflect.DeepEqual(iref.LoadedBoxIndex, lbi) != true {
		t.Fatalf("SetLoadedBoxIndex() did not set the LBI correctly.")
	}
}

func TestIrefBoxFactory_New(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			log.PrintError(err)

			t.Fatalf("Test failed.")
		}
	}()

	data := []byte{
		0, 0, 0, 0,
	}

	var b []byte
	bmfcommon.PushBox(&b, "iref", data)

	sb := rifs.NewSeekableBufferWithBytes(b)

	// Use zero length to prevent immediate parsing.
	resource := bmfcommon.NewBmfResource(sb, 0)

	headerSize := int64(8)
	box := bmfcommon.NewBox("iref", 0, headerSize+int64(len(data)), headerSize, resource)

	factory := irefBoxFactory{}

	cb, childrenOffset, err := factory.New(box)
	log.PanicIf(err)

	if childrenOffset != 4 {
		t.Fatalf("Children offset not correct.")
	}

	iref := cb.(*IrefBox)

	if iref.Version() != 0 {
		t.Fatalf("Parsed failed.")
	}
}

func TestIrefBoxFactory_Name(t *testing.T) {
	factory := irefBoxFactory{}

	if factory.Name() != "iref" {
		t.Fatalf("Name() not correct.")
	}
}
