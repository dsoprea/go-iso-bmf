package bmftype

import (
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestStblBox_SetLoadedBoxIndex(t *testing.T) {
	lbi := make(bmfcommon.Boxes, 0)

	stbl := new(StblBox)
	stbl.SetLoadedBoxIndex(lbi)

	if reflect.DeepEqual(stbl.LoadedBoxIndex, lbi.Index()) != true {
		t.Fatalf("SetLoadedBoxIndex() did not set the LBI correctly.")
	}
}

func TestStblBoxFactory_Name(t *testing.T) {
	name := stblBoxFactory{}.Name()

	if name != "stbl" {
		t.Fatalf("Name() not correct.")
	}
}

func TestStblBoxFactory_New(t *testing.T) {
	b := []byte{}
	bmfcommon.PushBox(&b, "stbl", nil)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file, err := bmfcommon.NewResource(sb, int64(len(b)))
	log.PanicIf(err)

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := stblBoxFactory{}.New(box)
	log.PanicIf(err)

	// Nothing else we can validate.
	_, ok := cb.(*StblBox)

	if ok != true {
		t.Fatalf("Expected an 'stbl' box.")
	}
}
