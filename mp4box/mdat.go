package mp4box

import (
	"github.com/dsoprea/go-logging"
)

// MdatBox is a "Media Data" box.
//
// A container box which can hold the actual media data for a presentation
// (mdat).
type MdatBox struct {
	Box
}

type mdatBoxFactory struct {
}

// Name returns the name of the type.
func (mdatBoxFactory) Name() string {
	return "mdat"
}

// New returns a new value instance.
func (mdatBoxFactory) New(box Box) (cb CommonBox, err error) {
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
	_ boxFactory = mdatBoxFactory{}
	_ CommonBox  = &MdatBox{}
)

func init() {
	registerAtom(mdatBoxFactory{})
}
