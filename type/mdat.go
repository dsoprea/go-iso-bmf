package bmftype

import (
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// MdatBox is the "Media Data" box.
//
// A container box which can hold the actual media data for a presentation
// (mdat).
type MdatBox struct {
	bmfcommon.Box
}

type mdatBoxFactory struct {
}

// Name returns the name of the type.
func (mdatBoxFactory) Name() string {
	return "mdat"
}

// New returns a new value instance.
func (mdatBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	mdatBox := &MdatBox{
		Box: box,
	}

	return mdatBox, nil
}

var (
	_ bmfcommon.BoxFactory = mdatBoxFactory{}
	_ bmfcommon.CommonBox  = &MdatBox{}
)

func init() {
	bmfcommon.RegisterBoxType(mdatBoxFactory{})
}
