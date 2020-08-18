package bmftype

import (
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestIdatBox_InlineString(t *testing.T) {
	data := []byte{1, 2, 3}

	idat := &IdatBox{
		data: data,
	}

	if idat.InlineString() != "NAME=[] PARENT=[ROOT] START=(0x0000000000000000) SIZE=(0) DATA-SIZE=(3)" {
		t.Fatalf("InlineString() is not correct: [%s]", idat.InlineString())
	}
}

func TestIdatBoxFactory_Name(t *testing.T) {
	factory := idatBoxFactory{}
	if factory.Name() != "idat" {
		t.Fatalf("Name() not correct.")
	}
}

func TestIdatBoxFactory_New(t *testing.T) {

	var b []byte
	bmfcommon.PushBox(&b, "idat", []byte{})

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file, err := bmfcommon.NewBmfResource(sb, int64(len(b)))
	log.PanicIf(err)

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, _, err := idatBoxFactory{}.New(box)
	log.PanicIf(err)

	_ = cb.(*IdatBox)
}
