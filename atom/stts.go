package atom

import (
	"github.com/dsoprea/go-logging"
)

// SttsBox is a "Decoding Time to Sample" box.
type SttsBox struct {
	Box

	version      byte
	flags        uint32
	entryCount   uint32
	sampleCounts []uint32
	sampleDeltas []uint32
}

func (sb *SttsBox) Version() byte {
	return sb.version
}

func (sb *SttsBox) Flags() uint32 {
	return sb.flags
}

func (sb *SttsBox) EntryCount() uint32 {
	return sb.entryCount
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

	data, err := b.readBoxData()
	log.PanicIf(err)

	b.version = data[0]
	b.flags = defaultEndianness.Uint32(data[0:4])

	count := defaultEndianness.Uint32(data[4:8])
	b.sampleCounts = make([]uint32, count)
	b.sampleDeltas = make([]uint32, count)

	for i := 0; i < int(count); i++ {
		b.sampleCounts[i] = defaultEndianness.Uint32(data[(8 + 8*i):(12 + 8*i)])
		b.sampleDeltas[i] = defaultEndianness.Uint32(data[(12 + 8*i):(16 + 8*i)])
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
func (sttsBoxFactory) New(box Box) (cb CommonBox, err error) {
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
	_ boxFactory = sttsBoxFactory{}
	_ CommonBox  = &SttsBox{}
)

func init() {
	registerAtom(sttsBoxFactory{})
}
