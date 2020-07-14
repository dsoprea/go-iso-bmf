package atom

import (
	"github.com/dsoprea/go-logging"
)

// StsdBox is a "Sample Description" box.
type StsdBox struct {
	Box

	version byte
	flags   uint32

	LoadedBoxIndex
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

	data, err := b.readBoxData()
	log.PanicIf(err)

	b.version = data[0]
	b.flags = defaultEndianness.Uint32(data[0:4])

	// Skip extra 8 bytes.
	boxes, err := b.Box.readBoxes(8)
	log.PanicIf(err)

	b.LoadedBoxIndex = boxes.Index()

	return nil
}

type stsdBoxFactory struct {
}

// Name returns the name of the type.
func (stsdBoxFactory) Name() string {
	return "stsd"
}

// New returns a new value instance.
func (stsdBoxFactory) New(box Box) (cb CommonBox, err error) {
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

	return stsdBox, nil
}

var (
	_ boxFactory = stsdBoxFactory{}
	_ CommonBox  = &StsdBox{}
)

func init() {
	registerAtom(stsdBoxFactory{})
}
