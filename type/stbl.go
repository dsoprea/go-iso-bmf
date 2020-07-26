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

// SetLoadedBoxIndex sets the child boxes after a box has been manufactured
// and the children have been parsed. This allows parent boxes to be
// registered before the child boxes can look for them.
func (stbl *StblBox) SetLoadedBoxIndex(lbi bmfcommon.LoadedBoxIndex) {

	// TODO(dustin): !! Add test

	stbl.LoadedBoxIndex = lbi
}

type stblBoxFactory struct {
}

// Name returns the name of the type.
func (stblBoxFactory) Name() string {
	return "stbl"
}

// New returns a new value instance.
func (stblBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	stblBox := &StblBox{
		Box: box,
	}

	return stblBox, 0, nil
}

var (
	_ bmfcommon.BoxFactory = stblBoxFactory{}
	_ bmfcommon.CommonBox  = &StblBox{}
)

func init() {
	bmfcommon.RegisterBoxType(stblBoxFactory{})
}
