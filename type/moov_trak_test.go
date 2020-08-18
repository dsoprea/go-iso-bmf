package bmftype

import (
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestTrakBox_SetLoadedBoxIndex(t *testing.T) {
	lbi := make(bmfcommon.Boxes, 0)

	trak := new(TrakBox)
	trak.SetLoadedBoxIndex(lbi)

	if reflect.DeepEqual(trak.LoadedBoxIndex, lbi.Index()) != true {
		t.Fatalf("SetLoadedBoxIndex() did not set the LBI correctly.")
	}
}

func TestTrakBoxFactory_Name(t *testing.T) {
	name := trakBoxFactory{}.Name()

	if name != "trak" {
		t.Fatalf("Name() not correct.")
	}
}

func TestTrakBoxFactory_New(t *testing.T) {
	b := []byte{}
	bmfcommon.PushBox(&b, "trak", nil)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file, err := bmfcommon.NewBmfResource(sb, int64(len(b)))
	log.PanicIf(err)

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := trakBoxFactory{}.New(box)
	log.PanicIf(err)

	// Nothing else we can validate.
	_, ok := cb.(*TrakBox)

	if ok != true {
		t.Fatalf("Expected an 'trak' box.")
	}
}
