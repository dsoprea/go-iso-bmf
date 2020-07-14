package atom

import (
	"github.com/dsoprea/go-logging"
)

// HmhdBox is a "Hint Media Header" box.
//
// Contains general information, independent of the protocol, for hint tracks.
type HmhdBox struct {
	Box

	version    byte
	maxPDUSize uint16
	avgPDUSize uint16
	maxBitrate uint32
	avgBitrate uint32
}

func (hb *HmhdBox) Version() byte {
	return hb.version
}

func (hb *HmhdBox) MaxPDUSize() uint16 {
	return hb.maxPDUSize
}

func (hb *HmhdBox) AvgPDUSize() uint16 {
	return hb.avgPDUSize
}

func (hb *HmhdBox) MaxBitrate() uint32 {
	return hb.maxBitrate
}

func (hb *HmhdBox) AvgBitrate() uint32 {
	return hb.avgBitrate
}

func (b *HmhdBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.readBoxData()
	log.PanicIf(err)

	b.version = data[0]
	b.maxPDUSize = defaultEndianness.Uint16(data[0:2])
	b.avgPDUSize = defaultEndianness.Uint16(data[2:4])
	b.maxBitrate = defaultEndianness.Uint32(data[4:8])
	b.avgBitrate = defaultEndianness.Uint32(data[8:12])

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
