package boxtype

import (
	"github.com/dsoprea/go-logging"
)

// Avc1Box defines the avc1 box structure.
type Avc1Box struct {
	Box

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
func (avc1BoxFactory) New(box Box) (cb CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	avc1Box := &Avc1Box{
		Box: box,
	}

	data, err := avc1Box.readBoxData()
	log.PanicIf(err)

	avc1Box.version = data[0]

	return avc1Box, nil
}

var (
	_ boxFactory = avc1BoxFactory{}
	_ CommonBox  = &Avc1Box{}
)

func init() {
	RegisterBoxType(avc1BoxFactory{})
}