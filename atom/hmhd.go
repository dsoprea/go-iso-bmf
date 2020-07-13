package atom

import (
	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

// HmhdBox is a "Hint Media Header" box.
//
// Contains general information, independent of the protocol, for hint tracks.
type HmhdBox struct {
	Box

	Version    byte
	MaxPDUSize uint16
	AvgPDUSize uint16
	MaxBitrate uint32
	AvgBitrate uint32
}

func (b *HmhdBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.readBoxData()
	log.PanicIf(err)

	b.Version = data[0]
	b.MaxPDUSize = binary.BigEndian.Uint16(data[0:2])
	b.AvgPDUSize = binary.BigEndian.Uint16(data[2:4])
	b.MaxBitrate = binary.BigEndian.Uint32(data[4:8])
	b.AvgBitrate = binary.BigEndian.Uint32(data[8:12])

	return nil
}

type hmhdBoxFactory struct {
}

// Name returns the name of the type.
func (hmhdBoxFactory) Name() string {
	return "hmhd"
}

// New returns a new value instance.
func (hmhdBoxFactory) New(box Box) (cb CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	hmhdBox := &HmhdBox{
		Box: box,
	}

	err = hmhdBox.parse()
	log.PanicIf(err)

	return hmhdBox, nil
}

var (
	_ boxFactory = hmhdBoxFactory{}
	_ CommonBox  = &HmhdBox{}
)

func init() {
	registerAtom(hmhdBoxFactory{})
}
