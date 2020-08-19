package bmftype

import (
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestMetaBox_SetLoadedBoxIndex(t *testing.T) {
	lbi := make(bmfcommon.Boxes, 0)

	meta := new(MetaBox)
	meta.SetLoadedBoxIndex(lbi)

	if reflect.DeepEqual(meta.LoadedBoxIndex, lbi.Index()) != true {
		t.Fatalf("SetLoadedBoxIndex() did not set the LBI correctly.")
	}
}

func TestMetaBoxFactory_Name(t *testing.T) {
	name := metaBoxFactory{}.Name()

	if name != "meta" {
		t.Fatalf("Name() not correct.")
	}
}

func TestMetaBoxFactory_New(t *testing.T) {
	b := []byte{}
	bmfcommon.PushBox(&b, "meta", nil)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file, err := bmfcommon.NewResource(sb, int64(len(b)))
	log.PanicIf(err)

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := metaBoxFactory{}.New(box)
	log.PanicIf(err)

	// Nothing else we can validate.
	_, ok := cb.(*MetaBox)

	if ok != true {
		t.Fatalf("Expected an 'meta' box.")
	}
}
