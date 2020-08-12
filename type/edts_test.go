package bmftype

import (
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestEdtsBox_SetLoadedBoxIndex(t *testing.T) {
	lbi := make(bmfcommon.LoadedBoxIndex)

	edts := new(EdtsBox)
	edts.SetLoadedBoxIndex(lbi)

	if reflect.DeepEqual(edts.LoadedBoxIndex, lbi) != true {
		t.Fatalf("SetLoadedBoxIndex() did not set the LBI correctly.")
	}
}

func TestEdtsBoxFactory_Name(t *testing.T) {
	ebf := edtsBoxFactory{}

	if ebf.Name() != "edts" {
		t.Fatalf("Name() not correct.")
	}
}

func TestEdtsBoxFactory_New(t *testing.T) {
	var b []byte
	bmfcommon.PushBox(&b, "edts", nil)

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := bmfcommon.NewBmfResource(sb, int64(len(b)))

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := edtsBoxFactory{}.New(box)
	log.PanicIf(err)

	// Nothing else we can validate.
	_, ok := cb.(*EdtsBox)

	if ok != true {
		t.Fatalf("Expected an 'edts' box.")
	}
}
