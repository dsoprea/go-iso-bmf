package bmftype

import (
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// MetaBox is the "Meta" box.
type MetaBox struct {
	// Box is the base inner box.
	bmfcommon.Box

	// LoadedBoxIndex contains this boxes children.
	bmfcommon.LoadedBoxIndex
}

type metaBoxFactory struct {
}

// Name returns the name of the type.
func (metaBoxFactory) Name() string {

	// TODO(dustin): Add test

	return "meta"
}

// New returns a new value instance.
func (metaBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	// Boxes follow the four-bytes with the version/flags.

	metaBox := &MetaBox{
		Box: box,
	}

	boxes, err := box.ReadBoxes(4, metaBox)
	log.PanicIf(err)

	metaBox.LoadedBoxIndex = boxes.Index()

	return metaBox, nil
}

var (
	_ bmfcommon.BoxFactory = metaBoxFactory{}
	_ bmfcommon.CommonBox  = &MetaBox{}
)

func init() {
	bmfcommon.RegisterBoxType(metaBoxFactory{})
}
