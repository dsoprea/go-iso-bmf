package bmftype

import (
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// SttsBox is a "Decoding Time to Sample" box.
type SttsBox struct {
	bmfcommon.Box

	version      byte
	flags        uint32
	sampleCounts []uint32
	sampleDeltas []uint32
}

func (sb *SttsBox) Version() byte {
	return sb.version
}

func (sb *SttsBox) Flags() uint32 {
	return sb.flags
}

func (sb *SttsBox) SampleCounts() []uint32 {
	return sb.sampleCounts
}

func (sb *SttsBox) SampleDeltas() []uint32 {
	return sb.sampleDeltas
}

func (b *SttsBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.ReadBoxData()
	log.PanicIf(err)

	b.version = data[0]
	b.flags = bmfcommon.DefaultEndianness.Uint32(data[0:4])

	count := bmfcommon.DefaultEndianness.Uint32(data[4:8])
	b.sampleCounts = make([]uint32, count)
	b.sampleDeltas = make([]uint32, count)

	offset := 8
	for i := 0; i < int(count); i++ {
		b.sampleCounts[i] = bmfcommon.DefaultEndianness.Uint32(data[offset : offset+4])
		b.sampleDeltas[i] = bmfcommon.DefaultEndianness.Uint32(data[offset+4 : offset+8])

		offset += 8
	}

	return nil
}

type sttsBoxFactory struct {
}

// Name returns the name of the type.
func (sttsBoxFactory) Name() string {
	return "stts"
}

// New returns a new value instance.
func (sttsBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	sttsBox := &SttsBox{
		Box: box,
	}

	err = sttsBox.parse()
	log.PanicIf(err)

	return sttsBox, nil
}

var (
	_ bmfcommon.BoxFactory = sttsBoxFactory{}
	_ bmfcommon.CommonBox  = &SttsBox{}
)

func init() {
	bmfcommon.RegisterBoxType(sttsBoxFactory{})
}
