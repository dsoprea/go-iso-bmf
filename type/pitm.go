package bmftype

import (
	"fmt"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// PitmBox is a "Handler Reference" box.
type PitmBox struct {
	bmfcommon.Box

	itemId uint32
}

// ItemId returns the primary item ID.
func (pitm *PitmBox) ItemId() uint32 {
	return pitm.itemId
}

// InlineString returns an undecorated string of field names and values.
func (pitm *PitmBox) InlineString() string {
	return fmt.Sprintf(
		"%s ID=(%d)",
		pitm.Box.InlineString(), pitm.itemId)
}

type pitmBoxFactory struct {
}

// Name returns the name of the type.
func (pitmBoxFactory) Name() string {
	return "pitm"
}

// New returns a new value instance.
func (pitmBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := box.ReadBoxData()
	log.PanicIf(err)

	version := data[0]

	var itemId uint32

	if version == 0 {
		itemId16 := bmfcommon.DefaultEndianness.Uint16(data[4:6])
		itemId = uint32(itemId16)
	} else if version == 1 {
		itemId = bmfcommon.DefaultEndianness.Uint32(data[4:8])
	} else {
		log.Panicf("version > 1 of PITM not supported")
	}

	pitm := &PitmBox{
		Box:    box,
		itemId: itemId,
	}

	return pitm, -1, nil
}

var (
	_ bmfcommon.BoxFactory = pitmBoxFactory{}
	_ bmfcommon.CommonBox  = &PitmBox{}
)

func init() {
	bmfcommon.RegisterBoxType(pitmBoxFactory{})
}
