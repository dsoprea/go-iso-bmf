package bmftype

import (
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// MdiaBox is the "Media" box.
//
// The media declaration container contains all the objects that declare information
// about the media data within a track.
type MdiaBox struct {
	// Box is the base box.
	bmfcommon.Box

	// LoadedBoxIndex contains this box's children.
	bmfcommon.LoadedBoxIndex
}

// SetLoadedBoxIndex sets the child boxes after a box has been manufactured
// and the children have been parsed. This allows parent boxes to be
// registered before the child boxes can look for them.
func (mdia *MdiaBox) SetLoadedBoxIndex(lbi bmfcommon.LoadedBoxIndex) {

	// TODO(dustin): !! Add test

	mdia.LoadedBoxIndex = lbi
}

type mdiaBoxFactory struct {
}

// Name returns the name of the type.
func (mdiaBoxFactory) Name() string {
	return "mdia"
}

// New returns a new value instance.
func (mdiaBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	mdiaBox := &MdiaBox{
		Box: box,
	}

	return mdiaBox, 0, nil
}

var (
	_ bmfcommon.BoxFactory = mdiaBoxFactory{}
	_ bmfcommon.CommonBox  = &MdiaBox{}
)

func init() {
	bmfcommon.RegisterBoxType(mdiaBoxFactory{})
}
