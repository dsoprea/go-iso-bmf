package bmftype

import (
	"fmt"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// IrefBox is a "Item Reference" box.
type IrefBox struct {
	bmfcommon.Box

	version byte

	// LoadedBoxIndex contains this boxes children.
	bmfcommon.LoadedBoxIndex
}

// Version retrns the structural version of the IREF box.
func (iref IrefBox) Version() byte {

	// TODO(dustin): Add test

	return iref.version
}

// InlineString returns an undecorated string of field names and values.
func (iref *IrefBox) InlineString() string {

	// TODO(dustin): Add test

	return fmt.Sprintf(
		"%s",
		// "%s VER=(%d) FROM-ITEM-ID=(%d) TO-ITEM-IDS=(%d)[%v]",
		iref.Box.InlineString(),
		// iref.version, iref.fromItemId, len(iref.toItemIds), iref.toItemIds
	)
}

// SetLoadedBoxIndex sets the child boxes after a box has been manufactured
// and the children have been parsed. This allows parent boxes to be
// registered before the child boxes can look for them.
func (iref *IrefBox) SetLoadedBoxIndex(lbi bmfcommon.LoadedBoxIndex) {

	// TODO(dustin): !! Add test

	iref.LoadedBoxIndex = lbi
}

type irefBoxFactory struct {
}

// Name returns the name of the type.
func (irefBoxFactory) Name() string {

	// TODO(dustin): Add test

	return "iref"
}

// New returns a new value instance.
//
// This contains other boxes, but the box-types are actually the reference-
// types (e.g. cdsc)..
func (irefBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	data, err := box.ReadBoxData()
	log.PanicIf(err)

	version := data[0]

	iref := &IrefBox{
		Box:     box,
		version: version,
	}

	return iref, 4, nil
}

var (
	_ bmfcommon.BoxFactory = irefBoxFactory{}
	_ bmfcommon.CommonBox  = &IrefBox{}
)

func init() {
	bmfcommon.RegisterBoxType(irefBoxFactory{})
}
