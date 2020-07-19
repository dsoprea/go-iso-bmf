package bmftype

import (
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// Avc1Box defines the avc1 box structure.
type Avc1Box struct {
	bmfcommon.Box

	version byte
}

// Version is the version.
func (b Avc1Box) Version() byte {
	return b.version
}

type avc1BoxFactory struct {
}

// Name returns the name of the type.
func (avc1BoxFactory) Name() string {
	return "avc1"
}

// New returns a new value instance.
func (avc1BoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	avc1Box := &Avc1Box{
		Box: box,
	}

	data, err := avc1Box.ReadBoxData()
	log.PanicIf(err)

	avc1Box.version = data[0]

	return avc1Box, nil
}

var (
	_ bmfcommon.BoxFactory = avc1BoxFactory{}
	_ bmfcommon.CommonBox  = &Avc1Box{}
)

func init() {
	bmfcommon.RegisterBoxType(avc1BoxFactory{})
}
