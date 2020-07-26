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

// SetLoadedBoxIndex sets the child boxes after a box has been manufactured
// and the children have been parsed. This allows parent boxes to be
// registered before the child boxes can look for them.
func (meta *MetaBox) SetLoadedBoxIndex(lbi bmfcommon.LoadedBoxIndex) {

	// TODO(dustin): !! Add test

	meta.LoadedBoxIndex = lbi
}

type metaBoxFactory struct {
}

// Name returns the name of the type.
func (metaBoxFactory) Name() string {

	// TODO(dustin): Add test

	return "meta"
}

// New returns a new value instance.
func (metaBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, childBoxSeriesOffset int, err error) {
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

	return metaBox, 4, nil
}

var (
	_ bmfcommon.BoxFactory = metaBoxFactory{}
	_ bmfcommon.CommonBox  = &MetaBox{}
)

func init() {
	bmfcommon.RegisterBoxType(metaBoxFactory{})
}
