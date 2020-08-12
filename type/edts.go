package bmftype

import (
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// EdtsBox is the "Edit" box.
type EdtsBox struct {
	bmfcommon.Box

	// LoadedBoxIndex contains this box's children.
	bmfcommon.LoadedBoxIndex
}

// SetLoadedBoxIndex sets the child boxes after a box has been manufactured
// and the children have been parsed. This allows parent boxes to be
// registered before the child boxes can look for them.
func (edts *EdtsBox) SetLoadedBoxIndex(lbi bmfcommon.LoadedBoxIndex) {

	// TODO(dustin): !! Add test

	edts.LoadedBoxIndex = lbi
}

type edtsBoxFactory struct {
}

// Name returns the name of the type.
func (edtsBoxFactory) Name() string {
	return "edts"
}

// New returns a new value instance.
func (edtsBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	edtsBox := &EdtsBox{
		Box: box,
	}

	return edtsBox, 0, nil
}

var (
	_ bmfcommon.BoxFactory = edtsBoxFactory{}
	_ bmfcommon.CommonBox  = &EdtsBox{}
)

func init() {
	bmfcommon.RegisterBoxType(edtsBoxFactory{})
}
