package bmftype

import (
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// EdtsBox is the "Edit" box.
type EdtsBox struct {
	bmfcommon.Box

	// LoadedBoxIndex contains this boxes children.
	bmfcommon.LoadedBoxIndex
}

func (b *EdtsBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	boxes, err := b.Box.ReadBoxes(0)
	log.PanicIf(err)

	b.LoadedBoxIndex = boxes.Index()

	return nil
}

type edtsBoxFactory struct {
}

// Name returns the name of the type.
func (edtsBoxFactory) Name() string {
	return "edts"
}

// New returns a new value instance.
func (edtsBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	edtsBox := &EdtsBox{
		Box: box,
	}

	err = edtsBox.parse()
	log.PanicIf(err)

	return edtsBox, nil
}

var (
	_ bmfcommon.BoxFactory = edtsBoxFactory{}
	_ bmfcommon.CommonBox  = &EdtsBox{}
)

func init() {
	bmfcommon.RegisterBoxType(edtsBoxFactory{})
}
