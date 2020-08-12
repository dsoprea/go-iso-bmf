package bmftype

import (
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// StsdBox is the "Sample Description" box.
type StsdBox struct {
	bmfcommon.Box

	version byte
	flags   uint32

	// LoadedBoxIndex contains this box's children.
	bmfcommon.LoadedBoxIndex
}

func (sb *StsdBox) Version() byte {
	return sb.version
}

func (sb *StsdBox) Flags() uint32 {
	return sb.flags
}

func (b *StsdBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.ReadBoxData()
	log.PanicIf(err)

	b.version = data[0]
	b.flags = bmfcommon.DefaultEndianness.Uint32(data[0:4])

	return nil
}

// SetLoadedBoxIndex sets the child boxes after a box has been manufactured
// and the children have been parsed. This allows parent boxes to be
// registered before the child boxes can look for them.
func (stsd *StsdBox) SetLoadedBoxIndex(lbi bmfcommon.LoadedBoxIndex) {

	// TODO(dustin): !! Add test

	stsd.LoadedBoxIndex = lbi
}

type stsdBoxFactory struct {
}

// Name returns the name of the type.
func (stsdBoxFactory) Name() string {
	return "stsd"
}

// New returns a new value instance.
func (stsdBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	stsdBox := &StsdBox{
		Box: box,
	}

	err = stsdBox.parse()
	log.PanicIf(err)

	return stsdBox, 8, nil
}

var (
	_ bmfcommon.BoxFactory = stsdBoxFactory{}
	_ bmfcommon.CommonBox  = &StsdBox{}
)

func init() {
	bmfcommon.RegisterBoxType(stsdBoxFactory{})
}
