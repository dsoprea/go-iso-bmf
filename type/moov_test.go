package bmftype

import (
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestMoovBox_IsFragmented_True(t *testing.T) {
	moov := &MoovBox{
		isFragmented: true,
	}

	if moov.IsFragmented() != true {
		t.Fatalf("IsFragmented() not correct.")
	}
}

func TestMoovBox_IsFragmented_False(t *testing.T) {
	moov := &MoovBox{
		isFragmented: false,
	}

	if moov.IsFragmented() != false {
		t.Fatalf("IsFragmented() not correct.")
	}
}

func TestMoovBox_SetLoadedBoxIndex(t *testing.T) {
	lbi := make(bmfcommon.LoadedBoxIndex)

	moov := new(MoovBox)
	moov.SetLoadedBoxIndex(lbi)

	if reflect.DeepEqual(moov.LoadedBoxIndex, lbi) != true {
		t.Fatalf("SetLoadedBoxIndex() did not set the LBI correctly.")
	}
}

func TestMoovBoxFactory_Name(t *testing.T) {
	name := moovBoxFactory{}.Name()

	if name != "moov" {
		t.Fatalf("Name() not correct.")
	}
}

func TestMoovBoxFactory_New(t *testing.T) {
	b := []byte{}
	bmfcommon.PushBox(&b, "moov", nil)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file, err := bmfcommon.NewBmfResource(sb, int64(len(b)))
	log.PanicIf(err)

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := moovBoxFactory{}.New(box)
	log.PanicIf(err)

	// Nothing else we can validate.
	_, ok := cb.(*MoovBox)

	if ok != true {
		t.Fatalf("Expected an 'moov' box.")
	}
}
