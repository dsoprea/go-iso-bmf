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

// SetLoadedBoxIndex sets the child boxes after a box has been manufactured
// and the children have been parsed. This allows parent boxes to be
// registered before the child boxes can look for them.
func (minf *MinfBox) SetLoadedBoxIndex(lbi bmfcommon.LoadedBoxIndex) {

	// TODO(dustin): !! Add test

	minf.LoadedBoxIndex = lbi
}

type minfBoxFactory struct {
}

// Name returns the name of the type.
func (minfBoxFactory) Name() string {
	return "minf"
}

// New returns a new value instance.
func (minfBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	minfBox := &MinfBox{
		Box: box,
	}

	return minfBox, 0, nil
}

var (
	_ bmfcommon.BoxFactory = minfBoxFactory{}
	_ bmfcommon.CommonBox  = &MinfBox{}
)

func init() {
	bmfcommon.RegisterBoxType(minfBoxFactory{})
}
