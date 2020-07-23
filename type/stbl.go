package bmftype

import (
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// StblBox is the "Sample Table" box.
type StblBox struct {
	bmfcommon.Box

	// LoadedBoxIndex contains this boxes children.
	bmfcommon.LoadedBoxIndex
}

func (b *StblBox) parse() (err error) {
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

type stblBoxFactory struct {
}

// Name returns the name of the type.
func (stblBoxFactory) Name() string {
	return "stbl"
}

// New returns a new value instance.
func (stblBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	stblBox := &StblBox{
		Box: box,
	}

	err = stblBox.parse()
	log.PanicIf(err)

	return stblBox, nil
}

var (
	_ bmfcommon.BoxFactory = stblBoxFactory{}
	_ bmfcommon.CommonBox  = &StblBox{}
)

func init() {
	bmfcommon.RegisterBoxType(stblBoxFactory{})
}
