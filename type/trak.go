package bmftype

import (
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// TrakBox is the "Track" box.
type TrakBox struct {
	bmfcommon.Box

	// SamplesDuration
	// SamplesSize
	// SampleGroupsInfo

	// chunks []Chunk
	// samples []Sample

	// LoadedBoxIndex contains this box's children.
	bmfcommon.LoadedBoxIndex
}

// SetLoadedBoxIndex sets the child boxes after a box has been manufactured
// and the children have been parsed. This allows parent boxes to be
// registered before the child boxes can look for them.
func (trak *TrakBox) SetLoadedBoxIndex(lbi bmfcommon.LoadedBoxIndex) {
	trak.LoadedBoxIndex = lbi
}

type trakBoxFactory struct {
}

// Name returns the name of the type.
func (trakBoxFactory) Name() string {
	return "trak"
}

// New returns a new value instance.
func (trakBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Return to this once we have a sample.

	trakBox := &TrakBox{
		Box: box,
	}

	return trakBox, -1, nil
}

var (
	_ bmfcommon.BoxFactory = trakBoxFactory{}
	_ bmfcommon.CommonBox  = &TrakBox{}
)

func init() {
	bmfcommon.RegisterBoxType(trakBoxFactory{})
}
