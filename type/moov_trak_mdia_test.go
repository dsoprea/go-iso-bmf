package bmftype

import (
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestMdiaBox_SetLoadedBoxIndex(t *testing.T) {
	lbi := make(bmfcommon.Boxes, 0)

	mdia := new(MdiaBox)
	mdia.SetLoadedBoxIndex(lbi)

	if reflect.DeepEqual(mdia.LoadedBoxIndex, lbi.Index()) != true {
		t.Fatalf("SetLoadedBoxIndex() did not set the LBI correctly.")
	}
}

func TestMdiaBoxFactory_Name(t *testing.T) {
	name := mdiaBoxFactory{}.Name()

	if name != "mdia" {
		t.Fatalf("Name() not correct.")
	}
}

func TestMdiaBoxFactory_New(t *testing.T) {
	b := []byte{}
	bmfcommon.PushBox(&b, "mdia", nil)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file, err := bmfcommon.NewBmfResource(sb, int64(len(b)))
	log.PanicIf(err)

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := mdiaBoxFactory{}.New(box)
	log.PanicIf(err)

	// Nothing else we can validate.
	_, ok := cb.(*MdiaBox)

	if ok != true {
		t.Fatalf("Expected an 'mdia' box.")
	}
}
