package bmftype

import (
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// MinfBox is the "Media Information" box.
//
// This box contains all the objects that declare characteristics information of
// the media in the track.
type MinfBox struct {
	bmfcommon.Box

	// LoadedBoxIndex contains this boxes children.
	bmfcommon.LoadedBoxIndex
}

func (b *MinfBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	boxes, err := b.Box.ReadBoxes(0, b)
	log.PanicIf(err)

	b.LoadedBoxIndex = boxes.Index()

	return nil
}

type minfBoxFactory struct {
}

// Name returns the name of the type.
func (minfBoxFactory) Name() string {
	return "minf"
}

// New returns a new value instance.
func (minfBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	minfBox := &MinfBox{
		Box: box,
	}

	err = minfBox.parse()
	log.PanicIf(err)

	return minfBox, nil
}

var (
	_ bmfcommon.BoxFactory = minfBoxFactory{}
	_ bmfcommon.CommonBox  = &MinfBox{}
)

func init() {
	bmfcommon.RegisterBoxType(minfBoxFactory{})
}
