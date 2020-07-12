package atom

import (
	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

// MvhdBox - Movie Header Box
// Box Type: mvhd
// Container: Movie Box (moov)
// Mandatory: Yes
// Quantity: Exactly one.
//
// This box defines overall information which is media-independent,
// and relevant to the entire presentationconsidered as a whole.
type MvhdBox struct {
	*Box

	Flags            uint32
	Version          uint8
	CreationTime     uint32
	ModificationTime uint32
	Timescale        uint32
	Duration         uint32
	Rate             Fixed32
	Volume           Fixed16
}

func (b *MvhdBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.readBoxData()
	log.PanicIf(err)

	b.Version = data[0]
	b.Timescale = binary.BigEndian.Uint32(data[12:16])
	b.Duration = binary.BigEndian.Uint32(data[16:20])
	b.Rate = fixed32(data[20:24])
	b.Volume = fixed16(data[24:26])

	return nil
}

type mvhdBoxFactory struct {
}

// Name returns the name of the type.
func (mvhdBoxFactory) Name() string {
	return "mvhd"
}

// New returns a new value instance.
func (mvhdBoxFactory) New(box *Box) (cb CommonBox, err error) {
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
	_ CommonBox  = MvhdBox{}
)

func init() {
	registerAtom(mvhdBoxFactory{})
}
