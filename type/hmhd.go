package bmftype

import (
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// HmhdBox is a "Hint Media Header" box.
//
// Contains general information, independent of the protocol, for hint tracks.
type HmhdBox struct {
	bmfcommon.Box

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

	data, err := b.ReadBoxData()
	log.PanicIf(err)

	// TODO(dustin): !! We're not convinced this is right (or the version for any of the boxes, for that fact; they all steal a byte from the adjacent field).
	b.version = data[0]

	b.maxPDUSize = bmfcommon.DefaultEndianness.Uint16(data[0:2])
	b.avgPDUSize = bmfcommon.DefaultEndianness.Uint16(data[2:4])
	b.maxBitrate = bmfcommon.DefaultEndianness.Uint32(data[4:8])
	b.avgBitrate = bmfcommon.DefaultEndianness.Uint32(data[8:12])

	return nil
}

type hmhdBoxFactory struct {
}

// Name returns the name of the type.
func (hmhdBoxFactory) Name() string {
	return "hmhd"
}

// New returns a new value instance.
func (hmhdBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, err error) {
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
	_ bmfcommon.BoxFactory = hmhdBoxFactory{}
	_ bmfcommon.CommonBox  = &HmhdBox{}
)

func init() {
	bmfcommon.RegisterBoxType(hmhdBoxFactory{})
}
