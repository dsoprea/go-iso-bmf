package bmftype

import (
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// MoovBox is the "Movie" box.
//
// The metadata for a presentation is stored in the single Movie bmfcommon.Box which occurs
// at the top-level of a file. Normally this box is close to the beginning or end
// of the file, though this is not required.
type MoovBox struct {
	bmfcommon.Box

	// TODO(dustin): Add test for this.
	isFragmented bool

	// LoadedBoxIndex contains this box's children.
	bmfcommon.LoadedBoxIndex
}

// IsFragmented return true if mvex box is present.
func (mv MoovBox) IsFragmented() bool {
	return mv.isFragmented
}

// SetLoadedBoxIndex sets the child boxes after a box has been manufactured
// and the children have been parsed. This allows parent boxes to be
// registered before the child boxes can look for them.
func (moov *MoovBox) SetLoadedBoxIndex(lbi bmfcommon.LoadedBoxIndex) {

	// TODO(dustin): !! Add test

	moov.LoadedBoxIndex = lbi
}

func (b *MoovBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	_, b.isFragmented = b.LoadedBoxIndex["mvex"]

	return nil
}

type moovBoxFactory struct {
}

// Name returns the name of the type.
func (moovBoxFactory) Name() string {
	return "moov"
}

// New returns a new value instance.
func (moovBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	moovBox := &MoovBox{
		Box: box,
	}

	return moovBox, 0, nil
}

var (
	_ bmfcommon.BoxFactory = moovBoxFactory{}
	_ bmfcommon.CommonBox  = &MoovBox{}
)

func init() {
	bmfcommon.RegisterBoxType(moovBoxFactory{})
}
