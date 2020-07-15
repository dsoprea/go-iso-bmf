package mp4box

import (
	"github.com/dsoprea/go-logging"
)

// MvhdBox is a "Movie Header" box.
//
// This box defines overall information which is media-independent,
// and relevant to the entire presentationconsidered as a whole.
type MvhdBox struct {
	Box

	flags   uint32
	version uint8
	// creationTime     uint32
	// modificationTime uint32
	timescale uint32
	duration  uint32
	rate      Fixed32
	volume    Fixed16
}

func (mb *MvhdBox) Flags() uint32 {
	return mb.flags
}

func (mb *MvhdBox) Version() uint8 {
	return mb.version
}

func (mb *MvhdBox) Timescale() uint32 {
	return mb.timescale
}

func (mb *MvhdBox) Duration() uint32 {
	return mb.duration
}

func (mb *MvhdBox) Rate() Fixed32 {
	return mb.rate
}

func (mb *MvhdBox) Volume() Fixed16 {
	return mb.volume
}

func (b *MvhdBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.readBoxData()
	log.PanicIf(err)

	b.version = data[0]
	b.timescale = defaultEndianness.Uint32(data[12:16])
	b.duration = defaultEndianness.Uint32(data[16:20])
	b.rate = fixed32(data[20:24])
	b.volume = fixed16(data[24:26])

	return nil
}

type mvhdBoxFactory struct {
}

// Name returns the name of the type.
func (mvhdBoxFactory) Name() string {
	return "mvhd"
}

// New returns a new value instance.
func (mvhdBoxFactory) New(box Box) (cb CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	mvhdBox := &MvhdBox{
		Box: box,
	}

	err = mvhdBox.parse()
	log.PanicIf(err)

	return mvhdBox, nil
}

var (
	_ boxFactory = mvhdBoxFactory{}
	_ CommonBox  = &MvhdBox{}
)

func init() {
	RegisterBoxType(mvhdBoxFactory{})
}
