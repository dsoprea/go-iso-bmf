package bmftype

import (
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// MdatBox is the "Media Data" box.
//
// A container box which can hold all of the actual media data. This is just a
// big space where the EXIF and image data offsets refer and has little value in
// being directly referenced.
type MdatBox struct {
	bmfcommon.Box
}

type mdatBoxFactory struct {
}

// Name returns the name of the type.
func (mdatBoxFactory) Name() string {
	return "mdat"
}

// New returns a new value instance. Since mdat is just the general space where
// all data referred to by everything is hosted, we don't capture it or directly
// referenced it. It's not generally useful.
func (mdatBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	mdatBox := &MdatBox{
		Box: box,
	}

	return mdatBox, -1, nil
}

var (
	_ bmfcommon.BoxFactory = mdatBoxFactory{}
	_ bmfcommon.CommonBox  = &MdatBox{}
)

func init() {
	bmfcommon.RegisterBoxType(mdatBoxFactory{})
}
